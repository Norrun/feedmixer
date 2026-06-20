package data

import (
	"errors"
	"os"
	"path"
)

const (
	portableSubDir  = "/feedmixdata/"
	systemSubDir    = "/feedmixer/"
	DbFileName      = "app.db"
	AppInfoFileName = "appinfo.json"
)

func Setup() *ServerState {
	panic("unimplemented")
}

func Update(dir string) {
	panic("unimplemented")
}

func Scan() (string, []error, uint8) {
	var errs []error
	var portable, system bool
	var dir string

	dir, err := portablePath()
	if err != nil {
		errs = append(errs, err)
	}
	err = nil
	dir, portable, err = Exists(dir)
	if err != nil {
		errs = append(errs, err)

	}
	err = nil
	if portable {
		return dir, errs, 0
	}

	dir, err = systemPath()
	if err != nil {
		errs = append(errs, err)
	}
	err = nil
	dir, system, err = Exists(dir)
	if err != nil {
		errs = append(errs, err)
	}
	err = nil

	if system {
		return dir, errs, 1
	}
	if len(errs) > 0 {
		return "", errs, 3
	}

	return "", nil, 2

}

func portablePath() (string, error) {
	dir, err := os.Executable()
	if err != nil {
		return "", err
	}
	di, _ := path.Split(dir)

	dir = path.Join(di, portableSubDir)

	return dir, nil

}

func Exists(file string) (string, bool, error) {
	if _, err := os.Stat(file); err == nil {
		return file, true, nil

	} else if errors.Is(err, os.ErrNotExist) {
		return "", false, nil

	} else {
		return file, false, err
	}

}
