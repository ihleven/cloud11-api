package main

import (
	"github.com/ihleven/cloud11-api/arbeit"
	arbeithandler "github.com/ihleven/cloud11-api/arbeit/handler"
	"github.com/ihleven/cloud11-api/arbeit/repository"
	"github.com/ihleven/cloud11-api/webserver"
)

func main() {

	repo, _ := repository.NewPostgresRepository()
	uc := arbeit.NewUsecase(repo)
	uc.SetupArbeitsjahr(2020, 1)
	// for _, k := range kalender.ListKalendertage(2020) {
	// 	err := repo.UpsertKalendertag(k)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }

	// arbeit.Repo = *repo
	//tag, _ := repo.RetrieveArbeitstag(2019, 6, 23, 1)
	//	defer arbeit.Repo.Close()
	//fmt.Println(" => ", tag)

	srv := webserver.New("", 8000)
	srv.Register("/arbeit", arbeithandler.ArbeitHandler(uc))

	srv.ListenAndServe()
}
