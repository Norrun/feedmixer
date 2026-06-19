package feedmixer

import (
	"embed"
	"io/fs"
)

var prod bool
var envSet bool

//go:embed sql/schema/*.sql
var sqldata embed.FS

func GetFileSys() fs.FS {

	return sqldata

}
