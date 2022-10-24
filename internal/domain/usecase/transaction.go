package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/ktrntrsv/userBalanceService/internal/domain/models"
)

type transactionRepository interface {
	GetById(ctx context.Context, id uuid.UUID) (models.Transaction, error)
	UpdateStatus(ctx context.Context, id uuid.UUID) error
	StartTransaction(ctx context.Context, dto models.TransactionStartDTO) (models.Transaction, error)
}

type TransactionUsecase struct {
	transactRepo      transactionRepository
	accountRepository accountRepository
}

func NewTransactionUsecase(transactRepo transactionRepository, accRepo accountRepository) *TransactionUsecase {
	return &TransactionUsecase{transactRepo: transactRepo, accountRepository: accRepo}
}

func (t *TransactionUsecase) StartTransaction(ctx context.Context, dto models.TransactionStartDTO) (uuid.UUID, error) {
	sender, err := t.accountRepository.GetById(ctx, dto.AccountFromId)
	if err != nil {
		// todo обработать ошибку (сетевая проблема или пользователь не существует)
		return uuid.UUID{}, err
	}

	_, err = t.accountRepository.GetById(ctx, dto.AccountToId)
	if err != nil {
		// todo обработать ошибку (сетевая проблема или пользователь не существует)
		return uuid.UUID{}, err
	}

	if err := sender.UpdateBalance(-dto.Amount); err != nil {
		return uuid.UUID{}, err
	}

	transaction, err := t.transactRepo.StartTransaction(ctx, dto)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = t.accountRepository.UpdateBalance(ctx, sender)
	if err != nil {
		return uuid.UUID{}, err
	}

	return transaction.Id, nil
}

func (t *TransactionUsecase) ApproveTransaction(ctx context.Context, transactID uuid.UUID) error {
	transact, err := t.transactRepo.GetById(ctx, transactID)
	if err != nil {
		return err
	}

	transact.Status = models.Approved
	err = t.transactRepo.UpdateStatus(ctx, transactID)
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
}

func (t *TransactionUsecase) AbortTransaction(ctx context.Context, transactID uuid.UUID) error {
	transact, err := t.transactRepo.GetById(ctx, transactID)
	if err != nil {
		return err
	}

	transact.Status = models.Abort
	err = t.transactRepo.UpdateStatus(ctx, transactID)
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
}
