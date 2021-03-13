package main

import (
	"fmt"

	"github.com/ihleven/cloud11-api/arbeit/repository"
)

func main() {
	fmt.Println("Hallo Welt!")

	repo, _ := repository.NewPostgresRepository()
	//kalendertage, err := repo.MigrateKalendertage()
	err := repo.SetupJobs()
	err = repo.LegacyArbeitsjahre()
	err = repo.LegacyArbeitstage()

	// for _, a := range arbeitstage {
	// 	d := a.ID / 1000
	// 	jahr := d / 10000
	// 	monat := d % 10000 / 100
	// 	tag := d % 100
	// 	datum := time.Date(jahr, time.Month(monat), tag, 0, 0, 0, 0, time.UTC)
	// 	kwyear, kw := datum.ISOWeek()
	// 	feiertag := ""
	// 	var arbeitstag = arbeit.Arbeitstag{
	// 		Datum: arbeit.Datum{
	// 			Datum:    arbeit.Date(datum),
	// 			Jahr:     int16(jahr),
	// 			Monat:    uint8(monat),
	// 			Tag:      uint8(tag),
	// 			Jahrtag:  uint16(datum.YearDay()),
	// 			KwJahr:   int16(kwyear),
	// 			KwNr:     uint8(kw),
	// 			KwTag:    uint8(datum.Weekday()),
	// 			Feiertag: &feiertag,
	// 		},
	// 	}
	// 	fmt.Println(a.ID, a.Start, a.Ende, err, arbeitstag, d, jahr, monat, tag, datum)
	// }
	fmt.Println(err)
}
