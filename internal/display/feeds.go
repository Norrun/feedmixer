package display

import (
	"strconv"
	"strings"

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
	Title   string
	Id      int
	Checked bool
}
type Tag struct {
	Text    string
	Id      int
	Checked bool
	Related []Tag
}

type CentralData struct {
	Items []Item
	Feeds []Feed
	Tags  []Tag
}

func HTMLID(ctxprefix string, ids []int) string {
	var builder strings.Builder
	builder.WriteString(ctxprefix)
	for _, id := range ids {
		builder.WriteString("-")
		builder.WriteString(strconv.Itoa(id))
	}
	return builder.String()
}

func IDPathParam(ids []int) string {
	var builder strings.Builder
	for i, id := range ids {
		if i > 0 {
			builder.WriteRune('+')

		}
		builder.WriteString(strconv.Itoa(id))
	}
	return builder.String()
}
