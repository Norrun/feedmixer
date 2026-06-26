package display

import (
	"github.com/mmcdole/gofeed"
)

type Item struct {
	Title       string
	Description string
	Url         string
	Img         gofeed.Image
	Author      gofeed.Person
}

type Feed struct {
	Title string
	Id    string // Pre-processed int
}
type Tag struct {
	Text string
	Id   string // Pre-processed int
}

type CentralData struct {
	Items []Item
	Feeds []Feed
	Tags  []Tag
}
