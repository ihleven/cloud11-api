package auth

import "net/http"

// Responder
type Responder interface {
	Respond(http.ResponseWriter, *http.Request, interface{}) error
}
type ResponderFunc func(http.ResponseWriter, *http.Request, interface{}) error

// ServeHTTP calls f(w, r).
func (f ResponderFunc) Respond(w http.ResponseWriter, r *http.Request, data interface{}) error {
	return f(w, r, data)
}
