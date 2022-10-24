package models

import "github.com/google/uuid"

type TransactionStartDTO struct {
	Amount        float64   `json:"amount" binding:"required"`
	AccountToId   uuid.UUID `json:"accountToId" binding:"required"`
	AccountFromId uuid.UUID `json:"accountFromId" binding:"required"`
	ServiceId     *string   `json:"serviceId"`
}

type Transaction struct {
	Id            uuid.UUID
	Amount        float64
	AccountToId   uuid.UUID
	AccountFromId uuid.UUID
	Status        Status
	ServiceId     *string
}

type Status int

const (
	Pending Status = iota
	Approved
	Abort
)

func (t *Transaction) UpdateStatus(status Status) {
	t.Status = status
}
