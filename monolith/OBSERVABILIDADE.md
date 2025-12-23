# Observabilidade em Tempo Real - Monitor PIX

## 📊 Visão Geral

Esta aplicação implementa **observabilidade em tempo real** usando **Server-Sent Events (SSE)** para monitorar mudanças de status de pagamentos PIX instantaneamente, sem precisar consultar logs ou fazer polling.

## 🎯 Problema Resolvido

**Antes:** As mudanças de status aconteciam em background e eram difíceis de monitorar sem logs ou polling constante.

**Agora:** Interface visual em tempo real que mostra cada mudança de status instantaneamente via Server-Sent Events (SSE).

## 🚀 Como Usar

### 1. Acessar a Página de Monitoramento

Abra no navegador:
```
http://localhost:8080/monitor
```

### 2. Criar um Pagamento PIX

Em outro terminal ou aba, crie um pagamento:

```bash
curl -X POST http://localhost:8080/payments/pix \
  -H 'Content-Type: application/json' \
  -d '{"amount": 123.45}'
```

**Anote o ID retornado** (ex: `{"id": 1, ...}`)

### 3. Monitorar em Tempo Real

1. Na página de monitoramento, digite o ID do pagamento
2. Clique em "Iniciar Monitoramento"
3. **Aguarde** enquanto o pagamento é processado em background
4. Veja as mudanças de status em tempo real:
   - 🟡 **CREATED** - Pagamento criado (imediato)
   - 🔵 **AUTHORIZED** - Autorizado pelo BACEN (após ~3 segundos)
   - 🟢 **SETTLED** - Liquidado (após ~6 segundos)

## 📡 Endpoints de Observabilidade

### Página de Monitoramento (HTML)
```
GET /monitor
```
Interface visual completa para monitorar pagamentos.

### SSE - Server-Sent Events
```
GET /payments/pix/monitor/{id}
```
Endpoint SSE que envia eventos em tempo real quando o status muda.

**Exemplo de uso direto:**
```bash
curl -N http://localhost:8080/payments/pix/monitor/1
```

Você verá eventos como:
```
event: initial
data: {"payment_id":1,"status":"CREATED","amount":123.45,"timestamp":"2024-01-15T10:30:00Z","message":"Status inicial do pagamento"}

event: status_change
data: {"payment_id":1,"status":"AUTHORIZED","amount":123.45,"timestamp":"2024-01-15T10:30:00.5Z","message":"Pagamento PIX autorizado pelo BACEN"}

event: status_change
data: {"payment_id":1,"status":"SETTLED","amount":123.45,"timestamp":"2024-01-15T10:30:01.5Z","message":"Pagamento PIX liquidado com sucesso"}
```

## 🎨 Interface de Monitoramento

A página `/monitor` oferece:

- ✅ **Input para ID do pagamento**
- ✅ **Status atual visual** com badges coloridos
- ✅ **Log de eventos em tempo real**
- ✅ **Timestamps precisos** de cada mudança
- ✅ **Design moderno e responsivo**

### Cores dos Status

- 🟡 **CREATED** - Amarelo (pagamento criado)
- 🔵 **AUTHORIZED** - Azul (autorizado)
- 🟢 **SETTLED** - Verde (liquidado)

## 🔧 Como Funciona

### 1. Sistema de Eventos

```go
// EventBroadcaster gerencia clientes SSE
type EventBroadcaster struct {
    clients map[int64]map[chan PaymentEvent]bool
    mu      sync.RWMutex
}
```

### 2. Emissão de Eventos

Quando o status muda no use case:
```go
// Emite evento de mudança
eventBroadcaster.Broadcast(paymentID, PaymentEvent{
    PaymentID: paymentID,
    Status:    newStatus,
    Amount:    amount,
    Timestamp: time.Now(),
    Message:   "Pagamento autorizado",
})
```

### 3. Clientes SSE

Cada cliente conectado recebe eventos em tempo real:
```go
eventSource = new EventSource('/payments/pix/monitor/' + id);
eventSource.addEventListener('status_change', function(e) {
    const event = JSON.parse(e.data);
    updateStatus(event);
});
```

## 📊 Fluxo Completo

```
1. Cliente cria pagamento PIX (POST retorna imediatamente com CREATED)
   ↓
2. Use case emite evento CREATED
   ↓
3. SSE envia evento para clientes conectados
   ↓
4. Interface atualiza status visual
   ↓
5. Processamento em background inicia (goroutine)
   ↓
6. Após ~3 segundos: Use case emite evento AUTHORIZED
   ↓
7. SSE envia evento para clientes
   ↓
8. Interface atualiza status
   ↓
9. Após ~6 segundos: Use case emite evento SETTLED
   ↓
10. SSE envia evento final
   ↓
11. Interface mostra status final
```

## 🎓 Casos de Uso

### 1. Monitoramento Durante Desenvolvimento
- Veja exatamente quando cada status muda
- Entenda o timing do processo
- Debug visual do fluxo

### 2. Demonstração para Clientes
- Mostre o processo em tempo real
- Visual profissional e moderno
- Fácil de entender

### 3. Testes e Validação
- Verifique se os status mudam corretamente
- Confirme os tempos de processamento
- Valide o fluxo completo

## 💡 Vantagens sobre Logs

| Aspecto | Logs | SSE + Interface |
|---------|------|-----------------|
| **Tempo Real** | ❌ Precisa consultar | ✅ Instantâneo |
| **Visual** | ❌ Texto | ✅ Interface gráfica |
| **Fácil de Usar** | ❌ Terminal | ✅ Navegador |
| **Múltiplos Observadores** | ❌ Difícil | ✅ Múltiplos clientes |
| **Histórico** | ✅ Sim | ✅ Log de eventos |

## 🔍 Exemplo Prático

### Passo a Passo

1. **Inicie a aplicação:**
   ```bash
   cd monolith
   docker compose up -d
   ```

2. **Abra o monitor:**
   ```
   http://localhost:8080/monitor
   ```

3. **Em outro terminal, crie um pagamento:**
   ```bash
   curl -X POST http://localhost:8080/payments/pix \
     -H 'Content-Type: application/json' \
     -d '{"amount": 250.75}'
   ```
   
   Resposta: `{"id": 1, "amount": 250.75, ...}`

4. **No monitor, digite o ID (1) e clique "Iniciar Monitoramento"**

5. **Observe as mudanças em tempo real:**
   - Status muda de CREATED → AUTHORIZED → SETTLED
   - Log de eventos mostra cada mudança
   - Timestamps precisos

## 🛠️ Arquitetura Técnica

### Componentes

1. **EventBroadcaster** (`http/events.go`)
   - Gerencia clientes SSE
   - Distribui eventos
   - Thread-safe

2. **Use Case** (`application/create_pix_payment_usecase.go`)
   - Emite eventos em cada mudança
   - Integrado com o broadcaster

3. **SSE Endpoint** (`http/payments_facade.go`)
   - `/payments/pix/monitor/{id}`
   - Stream de eventos em tempo real

4. **Interface HTML** (`http/payments_facade.go`)
   - `/monitor`
   - Página completa de monitoramento

## 📈 Próximos Passos

Para melhorar ainda mais a observabilidade:

1. **Métricas Prometheus**
   - Contadores de pagamentos por status
   - Tempo médio de processamento
   - Taxa de sucesso

2. **Tracing Distribuído**
   - Jaeger ou Zipkin
   - Rastreamento completo do fluxo

3. **Alertas**
   - Notificações quando algo falha
   - Alertas de tempo de processamento

4. **Dashboard**
   - Gráficos de status
   - Estatísticas em tempo real

## 🔗 Links Úteis

- **Monitor:** http://localhost:8080/monitor
- **Swagger:** http://localhost:8080/swagger/index.html
- **Health:** http://localhost:8080/health

