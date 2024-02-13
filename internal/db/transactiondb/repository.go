package transactiondb

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/moaabb/rinha-de-backend-2024-q1/internal/dto"
	"github.com/moaabb/rinha-de-backend-2024-q1/internal/models"
)

type TransactionRepository struct {
	db     *sql.DB
	logger *log.Logger
}

func NewTransactionRepository(db *sql.DB, logger *log.Logger) *TransactionRepository {
	return &TransactionRepository{
		db, logger,
	}
}

func (m *TransactionRepository) CreateTransaction(transaction dto.TransactionRequest, partyId int64) (*models.Party, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log.Fatal(err)
	}

	var party models.Party

	err = m.db.QueryRow(GetAccountBalance, partyId).Scan(
		&party.PartyId,
		&party.Limit,
		&party.Balance,
	)
	if err != nil {
		m.logger.Println("error getting account info")
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, NewErrorDefinition(err, "404")
		}
		return nil, NewErrorDefinition(err)
	}

	var newBalance int64
	if transaction.Type == models.Credit {
		newBalance = party.Balance + transaction.Value
	} else {
		newBalance = party.Balance - transaction.Value
	}

	if newBalance < -party.Limit {
		m.logger.Println("unsuported balance operation")
		return nil, NewErrorDefinition(errors.New("unsuported balance operation"))
	}

	_, err = m.db.Exec(CreateTransaction, transaction.Value, transaction.Type, transaction.Description, time.Now(), partyId)
	if err != nil {
		m.logger.Println("error creating transaction", err)
		tx.Rollback()
		return nil, NewErrorDefinition(err)
	}

	err = m.db.QueryRow(ChangeAccountBalance, newBalance, partyId).Scan(
		&party.PartyId,
		&party.Limit,
		&party.Balance,
	)
	if err != nil {
		m.logger.Println("error updating balance", err)
		tx.Rollback()
		return nil, NewErrorDefinition(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return &party, nil

}

func (m *TransactionRepository) GetAccountStatementByPartyId(partyId int64) (*models.Party, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	var party models.Party
	err := m.db.QueryRowContext(ctx, GetAccountBalance, partyId).Scan(
		&party.PartyId,
		&party.Limit,
		&party.Balance,
	)
	if err != nil {
		m.logger.Println("error getting party")
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, NewErrorDefinition(err, "404")
		}
		return nil, NewErrorDefinition(err)
	}

	rows, err := m.db.QueryContext(ctx, GetAccountStatementByPartyId, partyId)
	if err != nil {
		m.logger.Println("error getting transactions")
		return nil, NewErrorDefinition(err)
	}

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.TransactionId,
			&transaction.Value,
			&transaction.Type,
			&transaction.Description,
			&transaction.CreatedAt,
		)

		if err != nil {
			m.logger.Println("error parsing transaction")
			return nil, NewErrorDefinition(err)
		}

		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		m.logger.Println("error fetching transactions")
		return nil, NewErrorDefinition(err)

	}

	party.Transactions = transactions

	return &party, nil
}

func NewErrorDefinition(err error, errorCode ...string) *models.ErrorDefinition {
	code := "9999"

	if len(errorCode) != 0 {
		code = errorCode[0]
	}

	return &models.ErrorDefinition{
		Message:   err.Error(),
		ErrorCode: code,
	}
}
