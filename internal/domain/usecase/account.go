package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ktrntrsv/userBalanceService/internal/domain/models"
)

var ErrAccountNotFound = errors.New("account not found")

type accountRepository interface {
	GetById(ctx context.Context, id uuid.UUID) (models.Account, error)
	UpdateBalance(ctx context.Context, account models.Account) error
}

type AccountUsecase struct {
	accRepo accountRepository
}

func NewAccountUsecase(accRepo accountRepository) *AccountUsecase {
	return &AccountUsecase{accRepo: accRepo}
}

func (a *AccountUsecase) EnrollBalance(ctx context.Context, accountId uuid.UUID, sum float64) error {
	acc, err := a.accRepo.GetById(ctx, accountId)
	if err != nil {
		return fmt.Errorf("can not find account %v: %w", accountId, err)
	}
	if err := acc.UpdateBalance(sum); err != nil {
		return err
	}

	if err := a.accRepo.UpdateBalance(ctx, acc); err != nil {
		return err
	}

	return nil
}

func (a *AccountUsecase) GetBalance(ctx context.Context, accountId uuid.UUID) (float64, error) {
	acc, err := a.accRepo.GetById(ctx, accountId)
	if err != nil {
		return 0, err
	}

	return acc.Balance, nil
}
