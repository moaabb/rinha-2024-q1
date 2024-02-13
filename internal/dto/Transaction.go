package dto

import (
	"errors"

	"github.com/moaabb/rinha-de-backend-2024-q1/internal/models"
)

type TransactionRequest struct {
	Value       int64                  `json:"valor"`
	Type        models.TransactionType `json:"tipo"`
	Description string                 `json:"descricao"`
}

type TransactionResponse struct {
	Limit          int64 `json:"limite"`
	AccountBalance int64 `json:"saldo"`
}

func (t *TransactionRequest) Validate() error {
	if t.Value <= 0 {
		return InvalidDtoErr
	}

	if len(t.Description) > 10 || len(t.Description) == 0 {
		return InvalidDtoErr
	}

	if t.Type != models.Debit && t.Type != models.Credit {
		return InvalidDtoErr
	}

	return nil

}

var InvalidDtoErr = errors.New("invalid dto")
