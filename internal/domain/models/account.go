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
	balance float64
}

func (a *Account) UpdateBalance(amount float64) error {
	if a.balance+amount < 0 {
		return ErrorNotEnoughFunds
	}
	a.balance += amount
	return nil
}

func (a *Account) GetBalance() float64 {
	return a.balance
}
