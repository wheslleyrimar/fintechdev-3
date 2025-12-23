package domain

// NotificationClient é a interface para comunicação com o serviço de notificações
// No contexto de microsserviços, isso é uma chamada HTTP ou evento
type NotificationClient interface {
	SendPaymentCreatedNotification(paymentID int64, amount float64) error
	SendPaymentAuthorizedNotification(paymentID int64, amount float64) error
	SendPaymentSettledNotification(paymentID int64, amount float64) error
}
