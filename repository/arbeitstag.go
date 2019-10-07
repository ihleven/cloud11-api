package repository

import (
	"github.com/ihleven/cloud11-api/arbeit"

	"fmt"

	pq "github.com/lib/pq"
	"github.com/pkg/errors"
)

func (r Repository) ListArbeitstage(year, month, week int, accountID int) ([]arbeit.Arbeitstag, error) {

	query := `
		SELECT a.id, status, kategorie, krankmeldung, urlaubstage, soll, beginn, ende, brutto, pausen, extra, netto, differenz,
				k.jahr_id, k.monat, k.tag, k.datum, k.feiertag, k.kw_jahr, k.kw_nr , k.kw_tag, k.jahrtag, k.ordinal
		  FROM go_arbeitstag a, kalendertag k
		 WHERE a.tag_id=k.id
		   AND k.jahr_id=$1
	  ORDER BY a.id
	`
	aa := []arbeit.Arbeitstag{}
	err := r.DB.Select(&aa, query, year)
	fmt.Println(aa, year)
	return aa, err

}

func (r Repository) ReadArbeitstag(id int) (*arbeit.Arbeitstag, error) {

	query := `
		SELECT a.id, status, kategorie, krankmeldung, urlaubstage, soll, beginn, ende, brutto, pausen, extra, netto, differenz,
		  		k.jahr_id, k.monat, k.tag, k.datum, k.feiertag, k.kw_jahr, k.kw_nr , k.kw_tag, k.jahrtag, k.ordinal
		  FROM go_arbeitstag  a, kalendertag k
		 WHERE a.tag_id=k.id AND a.id = $1
	`

	a := arbeit.Arbeitstag{}
	err := r.DB.Get(&a, query, id)

	return &a, err
	// pausen_query := `
	// SELECT nr, von, bis
	//   FROM go_zeitspanne
	//  WHERE arbeitstag_id = $1
	// `
	// a.Zeitspannen = []arbeit.Zeitspanne{}
	// err = r.DB.Select(&a.Zeitspannen, pausen_query, id)
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "Could not Select  arbeits_zeitspanne %v", a.Zeitspannen)
	// }
	// fmt.Println(a)
	// return &a, nil
}

func (r Repository) ListZeitspannen(arbeitstag_id int) ([]arbeit.Zeitspanne, error) {
	query := `
		SELECT nr, von, bis
		  FROM arbeits_zeitspanne
		 WHERE arbeitstag_id=arbeitstag_id
	  ORDER BY a.id
	`
	zs := []arbeit.Zeitspanne{}
	err := r.DB.Select(&zs, query, arbeitstag_id)

	return zs, err
}

func (r Repository) UpsertZeitspanne(arbeitstagID int, z *arbeit.Zeitspanne) error {
	stmt := `
		INSERT INTO go_zeitspanne (nr,typ,von,bis,dauer,arbeitstag_id)
		                   VALUES ($1,$2,$3,$4,$5,$6)
	`
	_, err := r.DB.Exec(stmt, z.Nr, z.Typ, z.Von, z.Bis, z.Dauer, arbeitstagID)
	if err != nil {
		if pqErr := err.(*pq.Error); pqErr.Code != "23505" { //"23505": "unique_violation",
			return errors.Wrapf(err, "Could not insert go_zeitspanne %s", z.Nr)
		}
	}

	stmt = `
		UPDATE go_zeitspanne 
	   	   SET typ=$1,von=$2,bis=$3,dauer=$4
	 	 WHERE arbeitstag_id=$5 AND nr=$6
	`
	_, err = r.DB.Exec(stmt, z.Typ, z.Von, z.Bis, z.Dauer, arbeitstagID, z.Nr)
	if err != nil {
		return errors.Wrapf(err, "Could not update go_zeitspanne %s", z.Nr)
	}
	return nil
}
func (r Repository) DeleteZeitspanne(zeitspanne *arbeit.Zeitspanne) error {

	stmt := `DELETE FROM go_zeitspanne WHERE id=$1`
	_, err := r.DB.Exec(stmt, zeitspanne.Nr)
	if err != nil {
		return errors.Wrapf(err, "Could not delete go_zeitspanne %s", zeitspanne.Nr)
	}
	return nil
}

func (r Repository) UpdateArbeitstag(id int, a *arbeit.Arbeitstag) error {

	stmt := `
		UPDATE go_arbeitstag 
		   SET status=$1, kategorie=$2, krankmeldung=$3, urlaubstage=$4, soll=$5, beginn=$6, ende=$7, brutto=$8, netto=$9, differenz=$10
		 WHERE id = $11
	`
	res, err := r.DB.Exec(stmt, a.Status, a.Kategorie, a.Krankmeldung, a.Urlaubstage, a.Soll, a.Start, a.Ende, a.Brutto, a.Netto, id)
	if err != nil {
		return errors.Wrapf(err, "Could not exec sql update statement for id=%s", id)
	}
	_, err = res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "Could not get number of affected rows")
	}

	return nil
}
