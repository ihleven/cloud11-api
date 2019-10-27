package http

import (
	"fmt"
	"net/http"
	"time"

	arbeithandler "github.com/ihleven/cloud11-api/arbeit/handler"
	"github.com/ihleven/cloud11-api/drive/actions"
	"github.com/ihleven/cloud11-api/drive/fs"
)

// NewServer returns a new instance of Server.
func NewServer(address string) *http.Server {
	handler := arbeithandler.ArbeitHandler{}

	mux := http.NewServeMux()

	//mux.Handle("/assets/", assetHandler("assets", "_static/assets"))
	mux.Handle("/arbeit/", handler)
	//mux.HandleFunc("/login", Login)
	//mux.HandleFunc("/logout", Logout)
	mux.HandleFunc("/home/", actions.Dispatch(&fs.Drive))
	mux.HandleFunc("/serve/home/", fs.Serve(fs.Drive))
	mux.HandleFunc("/static/", fs.DispatchRaw(fs.Drive, "/static"))

	return &http.Server{
		Handler:      mux,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func assetHandler(prefix, location string) http.Handler {
	//fmt.Println("assetHandler", prefix, location)
	return http.StripPrefix(fmt.Sprintf("/%s/", prefix), http.FileServer(http.Dir(fmt.Sprintf("./%s", location))))
}
