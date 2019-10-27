package arbeit

import (
	"fmt"

	"github.com/pkg/errors"
)

func GetArbeitsjahr(year int, accountID int) (*Arbeitsjahr, error) {

	return Repo.RetrieveArbeitsjahr(year, accountID)
}

func GetArbeitsmonat(year int, month int, accountID int) (*Arbeitsmonat, error) {

	//m, _ := Repo.ListArbeitstage(year, month, 0, accountID)
	//return &Arbeitsmonat{m}, nil
	return Repo.RetrieveArbeitsmonat(year, month, 1)
}

func GetArbeitstag(year, month, day int, accountID int) (*Arbeitstag, error) {

	at, err := Repo.ReadArbeitstag(((year*100+month)*100+day)*1000 + accountID)
	if err != nil {
		//return &Arbeitstag{}, nil
		return nil, errors.Wrapf(err, "Could not read Arbeitstag: %d/%d/%d, %d", year, month, day, accountID)
	}
	return at, nil
}

func UpdateArbeitstag(year, month, day int, accountID int, arbeitstag *Arbeitstag) error {

	fmt.Println("usecase update arbeitstag", year, month, day, accountID, *arbeitstag)

	id := ((year*100+month)*100+day)*1000 + accountID

	arbeitstag.Pausen, _ = UpdateZeitspannen(id, arbeitstag.Zeitspannen)
	arbeitstag.Extra = 0

	// arbeitstagDB, err := Repo.ReadArbeitstag(id)
	// if err != nil {
	// 	fmt.Println("read at error:", err)
	// 	return errors.Wrapf(err, "Could not retrieve Arbeitstag %s%s%s", year, month, day)
	// }
	if arbeitstag.Start != nil && arbeitstag.Ende != nil {
		arbeitstag.Brutto = arbeitstag.Ende.Sub(*arbeitstag.Start).Hours()
		arbeitstag.Netto = arbeitstag.Brutto - arbeitstag.Pausen + arbeitstag.Extra
		arbeitstag.Differenz = arbeitstag.Soll - arbeitstag.Netto
	}

	err := Repo.UpdateArbeitstag(id, arbeitstag)
	if err != nil {
		return errors.Wrapf(err, "Could not update Arbeitstag %d", id)
	}
	fmt.Println("sucess update arbeitstag", id)
	return nil
}

func UpdateZeitspannen(arbeitstagId int, zeitspannen []Zeitspanne) (float64, error) {

	pausen := 0.0

	// Zeitspannen in der DB loeschen, deren Nr. es nicht mehr gibt
	dbZeitspannen, _ := Repo.ListZeitspannen(arbeitstagId)
	for _, dbZeitspanne := range dbZeitspannen {
		if !IsContained(zeitspannen, dbZeitspanne) {
			Repo.DeleteZeitspanne(arbeitstagId, dbZeitspanne.Nr)
		}
	}
	// Insert oder Update Zeitspannen
	for _, zeitspanne := range zeitspannen {
		dauer := zeitspanne.Ende.Sub(*zeitspanne.Start).Hours()

		zeitspanne.Dauer = dauer
		pausen += dauer
		fmt.Println("Dauer: ", dauer, zeitspanne)
		err := Repo.UpsertZeitspanne(arbeitstagId, &zeitspanne)
		if err != nil {
			fmt.Println("error upsert:", err)
			return 0.0, err
		}
	}
	return pausen, nil
}

func IsContained(haystack []Zeitspanne, needle Zeitspanne) bool {
	for _, n := range haystack {
		if n.Nr == needle.Nr {
			return true
		}
	}
	return false
}
