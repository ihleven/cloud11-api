package repository

import (
	"database/sql"

	"github.com/ihleven/cloud11-api/arbeit"

	"github.com/pkg/errors"
)

func (r Repository) RetrieveArbeitsjahr(year, accountID int) (*arbeit.Arbeitsjahr, error) {
	// tage_freizeitausgleich, tage_krank, tage_arbeit, tage_buero, tage_dienstreise, tage_homeoffice, tage_frei,
	// jahr_id, user_id
	query := `
		SELECT id, urlaub_vorjahr, urlaub_anspruch, urlaub_tage, urlaub_geplant, urlaub_rest,
			   soll, ist, differenz
		  FROM arbeitsjahr 
		 WHERE id = $1
	`
	id := year*1000 + accountID

	a := arbeit.Arbeitsjahr{}
	err := r.DB.QueryRow(query, id).Scan(
		&a.ID,
		&a.UrlaubVorjahr,
		&a.UrlaubAnspruch,
		&a.UrlaubTage,
		&a.UrlaubGeplant,
		&a.UrlaubRest,
		&a.Soll,
		&a.Ist,
		&a.Differenz,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// special case: there was no row
		} else {
			return nil, errors.Wrapf(err, "Could not QueryRow and Scan for Arbeitstag %v", id)
		}
	}
	return &a, nil
}
