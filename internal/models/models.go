package models

import "time"

type TransactionType string

const (
	Credit TransactionType = "c"
	Debit  TransactionType = "d"
)

type Transaction struct {
	TransactionId int64
	Value         int64
	Type          TransactionType
	PartyId       int64
	Party         Party
	Description   string
	CreatedAt     time.Time
}

type Party struct {
	PartyId      int64
	Limit        int64
	Balance      int64
	Transactions []Transaction
}
