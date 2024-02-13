package repository

type TransactionRepository interface {
	GetTransactionsByPartyId()
	CreateTransaction()
}
