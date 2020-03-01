package kalender

import (
	"fmt"
	"time"
)

//     datum    |          feiertag          | kw_jahr | kw_nr | kw_tag | jahrtag | ordinal | monatsname | tagesname   | kalendermonat_id | kalenderwoche_id

type Kalendertag struct {
	ID      int
	Datum   time.Time `json:"date,omitempty"`
	Jahr    int16     `db:"jahr_id" json:"year,omitempty"`
	Monat   uint8     `db:"monat" json:"month,omitempty"`
	Tag     uint8     `db:"tag" json:"day,omitempty"`
	Jahrtag uint16    `db:"jahrtag" json:"jahrtag,omitempty"`
	KwJahr  int16     `db:"kw_jahr" json:"kw_jahr,omitempty"`
	KwNr    uint8     `db:"kw_nr" json:"kw_nr,omitempty"`
	KwTag   uint8     `db:"kw_tag" json:"kw_tag,omitempty"`
	// Ordinal  int       `json:"ord,omitempty"`
	Feiertag string `json:"feiertag,omitempty"`
}

// Gauss algorithm to calculate the date of Easter in a given year
// returns day, month, year as integers
func getEaster2(year int) (int, int, int) {
	// don"t go below start of Gregorian calendar
	if year < 1583 {
		year = 1583
	}
	// for type (by inference) and value assignment use :=
	// shorthand for   var month int = 3
	month := 3
	// determine the Golden number
	golden := (year % 19) + 1
	// determine the century number
	century := year/100 + 1
	// correct for the years that are not leap years
	xx := (3*century)/4 - 12
	// moon correction
	yy := (8*century+5)/25 - 5
	// find Sunday
	zz := (5*year)/4 - xx - 10
	// determine epact
	// age of moon on January 1st of that year
	// (follows a cycle of 19 years)
	ee := (11*golden + 20 + yy - xx) % 30
	if ee == 24 {
		ee += 1
	}
	if (ee == 25) && (golden > 11) {
		ee += 1
	}
	// get the full moon
	moon := 44 - ee
	if moon < 21 {
		moon += 30
	}
	// up to Sunday
	day := (moon + 7) - ((zz + moon) % 7)
	// possibly up a month in easter_date
	if day > 31 {
		day -= 31
		month = 4
	}
	return day, month, year
}

func GausOstern(jahr int) int {

	säkularzahl := jahr / 100
	fmt.Println("säkularzahl", säkularzahl)

	// die säkulare Mondschaltung	M(K) = 15 + (3K + 3) div 4 − (8K + 13) div 25
	säkMond := 15 + (3*säkularzahl+3)/4 - (8*säkularzahl+13)/25
	fmt.Println("säkMond", säkMond)

	// die säkulare Sonnenschaltung	S(K) = 2 − (3K + 3) div 4
	säkSonne := 2 - (3*säkularzahl + 3) // 4
	fmt.Println("säkSonne", säkSonne)

	// den Mondparameter	A(X) = X mod 19
	mondparameter := jahr % 19
	fmt.Println("mondparameter", mondparameter)

	// den Keim für den ersten Vollmond im Frühling	D(A,M) = (19A + M) mod 30
	keim := (19*mondparameter + säkMond) % 30
	fmt.Println("keim", keim)

	// die kalendarische Korrekturgröße	R(D,A) = (D + A div 11) div 29
	kalKorr := keim/29 + (keim/28-keim/29)*(mondparameter/11)
	fmt.Println("kalKorr", kalKorr, keim/29, keim/28-keim/29, mondparameter/11)

	// die Ostergrenze	OG(D,R) = 21 + D − R
	ostergrenze := 21 + keim - kalKorr
	fmt.Println("ostergrenze", ostergrenze)

	// den ersten Sonntag im März	SZ(X,S) = 7 − (X + X div 4 + S) mod 7
	sonntag := 7 - (jahr+jahr/4+säkSonne)%7
	fmt.Println("sonntag", sonntag)

	// die Entfernung des Ostersonntags von der Ostergrenze (Osterentfernung in Tagen)	OE(OG,SZ) = 7 − (OG − SZ) mod 7
	entf := 7 - (ostergrenze-sonntag)%7
	fmt.Println("entf", entf)

	// das Datum des Ostersonntags als Märzdatum (32. März = 1. April usw.)	OS = OG + OE
	märzdatum := ostergrenze + entf
	fmt.Println("märzdatum", märzdatum)

	var month, day int
	if märzdatum > 31 {
		month = 4
		day = märzdatum - 31
	} else {
		month = 3
		day = märzdatum
	}
	fmt.Println("Osterdatum", jahr, month, day)
	return märzdatum

}
