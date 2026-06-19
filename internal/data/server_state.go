package data

import (
	"database/sql"
	"sync"

	"github.com/Norrun/feedmixer/internal/datautils"
)

type ServerState struct {
	Data ServerData
}

type CacheKey int

const ()

type ServerData struct {
	DB    struct{} // Implemented later
	Cache struct {
		mu sync.RWMutex
		m  map[datautils.KeyTo[any, CacheKey]]any
	}
}

func Load(portable bool) *ServerState {
	dbstr := ""
	// should probably have a test tp see if portable instead.
	if portable {
		dbstr = "./app.db"
	} else {
		panic("unhandled case")
	}

	db, err := sql.Open("sqlite3", dbstr)
	if err != nil {
		// handle the goose stuff
	}

	_ = db

	return nil

}
