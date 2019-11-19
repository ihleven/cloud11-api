package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	arbeithandler "github.com/ihleven/cloud11-api/arbeit/handler"
	"github.com/ihleven/cloud11-api/drive"
	"github.com/pkg/errors"

	"github.com/ihleven/cloud11-api/drive/fs"
	"github.com/ihleven/cloud11-api/drive/hidrive"
)

// NewServer returns a new instance of Server.
func NewServer(address string) *http.Server {

	token, err := hidrive.NewToken()
	if err != nil {
		fmt.Println("token error:", err)
	}
	hidrive.HIDrive.Token = *token

	var router = Router{
		arbeit:  arbeithandler.ArbeitHandler{},
		home:    drive.Dispatch(&fs.Drive),
		serve:   drive.DispatchRaw(&fs.Drive),
		hidrive: drive.DispatchHandler(&hidrive.HIDrive),
		hiserve: hidrive.DispatchRaw(hidrive.HIDrive),
	}

	//mux.Handle("/assets/", assetHandler("assets", "_static/assets"))
	//mux.HandleFunc("/login", Login)
	//mux.HandleFunc("/logout", Logout)
	// mux.HandleFunc("/home/", actions.Dispatch(&fs.Drive))
	// mux.HandleFunc("/serve/home/", fs.Serve(fs.Drive))
	// mux.HandleFunc("/static/", fs.DispatchRaw(fs.Drive, "/static"))

	return &http.Server{
		Handler:      &router,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

type Router struct {
	arbeit http.Handler
	//home   func(*http.Request) (*drive.DriveAction, error) //http.HandlerFunc
	home    func(*http.Request) (*drive.DriveAction, error)
	serve   http.HandlerFunc
	hidrive http.HandlerFunc
	hiserve http.HandlerFunc
}

func (r Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	var head string
	//var user = auth.CurrentUser
	var data interface{}
	var err error

	head, req.URL.Path = ShiftPath(req.URL.Path)
	switch head {

	case "serve":
		head, req.URL.Path = ShiftPath(req.URL.Path)
		r.serve(res, req)

	case "home":
		data, err = r.home(req) //, user)

	case "arbeit":
		r.arbeit.ServeHTTP(res, req)

	case "hidrive":
		r.hidrive(res, req)
		//_, err = hidrive.GetHandle(req.URL.Path)
	case "hiserve":
		r.hiserve(res, req)
		//_, err = hidrive.GetHandle(req.URL.Path)

	default:
		http.Error(res, "Not Found (route)", http.StatusNotFound)
	}
	if checkHTTPError(res, err) {
		return
	}
	if data != nil {
		err = respond(res, req, data)

	}
	if err != nil {
		http.Error(res, "Not Found (route)", http.StatusInternalServerError)
	}
}

func checkHTTPError(w http.ResponseWriter, err error) bool {
	if err != nil {
		status := http.StatusInternalServerError
		cause := errors.Cause(err)
		if os.IsNotExist(cause) {
			status = http.StatusNotFound
		} else if os.IsExist(cause) {
			status = http.StatusInternalServerError
		} else if os.IsPermission(cause) {
			status = http.StatusForbidden
		} else if e, ok := cause.(*os.PathError); ok {
			switch e {

			case os.ErrClosed:
				status = http.StatusGone
			case os.ErrNoDeadline:
				status = http.StatusInternalServerError
			}

			//http.Error(w, fmt.Sprintf("---%v %v %v", e.Op, e.Path, e.Err.Error()), 500)
		} else if cause.Error() == "Authemtication required" {
			status = http.StatusUnauthorized
		}
		http.Error(w, cause.Error(), status)
		return true
	}
	return false
}
func assetHandler(prefix, location string) http.Handler {
	//fmt.Println("assetHandler", prefix, location)
	return http.StripPrefix(fmt.Sprintf("/%s/", prefix), http.FileServer(http.Dir(fmt.Sprintf("./%s", location))))
}

func respond(w http.ResponseWriter, r *http.Request, data interface{}) (err error) {
	fmt.Println("respond:", data)
	enableCors(w)
	switch r.Header.Get("Accept") {
	case "application/json":
		err = SerializeJSON(w, http.StatusOK, data)
	default:
		err = SerializeJSON(w, http.StatusOK, data)
	}

	return
}
func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, DELETE, GET, HEAD")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}
func SerializeJSON(w http.ResponseWriter, status int, data interface{}) error {
	fmt.Println("SerializeJSON:", data)
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
