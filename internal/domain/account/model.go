package account

import (
	"errors"
)

type Account struct {
	Id      string
	balance float64
}

func (a *Account) UpdateBalance(sum float64) error {
	if a.balance+sum < 0 {
		return errors.New("insufficient funds")
	}

	a.balance += sum
	return nil
}

func (a *Account) GetBalance() float64 {
	return a.balance
}
