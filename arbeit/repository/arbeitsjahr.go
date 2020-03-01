package repository

import (
	"database/sql"
	"fmt"

	"github.com/ihleven/cloud11-api/arbeit"

	"github.com/pkg/errors"
)

func (r Repository) SetupArbeitsjahr(account int, job string, year int) error {
	// tage_freizeitausgleich, tage_krank, tage_arbeit, tage_buero, tage_dienstreise, tage_homeoffice, tage_frei,
	// jahr_id, user_id
	stmt := `
		INSERT INTO c11_arbeitsjahr (account, job, jahr) 
		VALUES ($1, $2, $3)
	`
	n, err := r.DB.Exec(stmt, account, job, year)
	if err != nil {
		return errors.Wrap(err, "Could not insert c11_arbeitsjahr")
	}
	fmt.Println("repo setup arbeitsjahr", account, job, year, n)
	return nil
}

// func (r Repository) RetrieveArbeitsjahr(year int, accountID int) (*arbeit.Arbeitsjahr, error) {
// 	query := `
// 	`
// 	id := year*1000 + accountID

// 	a := arbeit.Arbeitsjahr{}
// 	err := r.DB.QueryRow(query, id).Scan(
// 		&a.Urlaub.Vorjahr,
// 		&a.Urlaub.Anspruch,
// 		&a.Urlaub.Tage,
// 		&a.Urlaub.Geplant,
// 		&a.Urlaub.Rest,
// 		&a.Soll,
// 		&a.Ist,
// 		&a.Differenz,
// 	)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			// special case: there was no row
// 		} else {
// 			return nil, errors.Wrapf(err, "Could not QueryRow and Scan for Arbeitstag %v", id)
// 		}
// 	}
// 	return &a, nil
// }

type Arbeitsjahr struct {
	ID             int
	UrlaubVorjahr  sql.NullFloat64
	UrlaubAnspruch sql.NullFloat64
	UrlaubTage     sql.NullFloat64
	UrlaubGeplant  sql.NullFloat64
	UrlaubRest     sql.NullFloat64
	Soll           sql.NullFloat64
	Ist            sql.NullFloat64
	Differenz      sql.NullFloat64
	// TageArbeit            sql.NullFloat64
	// TageKrank             sql.NullFloat64
	// tageFreizeitausgleich sql.NullFloat64
	// tageBuero             sql.NullFloat64
	// tageDienstreise       sql.NullFloat64
	// tageHomeoffice        sql.NullFloat64
	// tageFrei              sql.NullFloat64
	// jahrID                sql.NullInt64
	// userID                sql.NullInt64
	// Monate                []Arbeitsmonat
}

func (r Repository) RetrieveArbeitsjahre(account, jahr int) ([]arbeit.Arbeitsjahr, error) {

	arbeitsjahre := make([]arbeit.Arbeitsjahr, 0)

	query := `
		SELECT account, job, jahr, urlaub_vorjahr, urlaub_anspruch, urlaub_tage, urlaub_geplant, urlaub_rest, soll, ist, diff
		  FROM c11_arbeitsjahr 
		 WHERE account=$1 
	`
	if jahr != 0 {
		query += " AND jahr=$2"
	} else {
		query += " AND jahr>$2"
	}
	// id := year*1000 + accountID
	rows, err := r.DB.Query(query, account, jahr)
	if err != nil {
		return nil, errors.Wrap(err, "Could not query for rows")
	}
	defer rows.Close()

	a := arbeit.Arbeitsjahr{}

	for rows.Next() {
		err := rows.Scan(
			&a.Account,
			&a.Job,
			&a.Jahr,
			&a.Urlaub.Vorjahr,
			&a.Urlaub.Anspruch,
			&a.Urlaub.Tage,
			&a.Urlaub.Geplant,
			&a.Urlaub.Rest,
			&a.Soll,
			&a.Ist,
			&a.Diff,
		)
		if err != nil {
			return nil, errors.Wrap(err, "Could not scan for row")
		}
		// arbeitsjahr := arbeit.Arbeitsjahr{Soll: a.Soll.Float64, Ist: a.Ist.Float64, Differenz: a.Differenz.Float64}
		// arbeitsjahr.Urlaub = arbeit.Urlaub{a.UrlaubVorjahr.Float64, a.UrlaubAnspruch.Float64, a.UrlaubTage.Float64, a.UrlaubGeplant.Float64, a.UrlaubRest.Float64, 0}
		arbeitsjahre = append(arbeitsjahre, a)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "rows error")
	}

	return arbeitsjahre, nil
}
