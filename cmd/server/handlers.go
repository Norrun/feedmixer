package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/Norrun/feedmixer/components"
	"github.com/Norrun/feedmixer/internal/data"
	"github.com/Norrun/feedmixer/internal/database"
	"github.com/Norrun/feedmixer/internal/datautils"
	"github.com/Norrun/feedmixer/internal/display"
	"github.com/Norrun/feedmixer/internal/wire"
	"github.com/a-h/templ"
)

type StandardHandlers struct {
	data.ServerState
	Cache  *datautils.SimpleCace
	Stacks struct {
		Central []templ.Component
		Tool    []templ.Component
	}
}

func (receiver *StandardHandlers) hxCentralView(w http.ResponseWriter, r *http.Request) {
	switch r.PathValue("going") {
	case "back":
		if len(receiver.Stacks.Central) == 0 {
			break
		}
		Back(receiver.Stacks.Central).Render(r.Context(), w)
		return

	default:

	}
	components.PostFeed(nil).Render(r.Context(), w)

}

func Back(comps []templ.Component) templ.Component {
	if len(comps) == 0 {
		panic("No page")
	}
	if len(comps) == 1 {
		return datautils.Peek(comps)
	} else {
		return datautils.Pop(comps)
	}
}

func (receiver StandardHandlers) mainPageHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := receiver.Data.DB.GetTagForest(r.Context())
	if err != nil {
		log.Panic(err, "tag forest deal with this later")
	}
	dbfeeds, err := receiver.Data.DB.GetAllFeeds(r.Context())
	if err != nil {
		log.Panic(err, "all feeds deal with this later")
	}
	feeds := datautils.ConvertSlice(dbfeeds, func(f database.Feed) display.Feed { return display.Feed{Title: f.Name, Id: int(f.ID)} })
	dbitems, err := receiver.Data.DB.GetAllItems(r.Context())
	if err != nil {
		log.Panic(err, "all items deal with this later")
	}
	items := datautils.ConvertSlice(dbitems, func(f database.Item) display.Item {
		return display.Item{Title: f.Title, Url: f.Url, Description: f.Description.String}
	})

	newVar := display.CentralData{Tags: tags, Feeds: feeds, Items: items}
	home := components.HomePage(newVar)

	if err := home.Render(r.Context(), w); err != nil {
		log.Print(err)
	}

	receiver.Stacks.Central = append(receiver.Stacks.Central, components.PostFeed(nil))

}

func (receiver StandardHandlers) hxEnableAddFeed(w http.ResponseWriter, r *http.Request) {

	form := components.AddingFeed("Joe Rogen", "https://feeds.megaphone.fm/GLT1412515089", "podcast, poular")
	button := components.AddFeedButton()

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

	feed, err := receiver.Data.DB.AddFeed(r.Context(), database.AddFeedParams{Name: r.FormValue("name"), Url: r.FormValue("url")})
	if err != nil {
		form := components.AddingFeed(r.FormValue("name"), r.FormValue("url"), r.FormValue("tags"))
		components.ErrorWithComponent("someting went wrong when saving your submission, may already exist", form).Render(r.Context(), w)
		return
	}
	log.Print(r.FormValue("name"), r.FormValue("url"), r.FormValue("tags"))
	for _, v := range strings.Split(r.FormValue("tags"), ",") {
		tag := strings.TrimSpace(v)
		dbtag, err := receiver.Data.DB.AddTag(r.Context(), tag)
		if err != nil {
			panic("deal with it later")
		}
		_, err = receiver.Data.DB.AttachTag(r.Context(), database.AttachTagParams{FeedID: feed.ID, TagID: dbtag.ID})
		if err != nil {
			panic("deal with it later")
		}
	}

	err = button.Render(r.Context(), w)
	if err != nil {
		log.Print(err)
	}

}

//...

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	const url = "https://www.youtube.com/feeds/videos.xml?channel_id=UCbRP3c757lWg9M-U7TyEkXA"
	var items []display.Item
	if aw, ok := w.(*wire.ApproveResponseWriter); ok {
		aw.Approve()
	}

	comp := components.Posts(items)

	if err := comp.Render(r.Context(), w); err != nil {
		log.Print(err)
	}

}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	_ = strings.ToLower(r.URL.Query().Get("search"))

}
