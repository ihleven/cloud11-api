package arbeit

import (
	"database/sql"
	"time"
)

// Job beschreibt eine Entitaet, fuer die Arbeitszeit erfasst werden soll.
type Job struct {
	name         string
	arbeitgeber  string
	arbeitnehmer string
	von          time.Time
	bis          time.Time
}

type Kalendertag struct {
	//id    int
	Jahr     int16     `db:"jahr_id" json:"year,omitempty"`
	Monat    uint8     `db:"monat" json:"month,omitempty"`
	Tag      uint8     `db:"tag" json:"day,omitempty"`
	Datum    time.Time `json:"date,omitempty"`
	Feiertag *string   `json:"feiertag,omitempty"`
	KwJahr   int       `db:"kw_jahr" json:"kw_jahr,omitempty"`
	KwNr     uint8     `db:"kw_nr" json:"kw_nr,omitempty"`
	KwTag    uint8     `db:"kw_tag" json:"kw_tag,omitempty"`

	Jahrtag uint16 `db:"jahrtag" json:"jahrtag,omitempty"`
	Ordinal int    `json:"ord,omitempty"`
	//monatsname string
	//tagesname  string
}

//func (t Kalendertag) String() string {
//	return t.Jahr
//}
func (t Kalendertag) Gestern() {

}
func (t Kalendertag) Morgen() {

}

type Arbeitsjahr struct {
	ID int
	//Account               *domain.Account `json:"-"`
	//Job                   *Job
	UrlaubVorjahr         sql.NullFloat64
	UrlaubAnspruch        sql.NullFloat64
	UrlaubTage            sql.NullFloat64
	UrlaubGeplant         sql.NullFloat64
	UrlaubRest            sql.NullFloat64
	Soll                  sql.NullFloat64
	Ist                   sql.NullFloat64
	Differenz             sql.NullFloat64
	tageFreizeitausgleich sql.NullFloat64
	tageKrank             sql.NullFloat64
	tageArbeit            sql.NullFloat64
	tageBuero             sql.NullFloat64
	tageDienstreise       sql.NullFloat64
	tageHomeoffice        sql.NullFloat64
	tageFrei              sql.NullFloat64
	jahrID                sql.NullInt64
	userID                sql.NullInt64
}

type Zeitspanne struct {
	Nr                                int
	Typ                               string
	Von, Bis                          *time.Time
	Dauer                             *float64
	Titel, Story, Beschreibung, Grund *string
	Arbeitszeit                       bool
}

type Arbeitstag struct {
	ID int `db:"id" json:"id"`
	//Account      int //domain.Account
	//Job          Job
	Status       string     `db:"status" json:"status,omitempty"`
	Kategorie    string     `db:"kategorie" json:"kategorie,omitempty"`
	Krankmeldung bool       `json:"krankmeldung,omitempty"`
	Urlaubstage  float64    `json:"urlaubstage,omitempty"`
	Soll         float64    `json:"soll,omitempty"`
	Start        *time.Time `db:"beginn" json:"beginn,omitempty"`
	Ende         *time.Time `db:"ende" json:"ende,omitempty"`
	Brutto       float64    `json:"brutto,omitempty"`
	Pausen       float64    `json:"pausen,omitempty"`
	Extra        float64    `json:"extra,omitempty"`
	Netto        float64    `json:"netto,omitempty"`
	Differenz    float64    `json:"diff,omitempty"`
	// Saldo        *float64
	// Zeitspannen  []Zeitspanne
	Kalendertag ` json:"tag,omitempty"`
}
