package arbeit

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func ArbeitHandler(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)
	//sessionUser, _ := session.GetSessionUser(r, w)

	var err error

	params := parseURL(r.URL.Path)

	switch {
	case params.day != 0:
		err = params.HandleArbeitstag(w, r, 1)
	case params.week != 0:
		arbeitswoche, err := Repo.RetrieveArbeitstage(params.year, 0, params.week, 1)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		json.NewEncoder(w).Encode(arbeitswoche)

	case params.month != 0:
		arbeitsmonat, err := GetArbeitsMonat(params.year, params.month, 1)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		json.NewEncoder(w).Encode(arbeitsmonat)

	case params.year != 0:
		arbeitsjahr, err := GetArbeitsjahr(params.year, 1)
		if err != nil {
			http.Error(w, "error with GetArbeitsJahr", 500)
		}
		json.NewEncoder(w).Encode(arbeitsjahr)

	default:
		http.Error(w, "could not parse url", 400)
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
	var regex = regexp.MustCompile(`arbeit\/(\d+)\/?(kw|KW|Kw|w|W)?(\d+)?\/?(\d+)?`)
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

func (p *params) HandleArbeitstag(w http.ResponseWriter, r *http.Request, accountID int) error {

	id := ((p.year*100+p.month)*100+p.day)*1000 + accountID

	switch r.Method {
	case http.MethodPut:

		var a Arbeitstag
		err := json.NewDecoder(r.Body).Decode(&a)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return nil
		}

		err = Repo.UpdateArbeitstag(id, &a)
		if err != nil {
			return errors.Wrapf(err, "could not update arbeitstag %d", id)
		}

		fallthrough

	case http.MethodGet:

		arbeitstag, err := Repo.ReadArbeitstag(id)
		if err != nil {
			return errors.Wrapf(err, "could not read arbeitstag %d", id)
		}
		json.NewEncoder(w).Encode(arbeitstag)

	default:
		http.Error(w, "only get or put", 400)
	}

	return nil
}
