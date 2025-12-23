package payments

type PixGateway interface {
	NotifyCreation(payment *PixPayment)
	Authorize(payment *PixPayment) error
	Settle(payment *PixPayment) error
}
