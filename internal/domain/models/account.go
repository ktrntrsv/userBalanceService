package models

import (
	"errors"
	"github.com/google/uuid"
)

var (
	ErrorNotEnoughFunds = errors.New("not enough funds")
)

type Account struct {
	Id      uuid.UUID
	Balance float64
}

func (a *Account) UpdateBalance(amount float64) error {
	if a.Balance+amount < 0 {
		return ErrorNotEnoughFunds
	}
	a.Balance += amount
	return nil
}
