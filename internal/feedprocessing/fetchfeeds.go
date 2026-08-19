package feedprocessing

import (
	"runtime"

	"github.com/Norrun/feedmixer/internal/database"
	"github.com/Norrun/feedmixer/internal/datautils"
	"github.com/Norrun/feedmixer/internal/utils"
	"github.com/mmcdole/gofeed"
)

func FetchFeeds(done <-chan struct{}, urls []database.Feed, db *database.Queries) {
	threads := runtime.GOMAXPROCS(0)
	chans := utils.FanOut(done, utils.Unloader(len(urls), urls...), len(urls), threads*4, func(s database.Feed) datautils.Result[*gofeed.Feed] {
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
			//db.item
		}
	})

}
