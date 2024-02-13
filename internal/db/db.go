package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(dsn string, logger *log.Logger) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to connect to db: %s", err.Error()))
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to db!")

	return db, nil
}
