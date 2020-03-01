package arbeit

import (
	"time"

	"github.com/ihleven/cloud11-api/kalender"
)

type Repository interface {
	// RetrieveArbeitsjahr(year int, accountID int) (*Arbeitsjahr, error)
	RetrieveArbeitsmonat(year, month, accountID int) (*Arbeitsmonat, error)
	ListArbeitstage(year, month, week int, accountID int) ([]Arbeitstag, error)
	ReadArbeitstag(int, int, int, int) (*Arbeitstag, error)
	// UpdateArbeitstag(int, *Arbeitstag) error

	ListZeitspannen(account int, datum Date) ([]Zeitspanne, error)
	UpsertZeitspanne(account int, datum Date, z *Zeitspanne) error
	DeleteZeitspanne(account int, datum Date, nr int) error

	Close()

	SelectArbeitsmonate(year, month int, accountID int) ([]Arbeitsmonat, error)
	RetrieveArbeitsjahre(account, jahr int) ([]Arbeitsjahr, error)

	UpsertKalendertag(k kalender.Kalendertag) error
	// UpsertArbeitstag(a *Arbeitstag) error
	UpsertArbeitstag(account int, job string, datum time.Time, a *Arbeitstag) error
	SetupArbeitsjahr(account int, job string, year int) error
	SetupArbeitsmonat(account int, job string, year, month int) error

	ListUrlaube(account, jahr, nr int) ([]Urlaub, error)
}

// var Repo Repository
