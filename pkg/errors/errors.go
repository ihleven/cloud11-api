package errors

import (
	"fmt"
	"io"
)

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// Errorf also records the stack trace at the point it was called.
func New(format string, args ...interface{}) error {
	return &fundamental{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}
func Code(code ErrorCode, format string, args ...interface{}) error {
	return &fundamental{
		msg:   fmt.Sprintf(format, args...),
		code:  code,
		stack: callers(),
	}
}

// fundamental is an error that has a message and a stack, but no caller.

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &withStack{
		err,
		callers(),
	}
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func Wrap(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	e := &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
	file, function, line, ok := caller()
	if ok {
		e.file = file
		e.function = function
		e.line = line
	}
	return &withStack{
		e,
		callers(),
	}
}

// WithMessagef annotates err with the format specifier.
// If err is nil, WithMessagef returns nil.
func WithMessage(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	wm := withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
	file, function, line, ok := caller()
	if ok {
		wm.file = file
		wm.function = function
		wm.line = line
	}

	return &wm
}

func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

type fundamental struct {
	msg  string
	code ErrorCode
	*stack
}

func (f *fundamental) Error() string { return f.msg }
func (f *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, f.msg)
			fmt.Fprintf(s, "\n   --- at:")
			f.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, f.msg)
	case 'q':
		fmt.Fprintf(s, "%q", f.msg)
	}
}

type withStack struct {
	error
	*stack
}

func (w *withStack) Cause() error { return w.error }
func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Cause())
			// fmt.Fprintf(s, "\nStacktrace: ")
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

type withMessage struct {
	cause    error
	msg      string
	file     string
	function string
	line     int
}

func (w *withMessage) Error() string { return w.msg + ": " + w.cause.Error() }
func (w *withMessage) Cause() error  { return w.cause }
func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// io.WriteString(s, w.msg)
			fmt.Fprintf(s, "%s\n", w.msg)
			// w.frame.Format(s, verb)
			fmt.Fprintf(s, "   --- at %v:%v (%v) ---\n", w.file, w.line, w.function)
			if w.cause != nil {
				fmt.Fprintf(s, " Caused by: ")

			}
			fmt.Fprintf(s, "%+v\n", w.Cause())
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

func Dump(err error) {
	fmt.Println("Dump:")
	if e, ok := err.(*fundamental); ok {
		fmt.Printf("fundamental: %+v\n", e.msg)
	} else if e, ok := err.(*withMessage); ok {
		fmt.Printf("withMessage: %s\n", e.msg)
		Dump(e.cause)
	} else if e, ok := err.(*withStack); ok {
		fmt.Printf("withStack: \n")
		Dump(e.error)
	} else {
		fmt.Println("other:", err)
	}

}
