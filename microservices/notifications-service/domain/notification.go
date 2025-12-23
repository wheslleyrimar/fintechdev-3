package domain

import "time"

type Notification struct {
	ID        int64             `json:"id"`
	Type      string            `json:"type"`
	Recipient string            `json:"recipient"`
	Message   string            `json:"message"`
	Status    NotificationStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
}

type NotificationStatus string

const (
	StatusPending NotificationStatus = "PENDING"
	StatusSent    NotificationStatus = "SENT"
	StatusFailed  NotificationStatus = "FAILED"
)

func NewNotification(notificationType, recipient, message string) *Notification {
	return &Notification{
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
