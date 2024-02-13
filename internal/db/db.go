package db

import (
	"database/sql"
	"errors"
	"fmt"
)

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to connect to db: %s", err.Error()))
	}

	err = db.Ping()

	return db, err
}
