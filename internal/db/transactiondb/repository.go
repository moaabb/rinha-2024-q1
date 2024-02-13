package transactiondb

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

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

func (m *TransactionRepository) CreateTransaction(transaction models.Transaction, partyId int64) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var party models.Party

	tx, err := m.db.Begin()
	if err != nil {
		return nil, err
	}

	err = m.db.QueryRowContext(ctx, GetAccountBalance, partyId).Scan(
		&party.PartyId,
		&party.Limit,
		&party.Balance,
	)
	if err != nil {
		return nil, err
	}

	var newBalance int64
	if transaction.Type == models.Credit {
		newBalance = party.Balance + transaction.Value
	} else {
		newBalance = party.Balance - transaction.Value
	}

	if newBalance < -party.Limit {
		return nil, errors.New("Unsuported balance")
	}

	var createdTransaction models.Transaction
	err = m.db.QueryRowContext(ctx, CreateTransaction, transaction.Value, transaction.Type, transaction.Description, time.Now(), partyId).Scan(
		&createdTransaction.Value,
		&createdTransaction.Type,
		&createdTransaction.Description,
		&createdTransaction.CreatedAt,
		&createdTransaction.PartyId,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = m.db.ExecContext(ctx, ChangeAccountBalance, newBalance, partyId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &createdTransaction, nil

}

func (m *TransactionRepository) GetAccountStatementByPartyId(partyId int64) (*models.Party, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, GetAccountStatementByPartyId, partyId)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var party models.Party
	err = m.db.QueryRowContext(ctx, GetAccountBalance, partyId).Scan(&party.PartyId, party.Balance, party.Limit)
	if err != nil {
		return nil, err
	}

	party.Transactions = transactions

	return &party, nil
}
