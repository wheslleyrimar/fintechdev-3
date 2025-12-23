package api

import (
	"fintech-payments-service/domain"
	"log"
	"sync"
)

// EventBroadcaster gerencia os clientes SSE conectados
type EventBroadcaster struct {
	clients map[int64]map[chan domain.PaymentEvent]bool
	mu      sync.RWMutex
}

var globalBroadcaster = &EventBroadcaster{
	clients: make(map[int64]map[chan domain.PaymentEvent]bool),
}

// Subscribe adiciona um cliente para receber eventos de um pagamento específico
func (eb *EventBroadcaster) Subscribe(paymentID int64) chan domain.PaymentEvent {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	ch := make(chan domain.PaymentEvent, 10) // Buffer de 10 eventos

	if eb.clients[paymentID] == nil {
		eb.clients[paymentID] = make(map[chan domain.PaymentEvent]bool)
	}
	eb.clients[paymentID][ch] = true

	log.Printf("Event: Cliente inscrito para pagamento %d (total: %d)", paymentID, len(eb.clients[paymentID]))
	return ch
}

// Unsubscribe remove um cliente
func (eb *EventBroadcaster) Unsubscribe(paymentID int64, ch chan domain.PaymentEvent) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if clients, ok := eb.clients[paymentID]; ok {
		delete(clients, ch)
		if len(clients) == 0 {
			delete(eb.clients, paymentID)
		}
	}
	close(ch)
	log.Printf("Event: Cliente removido para pagamento %d", paymentID)
}

// Broadcast envia um evento para todos os clientes de um pagamento
func (eb *EventBroadcaster) Broadcast(paymentID int64, event domain.PaymentEvent) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if clients, ok := eb.clients[paymentID]; ok {
		for ch := range clients {
			select {
			case ch <- event:
			default:
				// Canal cheio, remover cliente
				log.Printf("Event: Canal cheio para pagamento %d, removendo cliente", paymentID)
				go func(c chan domain.PaymentEvent) {
					eb.mu.Lock()
					delete(clients, c)
					close(c)
					eb.mu.Unlock()
				}(ch)
			}
		}
	}
}

// GetBroadcaster retorna o broadcaster global
func GetBroadcaster() *EventBroadcaster {
	return globalBroadcaster
}

