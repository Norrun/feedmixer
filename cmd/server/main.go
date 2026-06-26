package main

import (
	"io"
	"log"

	"net/http"

	"github.com/Norrun/feedmixer/internal/data"
	"github.com/Norrun/feedmixer/internal/wire"

	_ "modernc.org/sqlite"
)

func main() {

	//initialSetup()

	//feedmixer.SetEnv(false)
	//fileSys, _ = feedmixer.GetFileSys()

	state, err := data.Load(true)
	if err != nil {
		log.Fatal(err)
	}

	mux := routing(state)

	server := http.Server{
		Addr: ":8000",
		Handler: RenderError(func(sw *wire.SnitchResponseWriter, r *http.Request) {
			io.WriteString(sw, "500 internal error")
		}, mux),
	}
	server.ListenAndServe()
}

func routing(state data.ServerState) *http.ServeMux {
	mux := http.NewServeMux()
	handlers := StandardHandlers{ServerState: state}

	mux.HandleFunc("GET /", handlers.mainPageHandler)

	mux.HandleFunc("POST /hx/add-feed", handlers.hxAddFeed)
	mux.HandleFunc("GET /hx/add-feed", handlers.hxEnableAddFeed)
	mux.HandleFunc("GET /hx/central/{going}", handlers.hxCentralView)
	mux.HandleFunc("GET /hx/get-items", handlers.hxCentralView)
	return mux
}
