package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/Norrun/feedmixer"
	"github.com/Norrun/feedmixer/internal/must"
)

var fileSys fs.FS

type TutorialFilm struct {
	Title    string
	Director string
}

func main() {
	fmt.Println("starting service")
	feedmixer.SetEnv(false)
	fileSys, _ = feedmixer.GetFileSys()
	must.StartPeriodicDump(30*time.Second, 10*time.Second)
	defer func() {
		err := must.GetResidualErrorsError()
		if err != nil {
			log.Print(err)
		}
	}()
	defer must.LogIgnoredErrorsOnPanic()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", renderMainPage)

	server := http.Server{
		Addr:    ":8000",
		Handler: mux,
	}
	server.ListenAndServe()
}

func renderMainPage(w http.ResponseWriter, r *http.Request) {

}

func CreateAddfilm(w http.ResponseWriter, r *http.Request) {

}
