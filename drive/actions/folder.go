package actions

import (
	"fmt"
	"net/http"

	"github.com/ihleven/cloud11-api/drive"
)

type DirActionResponder struct {
	*drive.Folder
	wd          drive.WebDrive `json:"-"`
	Breadcrumbs string
}

func (a DirActionResponder) Handle(w http.ResponseWriter, r *http.Request) error {

	var err error
	switch r.Method {
	case http.MethodGet:
		err = a.GetAction(r, w)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a DirActionResponder) GetAction(r *http.Request, w http.ResponseWriter) error {
	fmt.Println(" * DirActionResponder: diraction")
	a.Breadcrumbs = "folder"
	return nil
}
