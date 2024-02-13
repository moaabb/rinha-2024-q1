package transactiondb

import (
	"database/sql"
	"log"
)

type TransactionRepository struct {
	db     *sql.DB
	logger *log.Logger
}

func NewTransactionRepository(db *sql.DB, logger *log.Logger) *TransactionRepository {
	return &TransactionRepository{
		db, logger,
	}
}
