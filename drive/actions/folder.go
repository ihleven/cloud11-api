package actions

import (
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/ihleven/cloud11-api/auth"
	"github.com/ihleven/cloud11-api/drive"
	"github.com/pkg/errors"
)

type DirActionResponder struct {
	*drive.File
	wd      drive.WebDrive     `json:"-"`
	Crumbs  []drive.Breadcrumb `json:"breadcrumbs"`
	Parents []drive.Breadcrumb `json:"parents"`
	// FileAction
	Entries []drive.File `json:"entries"`
}

func (a *DirActionResponder) Handle(w http.ResponseWriter, r *http.Request) (err error) {

	fmt.Println(" * DirActionResponder: Handle", r.Method)

	if r.Method == http.MethodPost {

		err = a.UploadFilesPostAction(r, w)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	if !a.File.Authorization.R {
		return fmt.Errorf("Missing read permissions for %s", a.File.URL)
	}

	a.Entries, err = a.wd.GetFiles(a.File, auth.CurrentUser)
	if err != nil {
		return err
	}
	a.Parents = a.wd.GenerateParents(a.URL)
	a.Crumbs = a.wd.GenerateBreadcrumbs(a.URL)

	return nil
}
func (a *DirActionResponder) UploadFilesPostAction(r *http.Request, w http.ResponseWriter) error {

	folder := a.File
	if !folder.Authorization.W {
		return fmt.Errorf("Missing write permissions for %s", folder.URL)
	}
	err := r.ParseMultipartForm(2000000)
	if err != nil {
		return err
	}

	formdata := r.MultipartForm

	for _, header := range formdata.File["files"] {

		file, err := header.Open()
		defer file.Close()
		if err != nil {
			return errors.Wrapf(err, "Could not open form file %v", header)
		}

		h, err := a.wd.CreateFile(folder, header.Filename)
		if err != nil {
			return errors.Wrapf(err, "Could not upload to folder '%v'. Unable to create the file for writing. Check your write access privilege", header.Filename)
		}

		_, err = io.Copy(h, file)
		if err != nil {
			return errors.Wrapf(err, "Unable to copy formfile")
		}
		fmt.Println(" * uploaded", h)
	}
	//for key, value := range formdata.Value["foo"] {
	fmt.Println(" * value:", formdata.Value["foo"])
	for _, name := range formdata.Value["folders"] {
		path := path.Join(folder.URL, name)
		fmt.Println(" * folder:", folder, path)
		fh, err := a.wd.Mkdir(path)
		fmt.Println("%v %v", fh, err)
	}

	//}
	return nil
}

// func (a *DirActionResponder) PutAction(r *http.Request, w http.ResponseWriter) error {

// 	file := a.File
// 	fmt.Printf("PutAction => Directory \"%s/\"\n", file.Name)

// 	//if !file.Permissions.Write {
// 	//	return errors.Errorf("no write permissions")
// 	//}

// 	var options struct {
// 		CreateThumbnails bool
// 	}
// 	err := json.NewDecoder(r.Body).Decode(&options)
// 	if err != nil {
// 		return errors.Wrap(err, "Error decoding put request body")
// 	}

// 	if options.CreateThumbnails {

// 		err := drive.MakeThumbs(file.Handle)
// 		if err != nil {
// 			return errors.Wrap(err, "Error making thumbnails")
// 		}
// 	}
// 	return nil
// }
