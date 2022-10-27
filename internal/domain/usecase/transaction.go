package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ktrntrsv/userBalanceService/internal/domain/models"
)

var ErrTransactionNotFound = errors.New("transaction not found")

type transactionRepository interface {
	GetById(ctx context.Context, id uuid.UUID) (models.Transaction, error)
	UpdateStatus(ctx context.Context, transaction models.Transaction) error
	StartTransaction(ctx context.Context, dto models.TransactionStartDTO) (uuid.UUID, error)
	WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error
}

type TransactionUsecase struct {
	transactRepository transactionRepository
	accountRepository  accountRepository
}

func NewTransactionUsecase(transactRepo transactionRepository, accRepo accountRepository) *TransactionUsecase {
	return &TransactionUsecase{
		transactRepository: transactRepo,
		accountRepository:  accRepo}
}

func (t *TransactionUsecase) StartTransaction(ctx context.Context, dto models.TransactionStartDTO) (uuid.UUID, error) {
	var uId uuid.UUID

	err := t.transactRepository.WithinTransaction(ctx, func(txCtx context.Context) error {
		sender, err := t.accountRepository.GetById(ctx, dto.AccountFromId)
		if err != nil {
			uId = uuid.UUID{}
			return err
		}

		_, err = t.accountRepository.GetById(ctx, dto.AccountToId)
		if err != nil {
			uId = uuid.UUID{}
			return err
		}

		if err := sender.UpdateBalance(-dto.Amount); err != nil {
			uId = uuid.UUID{}
			return err
		}

		id, err := t.transactRepository.StartTransaction(ctx, dto)
		if err != nil {
			uId = uuid.UUID{}
			return err
		}

		err = t.accountRepository.UpdateBalance(ctx, sender)
		if err != nil {
			uId = uuid.UUID{}
			return err
		}
		uId = id
		return nil
	})

	return uId, err
}

func (t *TransactionUsecase) ApproveTransaction(ctx context.Context, transactID uuid.UUID) error {

	return t.transactRepository.WithinTransaction(ctx, func(ctx context.Context) error {
		transact, err := t.transactRepository.GetById(ctx, transactID)
		if err != nil {
			return err
		}

		if transact.Status != models.Pending {
			return fmt.Errorf("can not approve transaction without pending status")
		}

		transact.Status = models.Approved
		err = t.transactRepository.UpdateStatus(ctx, transact)
		if err != nil {
			return err
		}

		receiver, err := t.accountRepository.GetById(ctx, transact.AccountToId)
		if err != nil {
			// todo обработать ошибку (сетевая проблема или пользователь не существует)
			return err
		}

		err = receiver.UpdateBalance(transact.Amount)
		if err != nil {
			return err
		}

		err = t.accountRepository.UpdateBalance(ctx, receiver)
		if err != nil {
			return err
		}

		return nil
	})

}

func (t *TransactionUsecase) AbortTransaction(ctx context.Context, transactID uuid.UUID) error {
	return t.transactRepository.WithinTransaction(ctx, func(ctx context.Context) error {
		transact, err := t.transactRepository.GetById(ctx, transactID)
		if err != nil {
			return err
		}

		if transact.Status != models.Pending {
			return fmt.Errorf("can not abort transaction without pending status")
		}

		transact.Status = models.Abort
		err = t.transactRepository.UpdateStatus(ctx, transact)
		if err != nil {
			return err
		}

		sender, err := t.accountRepository.GetById(ctx, transact.AccountFromId)
		if err != nil {
			// todo обработать ошибку (сетевая проблема или пользователь не существует)
			return err
		}

		err = sender.UpdateBalance(transact.Amount)
		if err != nil {
			return err
		}

		err = t.accountRepository.UpdateBalance(ctx, sender)
		if err != nil {
			return err
		}

		return nil
	})
}
