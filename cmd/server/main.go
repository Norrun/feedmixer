package main

import (
	"fmt"
	"io"
	"log"

	"net/http"
	"os"

	"github.com/Norrun/feedmixer/internal/data"
	feedui "github.com/Norrun/feedmixer/internal/feed"
	"github.com/Norrun/feedmixer/internal/wire"

	"github.com/Norrun/feedmixer/internal/serverutils"
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
		Handler: RenderError(func(aw *wire.ApproveResponseWriter, r *http.Request) {
			status := aw.Status()
			if 399 > status {
				status = 500
			}
			aw.Approve()
			io.WriteString(aw, fmt.Sprintf("<p>%d %s </p>", status, http.StatusText(status)))
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
	return mux
}

func initialSetup() {
	for {
		var start bool
		input := serverutils.GetInput()
		switch input[0] {
		case "feed":
			f, err := feedui.TestArea(input[1])
			fmt.Println(f, err)
		case "next":
			start = true
		case "exit":
			os.Exit(0)
		}
		if start {
			break
		}
	}
}
