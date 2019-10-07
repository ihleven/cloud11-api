package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/ihleven/cloud11-api/arbeit"
	"github.com/ihleven/cloud11-api/repository"

	. "github.com/logrusorgru/aurora"
)

type User struct {
	Id      string
	Balance uint64
}

func main() {
	u := User{Id: "US123", Balance: 8}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	res, _ := http.Post("https://httpbin.org/post", "application/json; charset=utf-8", b)
	var body struct {
		// httpbin.org sends back key/value pairs, no map[string][]string
		Headers map[string]string `json:"headers"`
		Origin  string            `json:"origin"`
	}
	json.NewDecoder(res.Body).Decode(&body)
	fmt.Println(body)
	fmt.Println("Hello,", Magenta("Aurora"))
	fmt.Println(Bold(Cyan("Cya!")))

	arbeit.Repo, _ = repository.NewPostgresRepository()
	defer arbeit.Repo.Close()
	//tag, _ := repo.RetrieveArbeitstag(2019, 6, 23, 1)
	//fmt.Println(" => ", tag)

	http.HandleFunc("/arbeit/", arbeit.ArbeitHandler)
	log.Fatal(http.ListenAndServe(":3001", nil))
}

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
