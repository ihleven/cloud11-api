package actions

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ihleven/cloud11-api/drive"
	"github.com/pkg/errors"
)

type FileAction struct {
	*drive.File
	wd      drive.WebDrive     `json:"-"`
	Crumbs  []drive.Breadcrumb `json:"breadcrumbs"`
	Parents []drive.Breadcrumb `json:"parents"`
	Content string             `json:"content"`
}

func (a *FileAction) Handle(w http.ResponseWriter, r *http.Request) error {
	fmt.Println(" * fileAction.Handle()")

	var err error
	switch r.Method {
	case http.MethodPut:
		err = a.PutAction(r, w)
		if err != nil {
			return err
		}
		fallthrough
	case http.MethodGet:
		err = a.GetAction(r, w)
		if err != nil {
			return err
		}
	case http.MethodPost:
		err = a.PostAction(r, w)
		if err != nil {
			return err
		}
	}
	a.Crumbs = a.wd.GenerateBreadcrumbs(a.URL)
	a.Parents = a.wd.GenerateParents(a.URL)
	return nil
}

func (a *FileAction) GetContent() ([]byte, error) { //offset, limit int) (e error) {

	var content = make([]byte, a.File.Size)

	// fd := a.File.OpenFile(0)
	// defer fd.Close()
	// fd.Seek(0, 0)

	bytes, err := a.File.Read(content)
	if err != nil {
		return nil, err
	}

	if int64(bytes) != a.File.Size {
		return content, errors.Errorf("read only %d of %d bytes", bytes, a.File.Size)
	}
	return content, nil
}

func (a *FileAction) GetAction(r *http.Request, w http.ResponseWriter) error {

	fmt.Println(" * FileActionResponder: GetAction")

	if a.File.Type.Mediatype == "text" {

		var content = make([]byte, a.File.Size)

		_, err := a.File.Read(content)
		if err != nil {
			return err
		}

		a.Content = string(content)
	}
	return nil
}

func (a *FileAction) PutAction(r *http.Request, w http.ResponseWriter) error {

	decoder := json.NewDecoder(r.Body)
	var f FileAction
	err := decoder.Decode(&f)
	if err != nil {
		return err
	}
	if f.Content != "" {
		io.WriteString(a.File, f.Content)
	}

	// body, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	return err
	// }
	// //i.update(body)
	// fmt.Println("CONTENT:", string(body))
	return nil
}

func (a *FileAction) PostAction(r *http.Request, w http.ResponseWriter) error {

	fmt.Println("\n\n * PostAction", a.File.Authorization)

	if !a.Authorization.W {
		return fmt.Errorf("no write permissions")
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	// fmt.Println("body:", string(body))

	_, err = a.File.Write(body)
	if err != nil {
		return err
	}
	// fmt.Println("bytes:", n, err)

	// buf := new(bytes.Buffer)
	// buf.ReadFrom(a.File)
	// fmt.Printf("content: %v", buf.String())

	// formfile, multipart, err := r.FormFile("file")
	// if err != nil {
	// 	return errors.Wrap(err, "parsing form")
	// }
	// defer formfile.Close()
	// _ = multipart.Header.Get("Content-Type")

	// _, err = ioutil.ReadAll(formfile)
	// if err != nil {
	// 	return errors.Wrap(err, "read form file")
	// }

	// err = a.File.SetUTF8Content(data)
	// if err != nil {
	// 	return errors.Wrap(err, "writing utf8 file")
	// }

	http.Redirect(w, r, a.File.URL, http.StatusFound)
	return nil
}
