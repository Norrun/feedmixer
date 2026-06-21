package data

import (
	"database/sql"
	"errors"
	"path"
	"sync"

	"github.com/Norrun/feedmixer/internal/datautils"
)

type ServerState struct {
	Data ServerData
}

type DBR = *sql.DB

type CacheKey int

const ()

type appInfo struct {
	App struct {
		Nr       [3]int `json:"number"`
		Addendum string `json:"addendum"`
	} `json:"app"`
	DBV int `json:"db"`
}

type ServerData struct {
	DB    DBR // Implemented later
	Cache struct {
		mu sync.RWMutex
		m  map[datautils.KeyTo[any, CacheKey]]any
	}
}

func Load(portable bool) (*ServerState, error) {
	dir, errs, code := Scan()

	switch code {
	case 0, 1:
		break
	case 2:
		return Setup()

	case 3:
		return nil, errors.Join(errs...)
	}

	db, err := sql.Open("sqlite3", path.Join(dir, DbFileName))
	if err != nil {
		return nil, err
	}

	return NewServerState(db), nil

}

func CheckVersion(dir string) bool {
	panic("unimplemented")
}

func NewServerState(db DBR) *ServerState {
	return &ServerState{ServerData{
		DB: db,
		Cache: struct {
			mu sync.RWMutex
			m  map[datautils.KeyTo[any, CacheKey]]any
		}{sync.RWMutex{},
			map[datautils.KeyTo[any, CacheKey]]any{}}}}
}
