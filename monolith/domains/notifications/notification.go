package notifications

import "time"

type Notification struct {
	ID        int64
	PaymentID int64 // Associação com o pagamento PIX
	Type      string
	Recipient string
	Message   string
	Status    NotificationStatus
	CreatedAt time.Time
}

type NotificationStatus string

const (
	StatusPending NotificationStatus = "PENDING"
	StatusSent    NotificationStatus = "SENT"
	StatusFailed  NotificationStatus = "FAILED"
)

func NewNotification(paymentID int64, notificationType, recipient, message string) *Notification {
	return &Notification{
		PaymentID: paymentID,
		Type:      notificationType,
		Recipient: recipient,
		Message:   message,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}
}

func (n *Notification) MarkAsSent() {
	n.Status = StatusSent
}

func (n *Notification) MarkAsFailed() {
	n.Status = StatusFailed
}
