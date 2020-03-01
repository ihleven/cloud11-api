package arbeit

import (
	"fmt"
	"time"

	"github.com/ihleven/cloud11-api/kalender"
	"github.com/ihleven/cloud11-api/pkg/errors"
)

func NewUsecase(repo Repository) *Usecase {
	return &Usecase{repo}
}

type Usecase struct {
	repo Repository
}

func (uc *Usecase) ListArbeitsjahre(account int) (Arbeitsjahre, error) {

	arbeitsjahre, err := uc.repo.RetrieveArbeitsjahre(account, 0)
	if err != nil {
		return nil, err
	}
	return arbeitsjahre, nil
}

func (uc *Usecase) Arbeitsjahr(year, account int) (*Arbeitsjahr, error) {

	arbeitsjahre, err := uc.repo.RetrieveArbeitsjahre(account, year)
	if err != nil {
		return nil, err
	}
	if len(arbeitsjahre) == 0 {
		return nil, errors.Code(errors.NotFound, "Arbeitsjahr %d for account %d not found", year, account)
	}
	arbeitsjahr := arbeitsjahre[0]

	arbeitsjahr.Monate, err = uc.repo.SelectArbeitsmonate(year, 0, account)
	if err != nil {
		return nil, err
	}
	arbeitsjahr.Urlaube, err = uc.repo.ListUrlaube(account, year, 0)
	if err != nil {
		return nil, err
	}
	return &arbeitsjahr, nil
}

func (uc *Usecase) SetupArbeitsjahr(year, account int) (*Arbeitsjahr, error) {

	err := uc.repo.SetupArbeitsjahr(account, "IC", year)
	if err != nil {
		fmt.Println(err)
		// return nil, err
	}
	for month := 1; month <= 12; month++ {
		err := uc.repo.SetupArbeitsmonat(account, "IC", year, month)
		if err != nil {
			fmt.Println(err)
			// return nil, err
		}
	}
	for _, k := range kalender.ListKalendertage(year) {
		fmt.Println("Tag:", k)
		err := uc.repo.UpsertKalendertag(k)
		if err != nil {
			return nil, err
		}
		arbeitstag := Arbeitstag{
			Account: 1, Job: "IC", Datum2: Date(k.Datum), Jahr: k.Jahr, Monat: k.Monat, Status: "A", Kategorie: "-", Soll: 8,
			Kommentar: "testkommentar",
		}
		if k.KwTag > 5 {
			arbeitstag.Status = "-"
			arbeitstag.Soll = 0
		}
		if k.Feiertag != "" {
			arbeitstag.Status = "F"
		}
		err = uc.repo.UpsertArbeitstag(1, "IC", k.Datum, &arbeitstag)
		if err != nil {
			fmt.Println(err)
		}
	}

	arbeitsjahre, err := uc.repo.RetrieveArbeitsjahre(account, year)
	if err != nil {
		return nil, err
	}
	if len(arbeitsjahre) == 0 {
		return nil, errors.Code(errors.NotFound, "Arbeitsjahr %d for account %d not found", year, account)
	}
	arbeitsjahr := arbeitsjahre[0]

	arbeitsjahr.Monate, err = uc.repo.SelectArbeitsmonate(year, 0, account)
	return &arbeitsjahr, nil
}

func (uc *Usecase) GetArbeitsmonat(year int, month int, accountID int) (*Arbeitsmonat, error) {

	am, err := uc.repo.SelectArbeitsmonate(year, month, 1)

	if len(am) == 0 {
		return nil, nil
	}
	m := am[0]
	m.Arbeitstage, err = uc.repo.ListArbeitstage(year, month, 0, accountID)
	if err != nil {
		fmt.Println("err:", err)
		return nil, err
	}
	//return &Arbeitsmonat{m}, nil

	return &m, err
}

/// ARBEITSTAG ///
func (uc *Usecase) GetArbeitstag(year, month, day int, accountID int) (*Arbeitstag, error) {

	at, err := uc.repo.ReadArbeitstag(accountID, year, month, day)
	if err != nil {
		//return &Arbeitstag{}, nil
		return nil, errors.Wrap(err, "Could not read Arbeitstag: %d/%d/%d, %d", year, month, day, accountID)
	}
	return at, nil
}

func (uc *Usecase) UpdateArbeitstag(year, month, day int, accountID int, arbeitstag *Arbeitstag) (*Arbeitstag, error) {

	fmt.Println("usecase update arbeitstag", year, month, day, accountID, arbeitstag.Zeitspannen)

	id := ((year*100+month)*100+day)*1000 + accountID
	datum := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	arbeitstag.Pausen, _ = uc.UpdateZeitspannen(accountID, Date(datum), arbeitstag.Zeitspannen)
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

	err := uc.repo.UpsertArbeitstag(1, "IC", datum, arbeitstag)
	if err != nil {
		return nil, errors.Wrap(err, "Could not update Arbeitstag %d", id)
	}
	fmt.Println("sucess update arbeitstag", id)
	return nil, nil
}

func (uc *Usecase) UpdateZeitspannen(account int, datum Date, zeitspannen []Zeitspanne) (float64, error) {

	pausen := 0.0

	// Zeitspannen in der DB loeschen, deren Nr. es nicht mehr gibt
	dbZeitspannen, err := uc.repo.ListZeitspannen(account, datum)
	if err != nil {
		return 0.0, err
	}
	for _, dbZeitspanne := range dbZeitspannen {
		if !IsContained(zeitspannen, dbZeitspanne) {
			uc.repo.DeleteZeitspanne(account, datum, dbZeitspanne.Nr)
		}
	}
	// Insert oder Update Zeitspannen
	for _, zeitspanne := range zeitspannen {
		dauer := zeitspanne.Ende.Sub(*zeitspanne.Start).Hours()

		zeitspanne.Dauer = dauer
		pausen += dauer
		fmt.Println("Dauer: ", dauer, zeitspanne)
		err := uc.repo.UpsertZeitspanne(account, datum, &zeitspanne)
		if err != nil {
			fmt.Println("error upsert:", err)
			return 0.0, err
		}
	}
	return pausen, nil
}
