package application

import (
	"fintech-monolith/domains/notifications"
	"fintech-monolith/domains/payments"
	"log"
	"time"
)

// CreatePixPaymentUseCase cria um pagamento PIX e simula o fluxo completo:
// 1. Cria o pagamento (CREATED)
// 2. Autoriza no BACEN (AUTHORIZED)
// 3. Liquida o pagamento (SETTLED)
// No monólito, tudo está no mesmo processo - comunicação direta
type CreatePixPaymentUseCase struct {
	paymentRepo      payments.PixPaymentRepository
	notificationRepo notifications.NotificationRepository
	gateway          payments.PixGateway
	eventBroadcaster payments.EventBroadcaster
}

func NewCreatePixPaymentUseCase(
	paymentRepo payments.PixPaymentRepository,
	notificationRepo notifications.NotificationRepository,
	gateway payments.PixGateway,
	eventBroadcaster payments.EventBroadcaster,
) *CreatePixPaymentUseCase {
	return &CreatePixPaymentUseCase{
		paymentRepo:      paymentRepo,
		notificationRepo: notificationRepo,
		gateway:          gateway,
		eventBroadcaster: eventBroadcaster,
	}
}

func (uc *CreatePixPaymentUseCase) Execute(amount float64) (*payments.PixPayment, error) {
	// 1. Criar pagamento com status CREATED
	payment, err := payments.NewPixPayment(amount)
	if err != nil {
		return nil, err
	}

	log.Printf("PIX: Criando pagamento de R$ %.2f", amount)

	// 2. Salvar no banco (compartilhado) com status CREATED
	saved, err := uc.paymentRepo.Save(payment)
	if err != nil {
		return nil, err
	}

	log.Printf("PIX: Pagamento criado - ID: %d, Status: %s", saved.ID, saved.Status)

	// Emitir evento de criação
	if uc.eventBroadcaster != nil {
		uc.emitStatusEvent(saved.ID, saved.Status, saved.Amount, "Pagamento PIX criado")
	}

	// IMPORTANTE: Processar o resto do fluxo em background (goroutine)
	// Isso permite que o POST retorne imediatamente e o SSE tenha tempo de conectar
	go uc.processPaymentFlow(saved)

	// Retornar imediatamente com o pagamento criado
	return saved, nil
}

// processPaymentFlow processa o fluxo completo do pagamento em background
func (uc *CreatePixPaymentUseCase) processPaymentFlow(saved *payments.PixPayment) {
	// Delay inicial para dar tempo do SSE conectar
	time.Sleep(1 * time.Second)

	// 3. Notificar criação ao BACEN (simulação)
	uc.gateway.NotifyCreation(saved)

	// 4. Criar notificação de criação
	notification := notifications.NewNotification(
		saved.ID,
		"PAYMENT_CREATED",
		"user@example.com",
		"Pagamento PIX criado com sucesso",
	)
	_, _ = uc.notificationRepo.Save(notification)

	// 5. Simular autorização no BACEN (após 2 segundos para facilitar visualização)
	time.Sleep(2 * time.Second)
	err := saved.Authorize()
	if err != nil {
		log.Printf("PIX: Erro ao autorizar pagamento %d: %v", saved.ID, err)
		return
	}

	// Atualizar status no banco
	err = uc.paymentRepo.UpdateStatus(saved.ID, saved.Status)
	if err != nil {
		log.Printf("PIX: Erro ao atualizar status para AUTHORIZED: %v", err)
	} else {
		log.Printf("PIX: Pagamento autorizado - ID: %d, Status: %s", saved.ID, saved.Status)
		// Emitir evento de autorização
		if uc.eventBroadcaster != nil {
			uc.emitStatusEvent(saved.ID, saved.Status, saved.Amount, "Pagamento PIX autorizado pelo BACEN")
		}
	}

	// 6. Criar notificação de autorização
	authNotification := notifications.NewNotification(
		saved.ID,
		"PAYMENT_AUTHORIZED",
		"user@example.com",
		"Pagamento PIX autorizado pelo BACEN",
	)
	_, _ = uc.notificationRepo.Save(authNotification)

	// 7. Simular liquidação no BACEN (após mais 3 segundos para facilitar visualização)
	time.Sleep(3 * time.Second)
	err = saved.Settle()
	if err != nil {
		log.Printf("PIX: Erro ao liquidar pagamento %d: %v", saved.ID, err)
		return
	}

	// Atualizar status no banco
	err = uc.paymentRepo.UpdateStatus(saved.ID, saved.Status)
	if err != nil {
		log.Printf("PIX: Erro ao atualizar status para SETTLED: %v", err)
	} else {
		log.Printf("PIX: Pagamento liquidado - ID: %d, Status: %s", saved.ID, saved.Status)
		// Emitir evento de liquidação
		if uc.eventBroadcaster != nil {
			uc.emitStatusEvent(saved.ID, saved.Status, saved.Amount, "Pagamento PIX liquidado com sucesso")
		}
	}

	// 8. Criar notificação de liquidação
	settleNotification := notifications.NewNotification(
		saved.ID,
		"PAYMENT_SETTLED",
		"user@example.com",
		"Pagamento PIX liquidado com sucesso",
	)
	_, _ = uc.notificationRepo.Save(settleNotification)

	// 9. Notificar liquidação ao BACEN
	uc.gateway.Settle(saved)

	log.Printf("PIX: Fluxo completo finalizado - ID: %d, Status: %s", saved.ID, saved.Status)
}

// emitStatusEvent emite um evento de mudança de status
func (uc *CreatePixPaymentUseCase) emitStatusEvent(paymentID int64, status payments.PaymentStatus, amount float64, message string) {
	event := payments.PaymentEvent{
		PaymentID: paymentID,
		Status:    status,
		Amount:    amount,
		Timestamp: time.Now(),
		Message:   message,
	}
	uc.eventBroadcaster.Broadcast(paymentID, event)
}
