package domain

import (
	"errors"
	"time"
)

type PixPayment struct {
	ID        int64         `json:"id"`
	Amount    float64       `json:"amount"`
	Status    PaymentStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
}

func NewPixPayment(amount float64) (*PixPayment, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be > 0")
	}
	return &PixPayment{Amount: amount, Status: StatusCreated}, nil
}

func (p *PixPayment) Authorize() error {
	if p.Status != StatusCreated {
		return errors.New("only CREATED payments can be authorized")
	}
	p.Status = StatusAuthorized
	return nil
}

func (p *PixPayment) Settle() error {
	if p.Status != StatusAuthorized {
		return errors.New("only AUTHORIZED payments can be settled")
	}
	p.Status = StatusSettled
	return nil
}
