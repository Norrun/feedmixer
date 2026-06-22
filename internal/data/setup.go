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
	"github.com/Norrun/feedmixer/internal/database"
	"github.com/Norrun/feedmixer/internal/serverutils"
	"github.com/pressly/goose/v3"
)

const (
	portableSubDir  = "/feedmixdata/"
	systemSubDir    = "/feedmixer/"
	DbFileName      = "app.db"
	AppInfoFileName = "appinfo.json"
)

func Setup() (ServerState, error) {
	dir, err := InteractiveSetup()
	if err != nil {
		return ServerState{}, fmt.Errorf("Interactive setup failed: %v", err)
	}

	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		return ServerState{}, fmt.Errorf("Error creating folder %s: %v", dir, err)
	}

	dbp := path.Join(dir, DbFileName)

	if _, err := os.Create(dbp); err != nil {
		return ServerState{}, fmt.Errorf("Error Creating %s: %v ", dbp, err)
	}
	db, err := sql.Open("sqlite", dbp)
	if err != nil {
		return ServerState{}, fmt.Errorf("Error openine %s as database: %v", dbp, err)
	}
	ver, err := Update(db)
	if err != nil {
		return ServerState{}, fmt.Errorf("Error when updating database: %v", err)
	}

	err = setAppInfo(dir, ver)
	if err != nil {
		return ServerState{}, fmt.Errorf("Error when creating appInfo: %v", err)
	}

	return NewServerState(database.New(db)), nil

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
		return fmt.Errorf("Error encoding AppInfo: %v", err)
	}

	p := path.Join(dir, AppInfoFileName)
	file, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("Error creating %s: %v", AppInfoFileName, err)
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("Error Writing to %s: %v", AppInfoFileName, err)
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
	goose.SetDialect("sqlite")
	if err := goose.Up(db, "sql/schema"); err != nil {
		return 0, fmt.Errorf("Failed to update database: %v", err)
	}
	v, err := goose.GetDBVersion(db)
	if err != nil {
		return 0, fmt.Errorf("Database version unknown: %v", err)
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
		return "", fmt.Errorf("Not Executable???: %v", err)
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
		return file, false, fmt.Errorf("Schrodinger's file: %s: %v", file, err)
	}

}
