package data

import "time"

var registry map[string]time.Time

var nameReg map[string]string

type FeedInfo struct {
	url       string
	lastFetch time.Time
}

