package main

import (
	"fmt"

	"net/http"
	"os"

	"github.com/Norrun/feedmixer/internal/data"
	feedui "github.com/Norrun/feedmixer/internal/feed"

	"github.com/Norrun/feedmixer/internal/serverutils"
)

func main() {

	//initialSetup()

	//feedmixer.SetEnv(false)
	//fileSys, _ = feedmixer.GetFileSys()

	state := data.Load()

	mux := routing(state)

	server := http.Server{
		Addr:    ":8000",
		Handler: RenderError(nil, mux),
	}
	server.ListenAndServe()
}

func routing(state *data.ServerState) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /", nil)

	mux.Handle("GET /hx/search", nil)
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
