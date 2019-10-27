package actions

import (
	"fmt"
	"net/http"
)

type FileAction struct {
}

func (a FileAction) Handle(w http.ResponseWriter, r *http.Request) error {
	fmt.Println(" * fileAction.Handle()")
	var err error
	switch r.Method {
	case http.MethodGet:
		//err = ar.GetAction(r, w)

	}
	if err != nil {
		return err
	}

	return nil
}
