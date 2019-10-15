package http

import (
	"net/http"
	"time"

	arbeithandler "github.com/ihleven/cloud11-api/arbeit/handler"
)

// NewServer returns a new instance of Server.
func NewServer(address string) *http.Server {
	handler := arbeithandler.ArbeitHandler{}
	return &http.Server{
		Handler:      handler,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
