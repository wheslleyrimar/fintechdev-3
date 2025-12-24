package pix

import (
	"fintech-monolith/domains/payments"
	"log"
	"time"
)

type BacenPixGateway struct{}

func NewBacenPixGateway() *BacenPixGateway {
	return &BacenPixGateway{}
}

func (g *BacenPixGateway) NotifyCreation(payment *payments.PixPayment) {
	// Simula notificação para o BACEN
	log.Printf("BACEN: Notificando criação de pagamento PIX - ID: %d, Valor: R$ %.2f", payment.ID, payment.Amount)
	// Simula latência de rede
	time.Sleep(100 * time.Millisecond)
	log.Printf("BACEN: Pagamento PIX registrado no sistema - ID: %d", payment.ID)
}

func (g *BacenPixGateway) Authorize(payment *payments.PixPayment) error {
	// Simula autorização no BACEN
	log.Printf("BACEN: Processando autorização de pagamento PIX - ID: %d", payment.ID)
	time.Sleep(200 * time.Millisecond)

	// Simula validações do BACEN
	if payment.Amount > 100000 {
		log.Printf("BACEN: Pagamento acima de R$ 100.000 requer análise adicional")
	}

	log.Printf("BACEN: Pagamento PIX autorizado - ID: %d", payment.ID)
	return nil
}

func (g *BacenPixGateway) Settle(payment *payments.PixPayment) error {
	// Simula liquidação no BACEN
	log.Printf("BACEN: Processando liquidação de pagamento PIX - ID: %d", payment.ID)
	time.Sleep(300 * time.Millisecond)
	log.Printf("BACEN: Pagamento PIX liquidado - ID: %d, Valor transferido: R$ %.2f", payment.ID, payment.Amount)
	return nil
}
