package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ihleven/cloud11-api/auth"
	"github.com/ihleven/cloud11-api/drive"
	"github.com/pkg/errors"
)

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, DELETE, GET, HEAD")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}
func Dispatch(wd drive.WebDrive) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Dispatch", r.URL.Path)
		enableCors(w)

		cleanedPath := filepath.Clean(strings.Replace(r.URL.Path, "|", ".", 1))

		file, err := wd.GetFile(cleanedPath, auth.CurrentUser)
		if err != nil {
			if err := errors.Cause(err); os.IsNotExist(err) {
				http.Error(w, err.Error(), 404)
			} else {
				http.Error(w, err.Error(), 500)
			}
			return
		}

		//filer, err := getFiler(cleanedPath)
		//
		// da := DriveAction{File: file, wd: wd, Account: auth.CurrentUser}
		var action Actioneer

		switch {
		case file.Type.Filetype == "F":
			action = &FileAction{File: file, wd: wd}

		case file.Type.Filetype == "D":
			action = &DirActionResponder{File: file, wd: wd}
		}
		err = action.Handle(w, r)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		SerializeJSON(w, r, action)

	}
}

type Actioneer interface {
	Handle(http.ResponseWriter, *http.Request) error
}

type HTTPGetter interface {
	HandleGet(*http.Request) error
}

func SerializeJSON(w http.ResponseWriter, r *http.Request, data interface{}) (err error) {

	// https://stackoverflow.com/questions/37863374/whats-the-difference-between-responsewriter-write-and-io-writestring
	js, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}
