package arbeit

type Repository interface {
	RetrieveArbeitsjahr(year int, accountID int) (*Arbeitsjahr, error)

	ListArbeitstage(year, month, week int, accountID int) ([]Arbeitstag, error)
	ReadArbeitstag(int) (*Arbeitstag, error)
	UpdateArbeitstag(int, *Arbeitstag) error
	ListZeitspannen(int) ([]Zeitspanne, error)
	UpsertZeitspanne(int, *Zeitspanne) error
	DeleteZeitspanne(int, int) error
	Close()
}

var Repo Repository
