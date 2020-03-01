package errors

import (
	"runtime"
	"strings"
)

/*
CleanPath function is applied to file paths before adding them to a stacktrace.
By default, it makes the path relative to the $GOPATH environment variable.
To remove some additional prefix like "github.com" from file paths in
stacktraces, use something like:
	stacktrace.CleanPath = func(path string) string {
		path = RemoveGoPath(path)
		path = strings.TrimPrefix(path, "github.com/")
		return path
	}
*/
var CleanPath = RemoveGoPath

type ErrorCode uint16

const NoCode ErrorCode = 500
const BadRequest ErrorCode = 400
const NotFound ErrorCode = 404

func GetCode(err error) (ErrorCode, string) {
	cause := Cause(err)
	if err, ok := cause.(*fundamental); ok {
		return err.code, err.msg
	}
	return NoCode, ""
}

// type stacktrace struct {
// 	message  string
// 	cause    error
// 	code     ErrorCode
// 	file     string
// 	function string
// 	line     int
// }

// func create(cause error, code ErrorCode, msg string, vals ...interface{}) error {
// 	// If no error code specified, inherit error code from the cause.
// 	if code == NoCode {
// 		code = GetCode(cause)
// 	}

// 	err := &stacktrace{
// 		message: fmt.Sprintf(msg, vals...),
// 		cause:   cause,
// 		code:    code,
// 	}

// 	// Caller of create is NewError or Propagate, so user's code is 2 up.
// 	pc, file, line, ok := runtime.Caller(2)
// 	if !ok {
// 		return err
// 	}
// 	if CleanPath != nil {
// 		file = CleanPath(file)
// 	}
// 	err.file, err.line = file, line

// 	f := runtime.FuncForPC(pc)
// 	if f == nil {
// 		return err
// 	}
// 	err.function = shortFuncName(f)

// 	return err
// }
func caller() (file string, function string, line int, ok bool) {

	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return
	}
	if CleanPath != nil {
		file = CleanPath(file)
	}

	f := runtime.FuncForPC(pc)
	if f == nil {
		return
	}
	function = shortFuncName(f)

	return
}

/* "FuncName" or "Receiver.MethodName" */
func shortFuncName(f *runtime.Func) string {
	// f.Name() is like one of these:
	// - "github.com/palantir/shield/package.FuncName"
	// - "github.com/palantir/shield/package.Receiver.MethodName"
	// - "github.com/palantir/shield/package.(*PtrReceiver).MethodName"
	longName := f.Name()

	withoutPath := longName[strings.LastIndex(longName, "/")+1:]
	withoutPackage := withoutPath[strings.Index(withoutPath, ".")+1:]

	shortName := withoutPackage
	shortName = strings.Replace(shortName, "(", "", 1)
	shortName = strings.Replace(shortName, "*", "", 1)
	shortName = strings.Replace(shortName, ")", "", 1)

	return shortName
}

// func (st *stacktrace) Error() string {
// 	return fmt.Sprint(st)
// }

// ExitCode returns the exit code associated with the stacktrace error based on its error code. If the error code is
// NoCode, return 1 (default); otherwise, returns the value of the error code.
// func (st *stacktrace) ExitCode() int {
// 	if st.code == NoCode {
// 		return 1
// 	}
// 	return int(st.code)
// }
