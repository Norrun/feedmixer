package feedprocessing

import (
	"context"
	"database/sql"
	"runtime"
	"time"

	"github.com/Norrun/feedmixer/internal/database"
	"github.com/Norrun/feedmixer/internal/datautils"
	"github.com/Norrun/feedmixer/internal/utils"
	"github.com/mmcdole/gofeed"
)

func FetchFeeds(done <-chan struct{}, feeds []database.Feed, db *database.Queries) {
	threads := runtime.GOMAXPROCS(0)
	chans := utils.FanOut(done, utils.Unloader(len(feeds), feeds...), len(feeds), threads*4, func(s database.Feed) datautils.Result[*gofeed.Feed] {
		parser := gofeed.NewParser()
		feed, err := parser.ParseURL(s.Url)
		return datautils.NewResult(feed, err)
	})
	utils.Fan(done, len(chans), threads, func(r datautils.Result[*gofeed.Feed]) any {
		if r.Err != nil {
			return r.Err
		}
		// add feed data
		for _, v := range r.Value.Items {
			ctx, can := context.WithCancel(context.Background())
			publishedNormalized := ""
			if v.PublishedParsed != nil {
				v.PublishedParsed.Format(time.RFC822)
			}
			db.AddItem(ctx, database.AddItemParams{
				Title:       v.Title,
				Url:         v.Link,
				Description: sql.NullString{String: v.Description, Valid: v.Description != ""},
				PublishedAt: v.PublishedParsed,
			})
		}
	})

}
