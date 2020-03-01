package repository

import (
	"database/sql"
	"fmt"

	"github.com/ihleven/cloud11-api/arbeit"
	"github.com/ihleven/cloud11-api/pkg/errors"

	_ "github.com/lib/pq"
)

type Arbeitsmonat struct {
	Soll      sql.NullFloat64
	Ist       sql.NullFloat64
	Differenz sql.NullFloat64
}

func (r Repository) SetupArbeitsmonat(account int, job string, year, month int) error {

	stmt := `
		INSERT INTO c11_arbeitsmonat (account, job, jahr, monat) 
		VALUES ($1, $2, $3, $4)
	`
	n, err := r.DB.Exec(stmt, account, job, year, month)
	if err != nil {
		return errors.Wrap(err, "Could not insert c11_arbeitsmonat")
	}
	fmt.Println("repo setup arbeitsmonat", account, job, year, n)
	return nil
}

func (r Repository) SelectArbeitsmonate(year, month int, accountID int) ([]arbeit.Arbeitsmonat, error) {

	arbeitsmonate := []arbeit.Arbeitsmonat{}

	query := `
		SELECT monat, soll, ist, diff
		  FROM c11_arbeitsmonat m 
		 WHERE m.account=$1 AND m.jahr=$2 
	`
	var err error
	if month != 0 {
		query += "AND m.monat=$3"
		err = r.DB.Select(&arbeitsmonate, query, accountID, year, month)
	} else {
		err = r.DB.Select(&arbeitsmonate, query, accountID, year)
	}
	if err != nil {
		return nil, errors.Wrap(err, "Could not Select  arbeitstage %v", arbeitsmonate)
	}

	// monate := make([]arbeit.Arbeitsmonat, len(arbeitsmonate))
	// for i, am := range arbeitsmonate {
	// 	if am.Soll.Valid {
	// 		monate[i].Soll = am.Soll.Float64
	// 	}
	// 	if am.Ist.Valid {
	// 		monate[i].Ist = am.Ist.Float64
	// 	}
	// 	if am.Differenz.Valid {
	// 		monate[i].Differenz = am.Differenz.Float64
	// 	}
	// }
	return arbeitsmonate, nil
}

func (r Repository) RetrieveArbeitsmonat(year, month int, accountID int) (*arbeit.Arbeitsmonat, error) {
	fmt.Println("RetrieveArbeitsmonat", year, month)

	arbeitsmonat := []arbeit.Arbeitstag{}

	query := `
		SELECT k.*, a.*
		  FROM kalendertag k LEFT OUTER JOIN go_arbeitstag a 
		    ON k.id=a.tag_id
		 WHERE k.jahr_id=$1 AND k.monat=$2
	`
	query = `
		SELECT a.id, status, kategorie, krankmeldung, urlaubstage, soll, beginn, ende, brutto, pausen, extra, netto, differenz,
				k.jahr_id, k.monat, k.tag, k.datum, k.feiertag, k.kw_jahr, k.kw_nr , k.kw_tag, k.jahrtag, k.ordinal
		  FROM go_arbeitstag a, kalendertag k
		 WHERE a.tag_id=k.id
		   AND k.jahr_id=$1 AND k.monat=$2
	  `
	err := r.DB.Select(&arbeitsmonat, query, year, month)

	if err != nil {
		err = errors.Wrap(err, "Could not Select  arbeitstage %v", arbeitsmonat)
	}
	//fmt.Printf("monat: %v\n", arbeitsmonat)
	return nil, nil
}

func (r Repository) RetrieveJahresArbeitsmonate(year, accountID int) ([]arbeit.Arbeitsmonat, error) {
	fmt.Println("RetrieveJahresArbeitsmonat", year)

	arbeitsmonat := []arbeit.Arbeitstag{}

	query := `
		SELECT k.*, a.*
		  FROM kalendertag k LEFT OUTER JOIN go_arbeitstag a 
		    ON k.id=a.tag_id
		 WHERE k.jahr_id=$1 AND k.monat=$2
	`
	query = `
		SELECT a.id, status, kategorie, krankmeldung, urlaubstage, soll, beginn, ende, brutto, pausen, extra, netto, differenz,
				k.jahr_id, k.monat, k.tag, k.datum, k.feiertag, k.kw_jahr, k.kw_nr , k.kw_tag, k.jahrtag, k.ordinal
		  FROM go_arbeitstag a, kalendertag k
		 WHERE a.tag_id=k.id
		   AND k.jahr_id=$1 AND k.monat=$2
	  `
	err := r.DB.Select(&arbeitsmonat, query, year)

	if err != nil {
		err = errors.Wrap(err, "Could not Select  arbeitstage %v", arbeitsmonat)
	}
	//fmt.Printf("monat: %v\n", arbeitsmonat)
	return nil, nil
}
