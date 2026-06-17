package feed

import (
	"github.com/mmcdole/gofeed"
)

func TestArea(url string) (FeedExt, error) {
	parser := gofeed.NewParser()
	original, err := parser.ParseURL(url)
	if err != nil {
		return FeedExt{}, err
	}
	extended := FeedExt{original}

	return extended, nil

}

type DisplayItem struct {
	Title       string
	Description string
	Url         string
	Img         gofeed.Image
	Author      gofeed.Person
}

type FeedHeader struct {
}

type FeedExt struct {
	*gofeed.Feed
}
type ItemExt struct {
}
