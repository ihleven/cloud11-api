package webserver

import (
	"fmt"
	"net/http"

	"time"

	"github.com/ihleven/cloud11-api/pkg/log"
)

func New(host string, port int) *httpServer {

	return &httpServer{
		addr:       fmt.Sprintf("%s:%d", host, port),
		dispatcher: NewDispatcher(nil),
	}
}

type httpServer struct {
	addr       string
	dispatcher *shiftPathDispatcher
	server     *http.Server
}

func (s *httpServer) ListenAndServe() {

	s.server = &http.Server{
		Addr:           s.addr,
		Handler:        s.dispatcher,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.server.ListenAndServe()
}

func (s *httpServer) Register(route string, handler interface{}) {

	switch handlerType := handler.(type) {
	case http.Handler:
		s.dispatcher.Register(route, handlerType)

	case func(w http.ResponseWriter, r *http.Request):
		s.dispatcher.Register(route, http.HandlerFunc(handlerType))

	default:
		log.Info("Could not register route '%v': unknown handler type %T", route, handler)
	}

}
