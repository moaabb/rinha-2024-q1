package models

type TransactionType rune

const (
	Credit TransactionType = 'c'
	Debit                  = 'd'
)

type Transaction struct {
	TransactionId int64
	Value         int64
	Type          TransactionType
	PartyId       int64
	Party         Party
	Description   string
	Date          string
}

type Party struct {
	PartyId int64
	Limit   int64
	Balance int64
}
