package db

import (
	"context"
	"github.com/ktrntrsv/userBalanceService/internal/domain/models"
	"github.com/ktrntrsv/userBalanceService/pkg/logger"
	"github.com/ktrntrsv/userBalanceService/pkg/postgresql"
)

type AccountRepository struct {
	client postgresql.Client
	logger logger.Interface
}

func NewAccountRepository(client postgresql.Client, logger logger.Interface) *AccountRepository {
	return &AccountRepository{
		client: client,
		logger: logger,
	}
}

func (r *AccountRepository) GetById(ctx context.Context, id string) (models.Account, error) {
	return models.Account{}, nil
}

func (r *AccountRepository) UpdateBalance(ctx context.Context, account models.Account) error {
	return nil
}
