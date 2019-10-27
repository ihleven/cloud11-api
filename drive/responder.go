package drive

import (
	"encoding/json"
	"net/http"
)

// Responder
type Responder interface {
	Respond(http.ResponseWriter, *http.Request, interface{}) error
}
type ResponderFunc func(http.ResponseWriter, *http.Request, interface{}) error

// ServeHTTP calls f(w, r).
func (f ResponderFunc) Respond(w http.ResponseWriter, r *http.Request, data interface{}) error {
	return f(w, r, data)
}

// TemplateResponder
type JSONResponder struct{}

func SerializeJSON(w http.ResponseWriter, r *http.Request, data interface{}) (err error) {

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

// TemplateResponder
type TemplateResponder struct {
	Template string
}

func (resp *TemplateResponder) Respond(w http.ResponseWriter, r *http.Request, data map[string]interface{}) (err error) {

	switch r.Header.Get("Accept") {
	case "application/json":
		// err = templates.SerializeJSON(w, http.StatusOK, data)
	default:
		// err = templates.Render(w, http.StatusOK, resp.Template, data)
	}

	if err != nil {
		// errors.Error(w, r, errors.Wrap(err, "render error"))
	}
	return
}

//func (resp *TemplateResponder) Render(w http.ResponseWriter, status int, data map[string]interface{}) error {
//
//	return Render(w, status, resp.template, data)
//}
