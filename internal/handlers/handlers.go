package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
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

func (h *RinhaHandler) GetAccountStatementByPartyId(c *fiber.Ctx) error {
	id := c.Params("id")

	party_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		h.logger.Println("error: id is not an integer")
		return c.Status(http.StatusUnprocessableEntity).JSON(&data.H{
			"mensagem": "id não é um numero inteiro",
		})
	}

	party, err := h.db.GetAccountStatementByPartyId(party_id)
	if err != nil {
		h.logger.Println("error: error fetching statement", err)
		return ErrorHandlerDB(c, err)
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

	return c.Status(http.StatusOK).JSON(response)
}

func (h *RinhaHandler) CreateTransaction(c *fiber.Ctx) error {
	id := c.Params("id")

	party_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		h.logger.Println("error: id is not an integer", err)
		return c.Status(http.StatusUnprocessableEntity).JSON(&data.H{
			"mensagem": "id não é um numero inteiro",
		})
	}

	var transaction dto.TransactionRequest
	err = c.BodyParser(&transaction)
	if err != nil {
		h.logger.Println("error: error decoding body", err)
		return c.Status(http.StatusUnprocessableEntity).JSON(&data.H{
			"mensagem": "dto inválido",
		})
	}

	err = transaction.Validate()
	if err != nil {
		h.logger.Println("invalid dto err")
		return c.Status(http.StatusUnprocessableEntity).JSON(&data.H{
			"mensagem": "dto inválido",
		})
	}

	party, err := h.db.CreateTransaction(transaction, party_id)
	if err != nil {
		h.logger.Println("error: error creating transaction", err)
		return ErrorHandlerDB(c, err)
	}

	response := dto.TransactionResponse{
		Limit:          party.Limit,
		AccountBalance: party.Balance,
	}

	return c.Status(http.StatusOK).JSON(response)
}

func ErrorHandlerDB(c *fiber.Ctx, err error) error {
	if err != nil {
		if err.(*models.ErrorDefinition).ErrorCode == "404" {
			return c.SendStatus(http.StatusNotFound)
		}
		return c.Status(http.StatusUnprocessableEntity).JSON(data.H{})
	}

	return nil
}
