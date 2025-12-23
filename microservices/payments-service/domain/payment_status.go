package domain

type PaymentStatus string

const (
	StatusCreated    PaymentStatus = "CREATED"
	StatusAuthorized PaymentStatus = "AUTHORIZED"
	StatusSettled    PaymentStatus = "SETTLED"
)
