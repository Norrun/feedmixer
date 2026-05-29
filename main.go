package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

type TutorialFilm struct {
	Title    string
	Director string
}

func main() {
	fmt.Println("starting service")
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", tutorialTemplateRender)
	mux.HandleFunc("POST /add-film/", CreateAddfilm)

	server := http.Server{
		Addr:    ":8000",
		Handler: mux,
	}
	server.ListenAndServe()
}

func tutorialTemplateRender(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/tutorial.html"))
	films := map[string][]TutorialFilm{
		"Films": {
			{Title: "The Godfather", Director: "Francis Ford Coppola"},
			{Title: "Blade Runner", Director: "Ridley Scott"},
			{Title: "The Thing", Director: "John Carpenter"},
		},
	}
	tmpl.Execute(w, films)
}

func CreateAddfilm(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	title := r.PostFormValue("title")
	director := r.PostFormValue("director")
	tmpl := template.Must(template.ParseFiles("templates/tutorial.html"))
	err := tmpl.ExecuteTemplate(w, "film-list-element", TutorialFilm{
		Title:    title,
		Director: director,
	})
	if err != nil {
		log.Print(err)
	}

}
