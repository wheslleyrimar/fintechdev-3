package notifications

import (
	"context"
	"fintech-monolith/domains/notifications"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PgNotificationRepository implementa NotificationRepository usando PostgreSQL

type PgNotificationRepository struct {
	pool *pgxpool.Pool
}

func NewPgNotificationRepository(pool *pgxpool.Pool) *PgNotificationRepository {
	return &PgNotificationRepository{pool: pool}
}

func (r *PgNotificationRepository) Save(notification *notifications.Notification) (*notifications.Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int64
	var createdAt time.Time
	err := r.pool.QueryRow(ctx,
		"INSERT INTO notifications (payment_id, type, recipient, message, status) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at",
		notification.PaymentID, notification.Type, notification.Recipient, notification.Message, string(notification.Status),
	).Scan(&id, &createdAt)

	if err != nil {
		return nil, err
	}

	notification.ID = id
	notification.CreatedAt = createdAt
	return notification, nil
}

func (r *PgNotificationRepository) FindByID(id int64) (*notifications.Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var notification notifications.Notification
	var status string
	err := r.pool.QueryRow(ctx,
		"SELECT id, payment_id, type, recipient, message, status, created_at FROM notifications WHERE id = $1",
		id,
	).Scan(&notification.ID, &notification.PaymentID, &notification.Type, &notification.Recipient, &notification.Message, &status, &notification.CreatedAt)

	if err != nil {
		return nil, err
	}

	notification.Status = notifications.NotificationStatus(status)
	return &notification, nil
}
