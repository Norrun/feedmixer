package feed

import (
	"github.com/Norrun/feedmixer/internal/shouldhave"
	"github.com/mmcdole/gofeed"
)

//import "github.com/mmcdole/gofeed"

func TestArea(url string) (*gofeed.Feed, error) {
	parser := gofeed.NewParser()

	return parser.ParseURL(url)

}

type DisplayItem struct {
	Title       string
	Description string
	Url         string
	Img         gofeed.Image
	Authou      gofeed.Person
}

type FeedHeader struct {
	shouldhave.Unimplemented
}
