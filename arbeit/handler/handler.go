package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/ihleven/cloud11-api/arbeit"

	"github.com/pkg/errors"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "PUT, DELETE, GET, HEAD")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

// ArbeitHandler ...
type ArbeitHandler struct {
}

func (a ArbeitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)
	//sessionUser, _ := session.GetSessionUser(r, w)

	var err error

	params := parseURL(r.URL.Path)

	switch {
	case params.day != 0:
		err = handleArbeitstag(w, r, params)

	case params.week != 0:
		//err = params.HandleArbeitWoche(w, r, params)

	case params.month != 0:
		//err = params.HandleMonat(w, r, params)

	case params.year != 0:
		handleJahr(w, r, params)

	default:
		http.NotFound(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

type params struct {
	year, month, week, day int
}

func parseURL(urlpath string) *params {
	var params params
	var regex = regexp.MustCompile(`arbeit\/(\d{4})\/?(kw|KW|Kw|w|W)?(\d{1,2})?\/?(\d{1,2})?`)
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

func handleArbeitstag(w http.ResponseWriter, r *http.Request, p *params) error {

	switch r.Method {
	case http.MethodPut:

		var body arbeit.Arbeitstag
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return nil
		}
		fmt.Printf("parsed Arbeitstag: %v", body)
		err = arbeit.UpdateArbeitstag(p.year, p.month, p.day, 1, &body)
		if err != nil {
			fmt.Printf("error Arbeitstag: %v", err)
			return errors.Wrapf(err, "could not update arbeitstag %v", body)
		} else {
			fmt.Println("success update at")
		}

		fallthrough

	case http.MethodGet:

		arbeitstag, err := arbeit.GetArbeitstag(p.year, p.month, p.day, 1)
		if err != nil {
			return errors.Wrapf(err, "could not read arbeitstag %d/%d/%d %d", p.year, p.month, p.day, 1)
		}
		//json.NewEncoder(w).Encode(arbeitstag)
		js, err := json.MarshalIndent(arbeitstag, "", "\t")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	default:
		//http.Error(w, "only get or put", 400)
		enableCors(&w)
	}

	return nil
}

// func handleArbeitswoche(w http.ResponseWriter, r *http.Request, p *params) {
// 	arbeitswoche, err := arbeit.RetrieveArbeitsWoche(params.year, params.week, 1)
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 	}
// 	json.NewEncoder(w).Encode(arbeitswoche)
// }

// func handleArbeitsmonat(w http.ResponseWriter, r *http.Request, p *params) {
// 	arbeitsmonat, err := arbeit.GetArbeitsMonat(params.year, params.month, 1)
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 	}
// 	json.NewEncoder(w).Encode(arbeitsmonat)
// }

func handleJahr(w http.ResponseWriter, r *http.Request, p *params) {

	arbeitsjahr, err := arbeit.GetArbeitsjahr(p.year, 1)
	if err != nil {
		http.Error(w, "error with GetArbeitsJahr", 500)
	}
	json.NewEncoder(w).Encode(arbeitsjahr)
}
