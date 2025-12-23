package persistence

import (
	"context"
	"fintech-notifications-service/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgNotificationRepository struct {
	pool *pgxpool.Pool
}

func NewPgNotificationRepository(pool *pgxpool.Pool) *PgNotificationRepository {
	return &PgNotificationRepository{pool: pool}
}

func (r *PgNotificationRepository) Save(notification *domain.Notification) (*domain.Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int64
	err := r.pool.QueryRow(ctx,
		"INSERT INTO notifications (payment_id, type, recipient, message, status, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		notification.PaymentID, notification.Type, notification.Recipient, notification.Message, string(notification.Status), notification.CreatedAt,
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	notification.ID = id
	return notification, nil
}

func (r *PgNotificationRepository) FindByID(id int64) (*domain.Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var notification domain.Notification
	var status string
	err := r.pool.QueryRow(ctx,
		"SELECT id, payment_id, type, recipient, message, status, created_at FROM notifications WHERE id = $1",
		id,
	).Scan(&notification.ID, &notification.PaymentID, &notification.Type, &notification.Recipient, &notification.Message, &status, &notification.CreatedAt)

	if err != nil {
		return nil, err
	}

	notification.Status = domain.NotificationStatus(status)
	return &notification, nil
}

func (r *PgNotificationRepository) FindAll() ([]*domain.Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx,
		"SELECT id, payment_id, type, recipient, message, status, created_at FROM notifications ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		var notification domain.Notification
		var status string
		if err := rows.Scan(&notification.ID, &notification.PaymentID, &notification.Type, &notification.Recipient, &notification.Message, &status, &notification.CreatedAt); err != nil {
			return nil, err
		}
		notification.Status = domain.NotificationStatus(status)
		notifications = append(notifications, &notification)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}
