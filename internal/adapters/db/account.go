package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/ktrntrsv/userBalanceService/internal/domain/models"
	"github.com/ktrntrsv/userBalanceService/internal/domain/usecase"
	"github.com/ktrntrsv/userBalanceService/pkg/logger"
)

type AccountRepository struct {
	*Database
	logger logger.Interface
}

func NewAccountRepository(client *Database, logger logger.Interface) *AccountRepository {
	return &AccountRepository{
		Database: client,
		logger:   logger,
	}
}

func (r *AccountRepository) GetById(ctx context.Context, id uuid.UUID) (models.Account, error) {
	query := `SELECT id, balance from account WHERE id = $1;`
	var account models.Account
	row := r.model(ctx).QueryRow(ctx, query, id)
	err := row.Scan(&account.Id, &account.Balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Account{}, usecase.ErrAccountNotFound
		}
		return models.Account{}, err
	}
	return account, nil
}

func (r *AccountRepository) UpdateBalance(ctx context.Context, account models.Account) error {
	query := `UPDATE account SET balance=$1 WHERE id=$2;`
	_, err := r.model(ctx).Exec(ctx, query, account.Balance, account.Id)
	if err != nil {
		return err
	}
	return nil
}
