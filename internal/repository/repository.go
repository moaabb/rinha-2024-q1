package repository

import (
	"github.com/moaabb/rinha-de-backend-2024-q1/internal/dto"
	"github.com/moaabb/rinha-de-backend-2024-q1/internal/models"
)

type TransactionRepository interface {
	GetAccountStatementByPartyId(id int64) (*models.Party, error)
	CreateTransaction(transaction dto.TransactionRequest, partyId int64) (*models.Party, error)
}
