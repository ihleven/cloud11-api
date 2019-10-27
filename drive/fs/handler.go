package fs

import (
	"net/http"
	"os"
	"path"

	"github.com/ihleven/cloud11-api/auth"
)

func Serve(webdrive FSWebDrive) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		//authuser, err := session.GetSessionUser(r, w)

		info, fd, err := webdrive.GetServeHandle(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer fd.Close()

		if info.IsDir() {
			r.URL.Path = path.Join(r.URL.Path, "index.html")
			Serve(webdrive)(w, r)
			return
		}

		http.ServeContent(w, r, info.Name(), info.ModTime(), fd)
	}
}

func DispatchRaw(webdrive FSWebDrive, prefix string) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		fh, err := webdrive.GetHandle(r.URL.Path, prefix)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}

		authuser := auth.CurrentUser //, err := session.GetSessionUser(r, w)

		if !fh.(*handle).HasReadPermission(nil) {
			http.Error(w, "Account '"+authuser.Username+"' has no read permission", 403)
			return
		}

		if fh.Mode().IsDir() {
			r.URL.Path = path.Join(r.URL.Path, "index.html")
			Serve(webdrive)(w, r)
			return
		}

		fd, _ := fh.OpenFile(os.O_RDONLY, 0)
		defer fd.Close()

		http.ServeContent(w, r, fh.Name(), fh.ModTime(), fd)
	}
}
