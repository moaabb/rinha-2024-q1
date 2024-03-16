package transactiondb

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moaabb/rinha-de-backend-2024-q1/internal/dto"
	"github.com/moaabb/rinha-de-backend-2024-q1/internal/models"
)

const CONTEXT_TIMEOUT = time.Second * 3

type TransactionRepository struct {
	db     *pgxpool.Pool
	logger *log.Logger
}

func NewTransactionRepository(db *pgxpool.Pool, logger *log.Logger) *TransactionRepository {
	return &TransactionRepository{
		db, logger,
	}
}

func (m *TransactionRepository) CreateTransaction(transaction dto.TransactionRequest, partyId int64) (*models.Party, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CONTEXT_TIMEOUT)
	defer cancel()
	tx, err := m.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return nil, NewErrorDefinition(err, tx)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()

	var party models.Party

	err = m.db.QueryRow(ctx, GetAccountBalance, partyId).Scan(
		&party.PartyId,
		&party.Limit,
		&party.Balance,
	)
	if err != nil {
		m.logger.Println("error getting account info")
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, NewErrorDefinition(err, tx, "404")
		}
		return nil, NewErrorDefinition(err, tx)
	}

	var query string
	if transaction.Type == models.Credit {
		query = AddToAccountBalance
	} else {
		ok := validateAccountBalance(party.Balance, party.Limit, transaction.Value)
		if !ok {
			m.logger.Println("unsuported balance operation")
			return nil, NewErrorDefinition(errors.New("unsuported balance operation"), tx)
		}
		query = SubtractFromAccountBalance
	}

	_, err = m.db.Exec(ctx, CreateTransaction, transaction.Value, transaction.Type, transaction.Description, time.Now(), partyId)
	if err != nil {
		m.logger.Println("error creating transaction", err)
		return nil, NewErrorDefinition(err, tx)
	}

	err = m.db.QueryRow(ctx, query, transaction.Value, partyId).Scan(
		&party.PartyId,
		&party.Limit,
		&party.Balance,
	)
	if err != nil {
		m.logger.Println("error updating balance", err)
		return nil, NewErrorDefinition(err, tx)
	}

	return &party, nil

}

func (m *TransactionRepository) GetAccountStatementByPartyId(partyId int64) (*models.Party, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CONTEXT_TIMEOUT)
	defer cancel()
	tx, err := m.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return nil, NewErrorDefinition(err, tx)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()

	var party models.Party
	err = m.db.QueryRow(ctx, GetAccountBalance, partyId).Scan(
		&party.PartyId,
		&party.Limit,
		&party.Balance,
	)
	if err != nil {
		m.logger.Println("error getting party")
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, NewErrorDefinition(err, tx, "404")
		}
		return nil, NewErrorDefinition(err, tx)
	}

	rows, err := m.db.Query(ctx, GetAccountStatementByPartyId, partyId)
	if err != nil {
		m.logger.Println("error getting transactions")
		return nil, NewErrorDefinition(err, tx)
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
			return nil, NewErrorDefinition(err, tx)
		}

		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		m.logger.Println("error fetching transactions")
		return nil, NewErrorDefinition(err, tx)

	}

	party.Transactions = transactions

	return &party, nil
}

func NewErrorDefinition(err error, tx pgx.Tx, errorCode ...string) *models.ErrorDefinition {
	code := "9999"

	if len(errorCode) != 0 {
		code = errorCode[0]
	}

	return &models.ErrorDefinition{
		Message:   err.Error(),
		ErrorCode: code,
	}
}

func validateAccountBalance(balance int64, limit int64, amount int64) bool {
	if (balance - amount) < -limit {
		return false

	}

	return true
}
