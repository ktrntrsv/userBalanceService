package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/ktrntrsv/userBalanceService/internal/domain/models"
	"github.com/ktrntrsv/userBalanceService/internal/domain/usecase"
	"github.com/ktrntrsv/userBalanceService/pkg/logger"
)

type TransactionRepository struct {
	*Database
	logger logger.Interface
}

func NewTransactionRepository(client *Database, logger logger.Interface) *TransactionRepository {
	return &TransactionRepository{
		Database: client,
		logger:   logger,
	}
}

func (r *TransactionRepository) GetById(ctx context.Context, id uuid.UUID) (models.Transaction, error) {
	query := `SELECT id, amount, account_to_id, account_from_id, status, service_id 
			  FROM transaction WHERE id = $1;`

	var transaction models.Transaction

	row := r.model(ctx).QueryRow(ctx, query, id)
	err := row.Scan(&transaction.Id, &transaction.Amount, &transaction.AccountToId, &transaction.AccountFromId, &transaction.Status, &transaction.ServiceId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Transaction{}, usecase.ErrTransactionNotFound
		}
		return models.Transaction{}, err
	}

	return transaction, nil
}

func (r *TransactionRepository) UpdateStatus(ctx context.Context, transaction models.Transaction) error {
	query := `UPDATE transaction SET status=$1 WHERE id=$2;`
	_, err := r.model(ctx).Exec(ctx, query, transaction.Status, transaction.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *TransactionRepository) StartTransaction(ctx context.Context, dto models.TransactionStartDTO) (uuid.UUID, error) {
	id := uuid.New()
	query := `INSERT INTO transaction (id, amount, account_to_id, account_from_id, status, service_id)
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.model(ctx).Exec(ctx, query, id, dto.Amount, dto.AccountToId, dto.AccountFromId, models.Pending, dto.ServiceId)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}
