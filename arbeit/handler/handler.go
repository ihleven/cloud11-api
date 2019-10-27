package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/ihleven/cloud11-api/arbeit"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "PUT, DELETE, GET, HEAD")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

// ArbeitHandler ...
type ArbeitHandler struct {
	//domain arbeit.Domain
}

func (a ArbeitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)
	//sessionUser, _ := session.GetSessionUser(r, w)

	var err error

	params := parseURL(r.URL.Path)

	switch {
	case params.day != 0:
		ArbeitstagAction(w, r, params)

	case params.week != 0:
		//err = params.HandleArbeitWoche(w, r, params)

	case params.month != 0:
		handleArbeitsmonat(w, r, params)

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

func ArbeitstagAction(w http.ResponseWriter, r *http.Request, p *params) {

	switch r.Method {
	case http.MethodPut:

		var body arbeit.Arbeitstag
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {

			http.Error(w, err.Error(), 400)
			return
		}

		err = arbeit.UpdateArbeitstag(p.year, p.month, p.day, 1, &body)
		if err != nil {
			http.Error(w, err.Error(), 500) //return errors.Wrapf(err, "could not update arbeitstag %v", body)
			return
		}

		fallthrough

	case http.MethodGet:

		arbeitstag, err := arbeit.GetArbeitstag(p.year, p.month, p.day, 1)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), 500)
			return //errors.Wrapf(err, "could not read arbeitstag %d/%d/%d %d", p.year, p.month, p.day, 1)
		}
		//json.NewEncoder(w).Encode(arbeitstag)
		js, err := json.MarshalIndent(arbeitstag, "", "\t")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	default:
		//http.Error(w, "only get or put", 400)
		enableCors(&w)
	}
}

// func handleArbeitswoche(w http.ResponseWriter, r *http.Request, p *params) {
// 	arbeitswoche, err := arbeit.RetrieveArbeitsWoche(params.year, params.week, 1)
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 	}
// 	json.NewEncoder(w).Encode(arbeitswoche)
// }

func handleJahr(w http.ResponseWriter, r *http.Request, p *params) {

	arbeitsjahr, err := arbeit.GetArbeitsjahr(p.year, 1)
	if err != nil {
		http.Error(w, "error with GetArbeitsJahr", 500)
	}
	json.NewEncoder(w).Encode(arbeitsjahr)
}

func handleArbeitsmonat(w http.ResponseWriter, r *http.Request, p *params) {
	fmt.Println("handleArbeitmonat", p.year, p.month)
	arbeitsmonat, err := arbeit.GetArbeitsmonat(p.year, p.month, 1)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	json.NewEncoder(w).Encode(arbeitsmonat)
}
