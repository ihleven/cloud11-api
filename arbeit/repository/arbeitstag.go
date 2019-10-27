package repository

import (
	"database/sql"

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
		SELECT a.id, status, kategorie, urlaubstage, soll, beginn, ende, brutto, pausen, extra, netto, differenz, kommentar,
		  		k.jahr_id, k.monat, k.tag, k.datum, k.feiertag, k.kw_jahr, k.kw_nr , k.kw_tag, k.jahrtag, k.ordinal
		  FROM go_arbeitstag  a, kalendertag k
		 WHERE a.tag_id=k.id AND a.id = $1
	`

	a := arbeit.Arbeitstag{}
	err := r.DB.Get(&a, query, id)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not Select  arbeitstag %v", id)
	}
	fmt.Printf("Arbeitstag: %v\n", a)

	//return &a, err
	// pausenQuery := `
	// 	SELECT nr, typ, von as start, bis as ende, dauer
	//       FROM go_zeitspanne
	//      WHERE arbeitstag_id = $1
	// `
	a.Zeitspannen, err = r.ListZeitspannen(id)
	//err = r.DB.Select(&a.Zeitspannen, pausenQuery, id)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("erronorows bei zeitspannen")
		} else {
			return nil, errors.Wrapf(err, "Could not Select  arbeits_zeitspanne %v", a.Zeitspannen)
		}
	}
	fmt.Println("zeitspannen:", a)
	return &a, nil
}

func (r Repository) ListZeitspannen(arbeitstag_id int) ([]arbeit.Zeitspanne, error) {
	query := `
		SELECT nr, status, start, ende, dauer
		  FROM go_zeitspanne a
		 WHERE arbeitstag_id=$1
	  ORDER BY a.nr
	`
	zs := []arbeit.Zeitspanne{}
	err := r.DB.Select(&zs, query, arbeitstag_id)

	return zs, err
}

func (r Repository) UpsertZeitspanne(arbeitstagID int, z *arbeit.Zeitspanne) error {
	stmt := `
		INSERT INTO go_zeitspanne (nr,status,start,ende,dauer,arbeitstag_id)
		                   VALUES ($1,$2,$3,$4,$5,$6)
	`
	_, err := r.DB.Exec(stmt, z.Nr, z.Status, z.Start, z.Ende, z.Dauer, arbeitstagID)
	if err != nil {
		if pqErr := err.(*pq.Error); pqErr.Code != "23505" { //"23505": "unique_violation",
			return errors.Wrapf(err, "Could not insert go_zeitspanne %s", z.Nr)
		}
	}

	stmt = `
		UPDATE go_zeitspanne 
	   	   SET status=$1,start=$2,ende=$3,dauer=$4
	 	 WHERE arbeitstag_id=$5 AND nr=$6
	`
	fmt.Println("usert:", z.Status)
	_, err = r.DB.Exec(stmt, z.Status, z.Start, z.Ende, z.Dauer, arbeitstagID, z.Nr)
	if err != nil {
		return errors.Wrapf(err, "Could not update go_zeitspanne %s", z.Nr)
	}
	return nil
}
func (r Repository) DeleteZeitspanne(arbeitstag_id int, nr int) error {

	stmt := `DELETE FROM go_zeitspanne WHERE arbeitstag_id=$1 AND nr=$2`
	_, err := r.DB.Exec(stmt, arbeitstag_id, nr)
	if err != nil {
		return errors.Wrapf(err, "Could not delete go_zeitspanne %d.%d", arbeitstag_id, nr)
	}
	return nil
}

func (r Repository) UpdateArbeitstag(id int, a *arbeit.Arbeitstag) error {

	stmt := `
		INSERT INTO go_arbeitstag (id,status,kategorie,urlaubstage,soll,beginn,ende,brutto,pausen,extra,netto,differenz,kommentar,tag_id)
		                   VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
	`
	_, err := r.DB.Exec(stmt, id, (*a).Status, (*a).Kategorie, (*a).Urlaubstage,
		(*a).Soll, (*a).Start, (*a).Ende, (*a).Brutto, (*a).Pausen, (*a).Extra, (*a).Netto, (*a).Differenz, (*a).Kommentar,
		id/1000)
	if err != nil {
		if pqErr := err.(*pq.Error); pqErr.Code != "23505" { //"23505": "unique_violation",
			return errors.Wrapf(err, "Could not insert go_arbeitstag %s", id)
		}
	}
	fmt.Println("repo update arbeitstag", id, *a)

	stmt = `
		UPDATE go_arbeitstag 
		   SET status=$1, kategorie=$2, urlaubstage=$3, 
		       soll=$4, beginn=$5, ende=$6, brutto=$7, pausen=$8, extra=$9, netto=$10, differenz=$11, kommentar=$12
		 WHERE id = $13
	`
	res, err := r.DB.Exec(stmt, (*a).Status, (*a).Kategorie, (*a).Urlaubstage,
		(*a).Soll, (*a).Start, (*a).Ende, (*a).Brutto, (*a).Pausen, (*a).Extra, (*a).Netto, (*a).Differenz, (*a).Kommentar, id)
	if err != nil {
		return errors.Wrapf(err, "Could not exec sql update statement for id=%s", id)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "Could not get number of affected rows")
	}
	if affected == 0 {
		return errors.Wrap(err, "no affected rows")

	}

	return nil
}
