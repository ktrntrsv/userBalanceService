package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ktrntrsv/userBalanceService/internal/domain/models"
	"github.com/ktrntrsv/userBalanceService/internal/domain/usecase"
	"github.com/ktrntrsv/userBalanceService/pkg/logger"
)

type AccountRepository struct {
	client *pgxpool.Pool
	logger logger.Interface
}

func NewAccountRepository(client *pgxpool.Pool, logger logger.Interface) *AccountRepository {
	return &AccountRepository{
		client: client,
		logger: logger,
	}
}

func (r *AccountRepository) GetById(ctx context.Context, id uuid.UUID) (models.Account, error) {
	query := `SELECT id, balance from account WHERE id = $1;`

	var account models.Account
	fmt.Println("client", r, r.client)
	row := r.client.QueryRow(ctx, query, id)
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
	_, err := r.client.Exec(ctx, query, account.Balance, account.Id)
	if err != nil {
		return err
	}
	return nil
}
