package repository

import (
	"database/sql"
	"time"

	"github.com/ihleven/cloud11-api/arbeit"
	"github.com/ihleven/cloud11-api/pkg/errors"

	"fmt"

	pq "github.com/lib/pq"
)

var basequery = `
	SELECT a.account, a.job, a.datum, a.jahr, a.monat, a.status, a.kategorie, 
		   a.soll, a.start, a.ende, a.brutto, a.pausen, a.extra, a.netto, a.diff, a.kommentar,
		   d.jahr, d.monat, d.tag, kw_jahr, d.kw, d.kw_tag, d.jahrtag, d.feiertag
      FROM c11_arbeitstag a, c11_datum d
     WHERE a.datum=d.datum AND a.account=$1
`

func (r Repository) ListArbeitstage(year, month, week int, account int) ([]arbeit.Arbeitstag, error) {

	arbeitstage := []arbeit.Arbeitstag{}
	var err error
	query := basequery
	if month != 0 {
		query += " AND d.jahr=$2 AND d.monat=$3 ORDER BY d.datum"
		err = r.DB.Select(&arbeitstage, query, account, year, month)
	} else if week != 0 {
		query += `AND d.kw_jahr=$2 AND d.kw=$3 ORDER BY d.datum`
		err = r.DB.Select(&arbeitstage, query, account, year, week)
	} else {
		query += `AND d.jahr=$2 ORDER BY d.datum`
		err = r.DB.Select(&arbeitstage, query, account, year)
	}
	return arbeitstage, err
}

func (r Repository) ReadArbeitstag(account, year, month, day int) (*arbeit.Arbeitstag, error) {

	query := basequery + " AND d.jahr=$2 AND d.monat=$3 AND d.tag=$4"

	arbeitstag := arbeit.Arbeitstag{}
	err := r.DB.Get(&arbeitstag, query, account, year, month, day)
	if err == sql.ErrNoRows {
		return nil, errors.Code(errors.NotFound, "Arbeitstag %d %d %d %d not found ", account, year, month, day)
	}
	if err != nil {
		fmt.Println(err)
		return nil, errors.Wrap(err, "Could not Select  arbeitstag")
	}

	datum := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	arbeitstag.Zeitspannen, err = r.ListZeitspannen(account, arbeit.Date(datum))
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("erronorows bei zeitspannen")
		} else {
			return nil, errors.Wrap(err, "Could not Select  arbeits_zeitspanne %v", arbeitstag.Zeitspannen)
		}
	}
	fmt.Println("zeitspannen:", arbeitstag)
	return &arbeitstag, nil
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
			return errors.Wrap(err, "Could not insert go_arbeitstag %s", id)
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
		return errors.Wrap(err, "Could not exec sql update statement for id=%s", id)
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

func (r Repository) UpsertArbeitstag(account int, job string, datum time.Time, a *arbeit.Arbeitstag) error {
	// account | job | datum | jahr | monat | status | kategorie | soll | start | ende | brutto | pausen | extra | netto | diff | kommentar
	fmt.Println("repo update arbeitstag", *a, a.Datum.Jahr)

	stmt := `INSERT INTO c11_arbeitstag
		(account,job,datum,jahr,monat,status,kategorie,soll,start,ende,brutto,pausen,extra,netto,diff,kommentar)
		VALUES
		($1,     $2, $3,    $4, $5,   $6,    $7,       $8,   $9,   $10, $11,   $12,   $13,  $14, $15, $16)
	`
	_, err := r.DB.Exec(stmt, account, job, datum, a.Datum.Jahr, a.Datum.Monat, (*a).Status, (*a).Kategorie, //(*a).Urlaubstage,
		(*a).Soll, (*a).Start, (*a).Ende, (*a).Brutto, (*a).Pausen, (*a).Extra, (*a).Netto, (*a).Differenz, (*a).Kommentar)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code != "23505" { //"23505": "unique_violation",
			return errors.Wrap(err, "Could not insert go_arbeitstag %v", a)
		}
	}

	stmt = `
		UPDATE c11_arbeitstag
		   SET jahr=$4, monat=$5, status=$6, kategorie=$7, 
		       soll=$8, start=$9, ende=$10, brutto=$11, pausen=$12, extra=$13, netto=$14, diff=$15, kommentar=$16
		 WHERE account=$1 AND job=$2 AND datum=$3
	`
	res, err := r.DB.Exec(stmt, account, job, datum, a.Datum.Jahr, a.Datum.Monat, (*a).Status, (*a).Kategorie,
		(*a).Soll, (*a).Start, (*a).Ende, (*a).Brutto, (*a).Pausen, (*a).Extra, (*a).Netto, (*a).Differenz, (*a).Kommentar)
	if err != nil {
		return errors.Wrap(err, "Could not exec sql update statement for id=%s", a.Datum)
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
