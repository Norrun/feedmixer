package feedmixer

import (
	"embed"
	"io/fs"
)

//go:embed sql/schema/*.sql
var sqldata embed.FS

func GetFileSys() fs.FS {

	return sqldata

}
