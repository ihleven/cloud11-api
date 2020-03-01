package webserver

import (
	"net/http"
	"path"
	"strings"
)

func NewDispatcher(notFoundHandler http.Handler) *shiftPathDispatcher {
	routes := make(map[string]http.Handler)
	if notFoundHandler == nil {
		notFoundHandler = http.HandlerFunc(http.NotFound)
	}
	return &shiftPathDispatcher{routes, notFoundHandler}
}

type shiftPathDispatcher struct {
	routes   map[string]http.Handler
	notFound http.Handler
}

func (d shiftPathDispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	head, tail := shiftPath(r.URL.Path)

	route, ok := d.routes[head]
	switch {
	case ok:
		r.URL.Path = tail
		route.ServeHTTP(w, r)
	default:
		d.notFound.ServeHTTP(w, r)
	}
}

func (d *shiftPathDispatcher) Register(route string, handler http.Handler) {

	var head, tail string
	head, tail = shiftPath(route)
	switch tail {
	case "/":
		d.routes[head] = handler
	default:
		d.registerSubRoute(head, tail[1:], handler)
	}
}

func (d *shiftPathDispatcher) registerSubRoute(head, tail string, handler http.Handler) {

	subDispatcher, ok := d.routes[head].(*shiftPathDispatcher)
	if !ok {
		subDispatcher = NewDispatcher(d.routes[head])
		d.routes[head] = subDispatcher
	}
	subDispatcher.Register(tail, handler)
}

func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
