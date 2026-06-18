package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/Norrun/feedmixer/components"
	"github.com/Norrun/feedmixer/internal/datautils"
	feedui "github.com/Norrun/feedmixer/internal/feed"
	"github.com/Norrun/feedmixer/internal/keyring"
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

	var cachMap sync.Map

	mux := http.NewServeMux()

	mux.Handle("GET /", injectCacheMiddleware(&cachMap, http.HandlerFunc(mainPageHandler)))

	mux.Handle("POST /hx/search", injectCacheMiddleware(&cachMap, http.HandlerFunc(searchHandler)))

	server := http.Server{
		Addr:    ":8000",
		Handler: RenderError(mux, nil),
	}
	server.ListenAndServe()
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
	cacher := r.Context().Value(datautils.MakeTypedKey[*sync.Map](keyring.Cache))
	if cacher != nil {
		c, ok := cacher.(*sync.Map)
		if !ok {
			panic("missused type key")
		}
		c.Store(datautils.MakeTypedKey[[]feedui.DisplayItem](keyring.Cache), items)

	}

	comp := components.Posts(items)

	if err := comp.Render(r.Context(), w); err != nil {
		log.Print(err)
	}

}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("search"))
	var results []feedui.DisplayItem
	tKey := datautils.MakeTypedKey[*sync.Map](keyring.Cache)
	that := r.Context().Value(tKey)
	if that == nil {
		//make some error responce
		return
	}
	cache, ok := that.(*sync.Map)
	if !ok {
		panic("type key misuse")
	}

	that, ok = cache.Load(datautils.MakeTypedKey[[]feedui.DisplayItem](keyring.Cache))
	if !ok {
		//make some error responce
		return
	}
	items, ok := that.([]feedui.DisplayItem)
	if !ok {
		panic("type key misuse")
	}
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Title), query) {
			results = append(results, item)
		}
	}
	comp := components.ListItems(results)
	if aw, ok := w.(*wire.ApproveResponseWriter); ok {
		aw.Approve()
	}

	_ = comp.Render(r.Context(), w)
}

func injectCacheMiddleware(cache *sync.Map, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), datautils.MakeTypedKey[*sync.Map](keyring.Cache), cache)
		nr := r.WithContext(ctx)
		next.ServeHTTP(w, nr)
	})
}
