package repository

import (
	"fmt"

	"github.com/ihleven/cloud11-api/arbeit"

	"github.com/pkg/errors"

	_ "github.com/lib/pq"
)

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
		err = errors.Wrapf(err, "Could not Select  arbeitstage %v", arbeitsmonat)
	}
	//fmt.Printf("monat: %v\n", arbeitsmonat)
	return nil, nil
}
