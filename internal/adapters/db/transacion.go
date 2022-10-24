package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/ktrntrsv/userBalanceService/internal/domain/models"
	"github.com/ktrntrsv/userBalanceService/pkg/logger"
	"github.com/ktrntrsv/userBalanceService/pkg/postgresql"
)

type TransactionRepository struct {
	client postgresql.Client
	logger logger.Interface
}

func NewTransactionRepository(client postgresql.Client, logger logger.Interface) *TransactionRepository {
	return &TransactionRepository{
		client: client,
		logger: logger,
	}
}

func (r *TransactionRepository) GetById(ctx context.Context, id string) (models.Transaction, error) {
	return models.Transaction{}, nil
}

func (r *TransactionRepository) UpdateStatus(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (r *TransactionRepository) StartTransaction(ctx context.Context, dto models.TransactionStartDTO) (models.Transaction, error) {
	return nil, nil
}
