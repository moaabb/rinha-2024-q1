package handlers

import (
	"database/sql"
	"log"
)

type RinhaHandler struct {
	logger *log.Logger
	db     *sql.DB
}

func NewRinhaHandler(l *log.Logger, db *sql.DB) *RinhaHandler {
	return &RinhaHandler{
		logger: l,
		db:     db,
	}
}
