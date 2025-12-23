# Fluxo de Pagamento PIX - Simulação Completa

## 📋 O que a Aplicação Faz

Esta aplicação **simula um fluxo completo de pagamento PIX**, desde a criação até a liquidação final. Ela demonstra como um monólito processa transações financeiras de forma síncrona.

## 🔄 Fluxo Completo do PIX

### 1. Criação do Pagamento (CREATED)
- ✅ Valida o valor (deve ser > 0)
- ✅ Cria o pagamento com status `CREATED`
- ✅ Salva no banco de dados
- ✅ Notifica o BACEN sobre a criação
- ✅ Cria notificação para o usuário

### 2. Autorização (AUTHORIZED)
- ✅ Após 1 segundo (delay inicial) + 2 segundos, autoriza o pagamento
- ✅ Atualiza status para `AUTHORIZED`
- ✅ Processa autorização no BACEN (simulação)
- ✅ Cria notificação de autorização

### 3. Liquidação (SETTLED)
- ✅ Após mais 3 segundos, liquida o pagamento
- ✅ Atualiza status para `SETTLED`
- ✅ Processa liquidação no BACEN (simulação)
- ✅ Cria notificação de liquidação
- ✅ Retorna pagamento com status final

## ⏱️ Timeline do Processo

```
0s     → Pagamento criado (CREATED) - retorna imediatamente
1s     → Delay inicial (para SSE conectar)
3s     → Pagamento autorizado (AUTHORIZED) - após 1s + 2s
6s     → Pagamento liquidado (SETTLED) - após 3s + 3s
```

**Total:** ~6 segundos para completar o fluxo completo (processado em background)

> **Nota:** O POST retorna imediatamente com status `CREATED`. O fluxo completo (autorização e liquidação) acontece em background usando goroutines, permitindo monitoramento em tempo real via Server-Sent Events (SSE).

## 📊 Status do Pagamento

| Status | Descrição | Quando Ocorre |
|--------|-----------|---------------|
| `CREATED` | Pagamento criado | Imediatamente após criação (POST retorna) |
| `AUTHORIZED` | Autorizado pelo BACEN | Após ~3 segundos (1s delay + 2s) |
| `SETTLED` | Liquidado e finalizado | Após ~6 segundos (3s + 3s) |

## 🎯 O que Acontece em Cada Etapa

### Etapa 1: Criação
```json
{
  "id": 1,
  "amount": 123.45,
  "status": "CREATED",
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Ações:**
- Validação do valor
- Criação do registro no banco
- Notificação ao BACEN
- Notificação ao usuário

### Etapa 2: Autorização
```json
{
  "id": 1,
  "amount": 123.45,
  "status": "AUTHORIZED",
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Ações:**
- Validação no BACEN (simulada)
- Atualização do status no banco
- Notificação de autorização

### Etapa 3: Liquidação
```json
{
  "id": 1,
  "amount": 123.45,
  "status": "SETTLED",
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Ações:**
- Processamento de liquidação no BACEN
- Atualização final do status
- Notificação de liquidação
- Transferência concluída (simulada)

## 🔍 Logs do Processo

Quando você cria um pagamento, verá logs como:

```
PIX: Criando pagamento de R$ 123.45
PIX: Pagamento criado - ID: 1, Status: CREATED
BACEN: Notificando criação de pagamento PIX - ID: 1, Valor: R$ 123.45
BACEN: Pagamento PIX registrado no sistema - ID: 1
PIX: Pagamento autorizado - ID: 1, Status: AUTHORIZED
BACEN: Processando autorização de pagamento PIX - ID: 1
BACEN: Pagamento PIX autorizado - ID: 1
PIX: Pagamento liquidado - ID: 1, Status: SETTLED
BACEN: Processando liquidação de pagamento PIX - ID: 1
BACEN: Pagamento PIX liquidado - ID: 1, Valor transferido: R$ 123.45
PIX: Fluxo completo finalizado - ID: 1, Status: SETTLED
```

## 📝 Notificações Criadas

Para cada pagamento, são criadas **3 notificações**:

1. **PAYMENT_CREATED** - "Pagamento PIX criado com sucesso"
2. **PAYMENT_AUTHORIZED** - "Pagamento PIX autorizado pelo BACEN"
3. **PAYMENT_SETTLED** - "Pagamento PIX liquidado com sucesso"

## 🧪 Como Testar

### Via Swagger UI

1. Acesse: `http://localhost:8080/swagger/index.html`
2. Use o endpoint `POST /payments/pix`
3. Envie: `{"amount": 123.45}`
4. A resposta retornará imediatamente com `status: "CREATED"`
5. O fluxo completo (AUTHORIZED → SETTLED) acontece em background (~6 segundos)
6. Use o monitor em tempo real: `http://localhost:8080/monitor` para ver as mudanças de status

### Via cURL

```bash
curl -X POST http://localhost:8080/payments/pix \
  -H 'Content-Type: application/json' \
  -d '{"amount": 123.45}'
```

### Verificar Status

```bash
# Listar todos os pagamentos
curl http://localhost:8080/payments/pix

# Buscar por ID
curl http://localhost:8080/payments/pix/1
```

### Verificar Notificações no Banco

```bash
docker exec -it monolith-db-1 psql -U fintech -d fintech -c \
  "SELECT id, type, message, status, created_at FROM notifications ORDER BY created_at DESC LIMIT 10;"
```

## 🎓 Conceitos Demonstrados

### No Monólito

- ✅ **Comunicação direta**: Tudo acontece no mesmo processo
- ✅ **Transações ACID**: Status atualizado atomicamente
- ✅ **Simplicidade**: Fácil de debugar e rastrear
- ✅ **Latência baixa**: Sem chamadas de rede entre componentes

### Diferenças para Microsserviços

No monólito:
- Tudo é síncrono e rápido
- Status atualizado imediatamente
- Logs fáceis de seguir

Em microsserviços (ver `../microservices/`):
- Comunicação via HTTP (mais lenta)
- Eventual consistency
- Logs distribuídos (mais complexo)

## 💡 Observações Importantes

1. **Simulação Realista**: O fluxo simula o processo real do PIX brasileiro
2. **Tempo de Processamento**: ~6 segundos para o fluxo completo (processado em background)
3. **Resposta Imediata**: O POST retorna imediatamente com status `CREATED`
4. **Processamento Assíncrono**: Autorização e liquidação acontecem em goroutine (background)
5. **Monitoramento em Tempo Real**: Use `http://localhost:8080/monitor` para ver mudanças de status via SSE
6. **Notificações**: Criadas automaticamente em cada etapa (CREATED, AUTHORIZED, SETTLED)
7. **Logs Detalhados**: Facilita entender o fluxo completo

## 🔄 Próximos Passos

Para ver como isso funciona em microsserviços:
- Veja `../microservices/` para a versão distribuída
- Compare a complexidade e latência
- Entenda os trade-offs de cada abordagem

