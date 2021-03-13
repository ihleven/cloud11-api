package main

import (
	"net/http"

	"github.com/ihleven/cloud11-api/arbeit"
	"github.com/ihleven/cloud11-api/arbeit/repository"
	"github.com/ihleven/pkg/web"
)

func main() {

	repo, _ := repository.NewPostgresRepository()
	_ = arbeit.NewUsecase(repo)
	// uc.SetupArbeitsjahr(2020, 1)
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

	// srv := webserver.New("", 8000)
	// srv.Register("/arbeit", arbeithandler.ArbeitHandler(uc))
	// srv.Register("/home", drive.DispatchHandler(&fs.Drive))
	// // srv.Register("/serve/home/", drive.DispatchRaw(fs.Drive))
	// srv.ListenAndServe()

	srv := web.NewServer(8001)
	srv.Register("/", Hello)

	srv.Run()

}

func Hello(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	w.Write([]byte("Hello Welt!"))
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, DELETE, GET, HEAD")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}
