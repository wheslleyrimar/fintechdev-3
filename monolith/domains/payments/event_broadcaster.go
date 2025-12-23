package payments

import "time"

// PaymentEvent representa um evento de mudança de status
type PaymentEvent struct {
	PaymentID int64         `json:"payment_id"`
	Status    PaymentStatus `json:"status"`
	Amount    float64       `json:"amount"`
	Timestamp time.Time     `json:"timestamp"`
	Message   string        `json:"message"`
}

// EventBroadcaster interface para emitir eventos de mudança de status
type EventBroadcaster interface {
	Broadcast(paymentID int64, event PaymentEvent)
}

