package payments

type PixPaymentRepository interface {
	Save(payment *PixPayment) (*PixPayment, error)
	FindByID(id int64) (*PixPayment, error)
	FindAll() ([]*PixPayment, error)
	UpdateStatus(id int64, status PaymentStatus) error
}
