package payments

import (
	"context"
	"fintech-monolith/domains/payments"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgPixPaymentRepository struct {
	pool *pgxpool.Pool
}

func NewPgPixPaymentRepository(pool *pgxpool.Pool) *PgPixPaymentRepository {
	return &PgPixPaymentRepository{pool: pool}
}

func (r *PgPixPaymentRepository) Save(payment *payments.PixPayment) (*payments.PixPayment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int64
	var createdAt time.Time
	err := r.pool.QueryRow(ctx,
		"INSERT INTO pix_payments (amount, status) VALUES ($1, $2) RETURNING id, created_at",
		payment.Amount, string(payment.Status),
	).Scan(&id, &createdAt)

	if err != nil {
		return nil, err
	}

	payment.ID = id
	payment.CreatedAt = createdAt
	return payment, nil
}

func (r *PgPixPaymentRepository) FindByID(id int64) (*payments.PixPayment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var payment payments.PixPayment
	var status string
	err := r.pool.QueryRow(ctx,
		"SELECT id, amount, status, created_at FROM pix_payments WHERE id = $1",
		id,
	).Scan(&payment.ID, &payment.Amount, &status, &payment.CreatedAt)

	if err != nil {
		return nil, err
	}

	payment.Status = payments.PaymentStatus(status)
	return &payment, nil
}

func (r *PgPixPaymentRepository) UpdateStatus(id int64, status payments.PaymentStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.pool.Exec(ctx,
		"UPDATE pix_payments SET status = $1 WHERE id = $2",
		string(status), id,
	)

	return err
}

func (r *PgPixPaymentRepository) FindAll() ([]*payments.PixPayment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.pool.Query(ctx,
		"SELECT id, amount, status, created_at FROM pix_payments ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paymentsList []*payments.PixPayment
	for rows.Next() {
		var payment payments.PixPayment
		var status string
		if err := rows.Scan(&payment.ID, &payment.Amount, &status, &payment.CreatedAt); err != nil {
			return nil, err
		}
		payment.Status = payments.PaymentStatus(status)
		paymentsList = append(paymentsList, &payment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return paymentsList, nil
}
