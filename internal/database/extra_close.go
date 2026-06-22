package database

import (
	"database/sql"
	"errors"
)

func (receiver *Queries) Close() error {
	if db, ok := receiver.db.(*sql.DB); ok {
		return db.Close()
	}
	return errors.New("nothing to close")
}
