package feedprocessing

import (
	"runtime"

	"github.com/Norrun/feedmixer/internal/database"
	"github.com/Norrun/feedmixer/internal/datautils"
	"github.com/Norrun/feedmixer/internal/utils"
	"github.com/mmcdole/gofeed"
)

func FetchFeeds(done <-chan struct{}, feeds []database.Feed) {
	threads := runtime.GOMAXPROCS(0)
	//fetching feeds
	chans := utils.FanOut(done, utils.Unloader(len(feeds), feeds...), len(feeds), threads*4, func(s database.Feed) datautils.Result[*gofeed.Feed] {
		parser := gofeed.NewParser()
		feed, err := parser.ParseURL(s.Url)
		return datautils.NewResult(struct {
			*gofeed.Feed
			db database.Feed
		}{feed, s}, err)
	})
	return utils.FanIn(done, chans)

}
