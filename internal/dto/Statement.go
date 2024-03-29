package dto

import (
	"github.com/moaabb/rinha-de-backend-2024-q1/internal/models"
)

type StatementResponse struct {
	Balance            AccountBalance `json:"saldo"`
	LatestTransactions []Transaction  `json:"ultimas_transacoes"`
}

type AccountBalance struct {
	Balance       int64  `json:"total"`
	StatementDate string `json:"data_extrato"`
	Limit         int64  `json:"limite"`
}

type Transaction struct {
	Value           int64                  `json:"valor"`
	Type            models.TransactionType `json:"tipo"`
	Description     string                 `json:"descricao"`
	TransactionDate string                 `json:"realizada_em"`
}
