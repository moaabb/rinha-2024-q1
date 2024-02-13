package handlers

import (
	"log"
	"net/http"

	"github.com/moaabb/rinha-de-backend-2024-q1/internal/repository"
)

type RinhaHandler struct {
	logger *log.Logger
	db     repository.TransactionRepository
}

func NewRinhaHandler(l *log.Logger, db repository.TransactionRepository) *RinhaHandler {
	return &RinhaHandler{
		logger: l,
		db:     db,
	}
}

func (h *RinhaHandler) GetAccountStatementByPartyId(w http.ResponseWriter, r *http.Request) {

}

func (h *RinhaHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {

}
