package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/moaabb/rinha-de-backend-2024-q1/internal/data"
	"github.com/moaabb/rinha-de-backend-2024-q1/internal/dto"
	"github.com/moaabb/rinha-de-backend-2024-q1/internal/models"
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
	id := r.PathValue("id")

	party_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		data.Response(w, http.StatusUnprocessableEntity, &data.H{
			"mensagem": "id não é um numero inteiro",
		})
		h.logger.Println("error: id is not an integer")
		return
	}

	party, err := h.db.GetAccountStatementByPartyId(party_id)
	if err != nil {
		ErrorHandlerDB(w, err)
		h.logger.Println("error: error fetching statement", err)
		return
	}

	var latestTransactions []dto.Transaction
	for _, v := range party.Transactions {
		transaction := dto.Transaction{
			Value:           v.Value,
			Type:            v.Type,
			Description:     v.Description,
			TransactionDate: v.CreatedAt.Format("2006-01-02T15:04:05.999999Z"),
		}

		latestTransactions = append(latestTransactions, transaction)
	}

	if len(latestTransactions) == 0 {
		latestTransactions = []dto.Transaction{}
	}

	response := dto.StatementResponse{
		Balance: dto.AccountBalance{
			Balance:       party.Balance,
			Limit:         party.Limit,
			StatementDate: time.Now().Format("2006-01-02T15:04:05.999999Z"),
		},
		LatestTransactions: latestTransactions,
	}

	data.Response(w, http.StatusOK, response)
}

func (h *RinhaHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	party_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		data.Response(w, http.StatusUnprocessableEntity, &data.H{
			"mensagem": "id não é um numero inteiro",
		})
		h.logger.Println("error: id is not an integer", err)
		return
	}

	var transaction dto.TransactionRequest
	err = json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		data.Response(w, http.StatusUnprocessableEntity, &data.H{
			"message": "error decoding body",
		})
		h.logger.Println("error: error decoding body", err)
		return
	}

	party, err := h.db.CreateTransaction(transaction, party_id)
	if err != nil {
		ErrorHandlerDB(w, err)
		h.logger.Println("error: error decoding body", err)
		return
	}

	response := dto.TransactionResponse{
		Limit:          party.Limit,
		AccountBalance: party.Balance,
	}

	data.Response(w, http.StatusOK, response)
}

func ErrorHandlerDB(w http.ResponseWriter, err error) {
	if err != nil {
		if err.(*models.ErrorDefinition).ErrorCode == "404" {
			data.Response(w, http.StatusNotFound, data.H{})
			return
		}

		data.Response(w, http.StatusUnprocessableEntity, data.H{})
	}
}
