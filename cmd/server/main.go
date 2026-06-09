package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/Norrun/feedmixer/components"
	feedui "github.com/Norrun/feedmixer/internal/feed"
	"github.com/Norrun/feedmixer/internal/serverutils"
	"github.com/Norrun/feedmixer/internal/wire"
)

var fileSys fs.FS

func main() {

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

	//feedmixer.SetEnv(false)
	//fileSys, _ = feedmixer.GetFileSys()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", mainPageHandler)

	server := http.Server{
		Addr:    ":8000",
		Handler: RenderError(mux, nil),
	}
	server.ListenAndServe()
}

func RenderError(next http.Handler, renderer func(aw *wire.ApproveResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aw := wire.NewApprovalWriter(w, func(arw *wire.ApproveResponseWriter) bool {
			return arw.Status() < 400 && arw.Status() > 0
		})
		defer func() {
			r := recover()
			if r != nil {
				log.Printf("panic occured in some handler: %v", r)
				if aw.Approved() {
					return
				}
				io.WriteString(w, "<p>500 Internal error</p>")
				w.WriteHeader(500)
			}
		}()
		next.ServeHTTP(aw, r)
		if aw.Approved() {
			return
		}
		if renderer != nil {
			renderer(aw, r)
		}

	}
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	const url = "https://www.youtube.com/feeds/videos.xml?channel_id=UCbRP3c757lWg9M-U7TyEkXA"
	var items []feedui.DisplayItem
	feed, err := feedui.TestArea(url)
	if aw, ok := w.(*wire.ApproveResponseWriter); ok {
		aw.Approve()
	}
	if err != nil {
		io.WriteString(w, "<p>400 Bad request or something</p>")
		log.Print(err)
		w.WriteHeader(400)
	}

	for _, v := range feed.Items {

		items = append(items, feedui.DisplayItem{Title: v.Title,
			Description: v.Description,
			Url:         v.Link})
	}
	comp := components.ListItems(items)

	if err := comp.Render(r.Context(), w); err != nil {
		log.Print(err)
	}

}
