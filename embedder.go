package feedmixer

import (
	"embed"
	"errors"
	"io/fs"
	"os"

	"github.com/Norrun/feedmixer/internal/must"
)

var prod bool
var envSet bool

//go:embed templates/*
var buildFiles embed.FS

var fileSys fs.FS

func GetFileSys() (_ fs.FS, erret error) {
	defer func() {
		erret = must.MustHandle(erret)
	}()

	if prod {
		return buildFiles, nil
	}
	if fileSys != nil {
		return fileSys, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fileSys = os.DirFS(wd)

	return fileSys, nil
}

func SetEnv(production bool) error {
	if envSet {
		return errors.New("already configured")
	}
	prod = production
	return nil
}
