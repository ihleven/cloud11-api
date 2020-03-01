package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/ihleven/cloud11-api/arbeit"
	"github.com/ihleven/cloud11-api/pkg/errors"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}
func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

type Marshaler interface {
	Render(bool) ([]byte, error)
}

// ArbeitHandler ...
func ArbeitHandler(usecase *arbeit.Usecase) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		enableCors(&w)
		//sessionUser, _ := session.GetSessionUser(r, w)

		var m Marshaler
		var err error

		params := parseURL(r.URL.Path)

		switch {
		case params.day != 0:
			// m, err = usecase.GetArbeitstag(params.year, params.month, params.day, 1)
			m, err = ArbeitstagAction(w, r, params, usecase)

		case params.week != 0:
			//err = params.HandleArbeitWoche(w, r, params)

		case params.month != 0:
			m, err = usecase.GetArbeitsmonat(params.year, params.month, 1)

		case params.year != 0:
			m, err = usecase.Arbeitsjahr(params.year, 1)

		default:
			m, err = usecase.ListArbeitsjahre(1)
			// http.NotFound(w, r)
		}
		if err != nil {
			code, msg := errors.GetCode(err)
			fmt.Println(code, msg, err)
			http.Error(w, msg, int(code))
		} else {

			bytes, err := m.Render(true)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(bytes)

		}
	}
}

type params struct {
	year, month, week, day int
}

func parseURL(urlpath string) *params {
	fmt.Println(urlpath)
	var params params
	var regex = regexp.MustCompile(`\/(\d{4})\/?(kw|KW|Kw|w|W)?(\d{1,2})?\/?(\d{1,2})?`)
	matches := regex.FindStringSubmatch(urlpath)
	if matches != nil {
		if year, err := strconv.Atoi(matches[1]); err == nil {
			params.year = year
		}
		if month, err := strconv.Atoi(matches[3]); err == nil {
			if matches[2] != "" {
				params.week = month
			} else {
				params.month = month
			}
		}
		if day, err := strconv.Atoi(matches[4]); err == nil {
			params.day = day
		}
	}
	return &params
}

func ArbeitstagAction(w http.ResponseWriter, r *http.Request, p *params, usecase *arbeit.Usecase) (Marshaler, error) {
	switch r.Method {
	case http.MethodPut:

		var body arbeit.Arbeitstag
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			return nil, errors.Code(errors.BadRequest, "Could not decode request body")
		}
		arbeitstag, err := usecase.UpdateArbeitstag(p.year, p.month, p.day, 1, &body)
		fmt.Println("--- ArbeitstagAction", p, r.Method, arbeitstag, err)
		if err != nil {
			return nil, errors.WithMessage(err, "Could not update arbeitstag")
		}
		// return arbeitstag, nil
		fallthrough

	case http.MethodGet:

		arbeitstag, err := usecase.GetArbeitstag(p.year, p.month, p.day, 1)
		if err != nil {
			return nil, errors.Wrap(err, "could not read arbeitstag %d/%d/%d %d", p.year, p.month, p.day, 1)
		}
		return arbeitstag, nil
		// js, err := json.MarshalIndent(arbeitstag, "", "\t")
		// if err != nil {
		// 	return nil, errors.Wrap(err, "could not marshal arbeitstag %d/%d/%d %d", p.year, p.month, p.day, 1)
		// }

		// w.Header().Set("Content-Type", "application/json")
		// w.Write(js)

	default:
		fmt.Println("--- def", p, r.Method)

		return nil, errors.Code(errors.BadRequest, "Only get or put")
		// enableCors(&w)
	}
	return nil, nil
}

// func handleArbeitswoche(w http.ResponseWriter, r *http.Request, p *params) {
// 	arbeitswoche, err := arbeit.RetrieveArbeitsWoche(params.year, params.week, 1)
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 	}
// 	json.NewEncoder(w).Encode(arbeitswoche)
// }
