package arbeit

import (
	"database/sql"
	"time"
)

// Job beschreibt eine Entitaet, fuer die Arbeitszeit erfasst werden soll.
type Job struct {
	Account int
	Nr      int
	//name         string
	arbeitgeber string
	//arbeitnehmer string
	von time.Time
	bis time.Time
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

type Arbeitsmonat struct {
	Arbeitstage []Arbeitstag
}

type Arbeitswoche struct {
	Jahr int
	Nr   int
	//Job                   *Job
	Arbeitstage []Arbeitstag
}

type Arbeitstag struct {
	ID int `db:"id" json:"id"`

	Job     Job
	Datum   ` json:"tag,omitempty"`
	Account int //domain.Account

	// Typ: A-Arbeitstag, A2-Halber Arbeitstag, U-Urluab, F-Feiertag,
	Status ArbeitstagStatus `db:"status" json:"status,omitempty"`
	// B-Büro, D-Dienstreise, H-Homeoffice, K-Krankheit, U-Urlaub, Z-Zeitausgleich
	Kategorie ArbeitstagKategorie `db:"kategorie" json:"kategorie,omitempty"`
	//Krankheit string              `json:"krankmeldung,omitempty"`
	Urlaubstage float64 `json:"urlaubstage,omitempty"`

	Soll      float64    `json:"soll,omitempty"`
	Start     *time.Time `db:"beginn" json:"beginn"` //084300 => 8h43:00
	Ende      *time.Time `db:"ende" json:"ende"`     //173000 => 17h30:00
	Brutto    float64    `json:"brutto"`             //099700
	Pausen    float64    `json:"pausen"`
	Extra     float64    `json:"extra"`
	Netto     float64    `json:"netto"` // Brutto + Extra - Pausen
	Differenz float64    `json:"diff"`  // Netto - Soll

	// ergibt sich aus Saldo Vortag + Differenz
	Saldo     float64 `json:"saldo"`
	Kommentar string  `json:"kommentar"`
	// Saldo        *float64
	Zeitspannen []Zeitspanne ` json:"zeitspannen,omitempty"`
}

type Datum struct {
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
type Arbeitstag2 struct {
	Datum  Datum
	Typ    *string //
	Status *string //
	//Zeitspannen []Zeitspanne
}

type Zeitspanne struct {
	Nr          int                 `json:"nr"`
	Status      ZeitspanneKategorie `json:"status"`
	Start       *time.Time          `json:"start"`
	Ende        *time.Time          `json:"ende"`
	Dauer       float64             `json:"dauer"`
	Arbeitszeit float64             `json:"arbeitszeit"`
	Netto       float64             `json:"netto"`
	//Kategorie    string              // AZ, Pause, Weg
	// Titel, Story, Beschreibung, Grund *string
}
