package feedprocessing

import (
	"context"
	"database/sql"
	"runtime"
	"time"

	"github.com/Norrun/feedmixer/internal/database"
	"github.com/Norrun/feedmixer/internal/datautils"
	"github.com/Norrun/feedmixer/internal/display"
	"github.com/Norrun/feedmixer/internal/utils"
	"github.com/mmcdole/gofeed"
)

func FetchFeedsAndSave(feeds []database.Feed, db *database.Queries) ([]display.Item, []struct {
	ref int64
	err error
}) {

	var errec []struct {
		ref int64
		err error
	}
	var items []display.Item
	for _, dbf := range feeds {
		parser := gofeed.NewParser()
		feed, err := parser.ParseURL(dbf.Url)
		if err != nil {
			errec = append(errec, struct {
				ref int64
				err error
			}{dbf.ID, err})
		}
		for _, v := range feed.Items {

			//defer can()
			hasPublishTime := v.PublishedParsed != nil
			publishedNormalized := ""
			if hasPublishTime {
				publishedNormalized = v.PublishedParsed.Format(time.RFC3339)
			}
			db.AddItem(context.Background(), database.AddItemParams{
				Title:       v.Title,
				Url:         v.Link,
				Description: sql.NullString{String: v.Description, Valid: v.Description != ""},
				PublishedAt: sql.NullString{String: publishedNormalized, Valid: hasPublishTime},
				FeedID:      dbf.ID,
			})
			items = append(items, display.Item{Title: v.Title, Description: v.Description, Url: v.Link, Img: *v.Image, Authors: *v.Authors})

		}
	}
	return items, errec
}

func FetchFeeds(done <-chan struct{}, feeds []database.Feed) <-chan datautils.Result[*gofeed.Feed] {
	threads := runtime.GOMAXPROCS(0)
	//fetching feeds
	chans := utils.FanOut(done, utils.Unloader(len(feeds), feeds...), len(feeds), threads*4, func(s database.Feed) datautils.Result[*gofeed.Feed] {
		parser := gofeed.NewParser()
		feed, err := parser.ParseURL(s.Url)
		return datautils.NewResult(feed, err)
	})
	return utils.FanIn[datautils.Result[*gofeed.Feed]](done, chans...)

}
func SaveItemsFromFeeds() {
	/*
		utils.FanOut(done, , threads, ) any {
			if r.Err != nil {
				return r.Err
			}
			// add feed data
			for _, v := range r.Value.Items {
				ctx, can := context.WithCancel(context.Background())
				utils.ChainDoneContextCancel(done, can)
				//defer can()
				hasPublishTime := v.PublishedParsed != nil
				publishedNormalized := ""
				if hasPublishTime {
					publishedNormalized = v.PublishedParsed.Format(time.RFC3339)
				}
				db.AddItem(ctx, database.AddItemParams{
					Title:       v.Title,
					Url:         v.Link,
					Description: sql.NullString{String: v.Description, Valid: v.Description != ""},
					PublishedAt: sql.NullString{String: publishedNormalized, Valid: hasPublishTime},
					FeedID:      r.Value.db.ID,
				})

			}
			return nil
		})*/
}
