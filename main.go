package main

import (
	"log"

	"github.com/ihleven/cloud11-api/arbeit"
	"github.com/ihleven/cloud11-api/arbeit/repository"
	"github.com/ihleven/cloud11-api/http"
)

func main() {

	repo, _ := repository.NewPostgresRepository()
	arbeit.Repo = *repo
	defer arbeit.Repo.Close()
	//tag, _ := repo.RetrieveArbeitstag(2019, 6, 23, 1)
	//fmt.Println(" => ", tag)

	srv := http.NewServer(":8000")
	log.Fatal(srv.ListenAndServe())
}
