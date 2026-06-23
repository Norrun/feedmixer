package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/Norrun/feedmixer/components"
	"github.com/Norrun/feedmixer/internal/data"
	"github.com/Norrun/feedmixer/internal/database"
	"github.com/Norrun/feedmixer/internal/datautils"
	feedui "github.com/Norrun/feedmixer/internal/feed"
	"github.com/Norrun/feedmixer/internal/keyring"
	"github.com/Norrun/feedmixer/internal/wire"
)

type StandardHandlers struct {
	ret *datautils.Pipe[func(aw *wire.ApproveResponseWriter, r *http.Request)]
	data.ServerState
}

func (receiver StandardHandlers) mainPageHandler(w http.ResponseWriter, r *http.Request) {
	home := components.HomePage()
	wire.Approve(w)
	if err := home.Render(r.Context(), w); err != nil {
		log.Print(err)
	}
}

func (receiver StandardHandlers) hxEnableAddFeed(w http.ResponseWriter, r *http.Request) {

	form := components.AddingFeed()
	button := components.AddFeedButton()
	wire.Approve(w)
	if r.URL.Query().Get("a") == "cancel" {
		err := button.Render(r.Context(), w)
		if err != nil {
			log.Print(err)
		}
		return
	}
	err := form.Render(r.Context(), w)
	if err != nil {
		log.Print(err)
	}

}

func (receiver StandardHandlers) hxAddFeed(w http.ResponseWriter, r *http.Request) {
	button := components.AddFeedButton()
	err := r.ParseForm()
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
		return
	}

	_, err = receiver.Data.DB.AddFeed(r.Context(), database.AddFeedParams{Name: r.FormValue("name"), Url: r.FormValue("url")})
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
		return
	}
	log.Print(r.FormValue("name"), r.FormValue("url"))
	wire.Approve(w)
	err = button.Render(r.Context(), w)
	if err != nil {
		log.Print(err)
	}

}

//...

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
	cacher := r.Context().Value(datautils.MakeKeyTo[*sync.Map](keyring.Cache))
	if cacher != nil {
		c, ok := cacher.(*sync.Map)
		if !ok {
			panic("missused type key")
		}
		c.Store(datautils.MakeKeyTo[[]feedui.DisplayItem](keyring.Cache), items)

	}

	comp := components.Posts(items)

	if err := comp.Render(r.Context(), w); err != nil {
		log.Print(err)
	}

}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("search"))
	var results []feedui.DisplayItem
	tKey := datautils.MakeKeyTo[*sync.Map](keyring.Cache)
	that := r.Context().Value(tKey)
	if that == nil {
		//make some error responce
		return
	}
	cache, ok := that.(*sync.Map)
	if !ok {
		panic("type key misuse")
	}

	that, ok = cache.Load(datautils.MakeKeyTo[[]feedui.DisplayItem](keyring.Cache))
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

	err := comp.Render(r.Context(), w)
	if err != nil {
		log.Print(err)
	}
}
