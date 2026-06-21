package data

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Norrun/feedmixer"
	"github.com/Norrun/feedmixer/internal/serverutils"
	"github.com/pressly/goose/v3"
)

const (
	portableSubDir  = "/feedmixdata/"
	systemSubDir    = "/feedmixer/"
	DbFileName      = "app.db"
	AppInfoFileName = "appinfo.json"
)

func Setup() (*ServerState, error) {
	dir, err := InteractiveSetup()
	if err != nil {
		return nil, err
	}

	if err := os.Mkdir(dir, os.ModeDir); err != nil {
		return nil, err
	}

	dbp := path.Join(dir, DbFileName)

	if _, err := os.Create(dbp); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", dbp)
	if err != nil {
		return nil, err
	}
	ver, err := Update(db)
	if err != nil {
		return nil, err
	}

	err = setAppInfo(dir, ver)
	if err != nil {
		return nil, err
	}

	return NewServerState(db), nil

}

func setAppInfo(dir string, dbv int) error {
	app := appInfo{
		App: struct {
			Nr       [3]int "json:\"number\""
			Addendum string "json:\"addendum\""
		}{Nr: [3]int{0, 0, 0}, Addendum: ""},
		DBV: dbv,
	}
	data, err := json.Marshal(app)
	if err != nil {
		return err
	}

	p := path.Join(dir, AppInfoFileName)
	file, err := os.Create(p)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func InteractiveSetup() (string, error) {
	for {
		fmt.Println("Do you want portable setup? yes/no")
		input := serverutils.GetInput()
		if strings.Contains(strings.ToLower(input[0]), "y") {
			return portablePath()
		}
		if strings.Contains(strings.ToLower(input[0]), "n") {
			return systemPath()
		}
		fmt.Println("Invalid answer, retry")
	}

}

func Update(db *sql.DB) (int, error) {
	files := feedmixer.GetFileSys()
	goose.SetBaseFS(files)
	if err := goose.Up(db, "./"); err != nil {
		return 0, err
	}
	v, err := goose.GetDBVersion(db)
	if err != nil {
		return 0, err
	}
	return int(v), nil
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
