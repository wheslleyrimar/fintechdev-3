package domain

// PixGateway é a interface para comunicação com o gateway do BACEN
type PixGateway interface {
	NotifyCreation(payment *PixPayment)
	Authorize(payment *PixPayment) error
	Settle(payment *PixPayment) error
}

