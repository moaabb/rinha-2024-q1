package dto

import "github.com/moaabb/rinha-de-backend-2024-q1/internal/models"

type TransactionRequest struct {
	Value       int64                  `json:"valor"`
	Type        models.TransactionType `json:"tipo"`
	Description string                 `json:"descricao"`
}

type TransactionResponse struct {
	Limit          int64 `json:"limite"`
	AccountBalance int64 `json:"saldo"`
}
