package dto

import (
	"time"

	"github.com/moaabb/rinha-de-backend-2024-q1/internal/models"
)

type StatementResponse struct {
	Balance          AccountBalance     `json:"saldo"`
	LastTransactions []LastTransactions `json:"ultimas_transacoes"`
}

type AccountBalance struct {
	Balance       int64     `json:"saldo"`
	StatementDate time.Time `json:"data_extrato"`
	Limit         int64     `json:"limite"`
}

type LastTransactions struct {
	Value           int64                  `json:"valor"`
	Type            models.TransactionType `json:"tipo"`
	Description     string                 `json:"descricao"`
	TransactionDate time.Time              `json:"realizada_em"`
}
