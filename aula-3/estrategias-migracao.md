# Estratégias de Migração para Microsserviços

##  Visão Geral

Este documento detalha as estratégias práticas para evoluir de um monólito para microsserviços de forma segura e incremental.

##  Princípios Fundamentais

### 1. Migração Incremental

-  Evitar big bang
-  Evolução contínua
-  Decisões reversíveis
-  Foco em risco controlado

### 2. O que Extrair Primeiro

Priorize extrair funcionalidades que tenham:

- **Alta taxa de mudança** - Funcionalidades que mudam frequentemente se beneficiam de deploy independente
- **Fronteira clara de dados** - Domínios com dados bem definidos são mais fáceis de extrair
- **Baixo risco sistêmico** - Comece por funcionalidades menos críticas
- **Dependências externas** - Integrações com sistemas externos são bons candidatos
- **Áreas de domínio bem definidas** - Logística, pagamentos, notificações, relatórios

### 3. Exemplos Comuns de Extração

- **Notificações** - Sistema de envio de emails/SMS
- **Antifraude** - Análise e detecção de fraudes
- **Relatórios** - Geração de relatórios e analytics
- **Integrações externas** - APIs de terceiros

##  Strangler Fig Pattern

O **Strangler Fig Pattern** é a estratégia mais segura para migração incremental.

### Como Funciona

1. **Novos serviços ao redor do monólito** - Criar novos serviços sem modificar o monólito
2. **Funcionalidades migradas aos poucos** - Extrair funcionalidades gradualmente
3. **Sistema legado vai sendo "estrangulado"** - O monólito diminui com o tempo

### Exemplo Prático - Este Repositório

Este repositório demonstra o **Strangler Fig Pattern** na prática:

**Fase 1: Monólito completo** (`../monolith/`)
```
┌─────────────────────────┐
│   Monólito              │
│  - Payments (PIX)       │
│  - Notifications        │
└─────────────────────────┘
         │
         ▼
    ┌──────────┐
    │  Banco   │
    │Compartilhado│
    └──────────┘
```

**Estrutura:**
- `monolith/domains/payments/` - Domínio de pagamentos
- `monolith/domains/notifications/` - Domínio de notificações
- `db/init-monolith.sql` - Banco compartilhado

**Fase 2: Extrair Notifications** (`../microservices/`)
```
┌─────────────────────────┐     HTTP     ┌──────────────┐
│   Monólito              │─────────────▶│ Notifications│
│  - Payments (PIX)       │              │   Service    │
└─────────────────────────┘              └───────┬──────┘
         │                                        │
         ▼                                        ▼
    ┌──────────┐                          ┌──────────┐
    │ Payments │                          │Notifications│
    │   DB     │                          │    DB     │
    └──────────┘                          └──────────┘
```

**Estrutura:**
- `microservices/payments-service/` - Serviço de pagamentos
- `microservices/notifications-service/` - Serviço de notificações extraído
- `db/init-payments.sql` - Banco do payments-service
- `db/init-notifications.sql` - Banco do notifications-service

**Como verificar no código:**

1. **Banco de dados:**
   ```bash
   # Monólito: banco compartilhado
   cat ../db/init-monolith.sql
   # Ambas as tabelas no mesmo banco, com foreign key
   
   # Microsserviços: bancos separados
   cat ../db/init-payments.sql
   cat ../db/init-notifications.sql
   # Cada serviço tem seu próprio banco, sem foreign keys
   ```

2. **Comunicação:**
   ```bash
   # Monólito: comunicação direta
   cat ../monolith/domains/payments/application/create_pix_payment_usecase.go
   # Linha 75-81: Cria notificação diretamente
   
   # Microsserviços: comunicação HTTP
   cat ../microservices/payments-service/application/create_pix_payment_usecase.go
   # Linha 74: Chama notificationClient via HTTP
   cat ../microservices/payments-service/infra/notifications/http_notification_client.go
   ```

## 🗄 Autonomia de Dados

### O Maior Erro: Banco Compartilhado

** ERRADO:**
```
┌──────────────┐     ┌──────────────┐
│ Payments      │     │ Notifications│
│ Service       │     │   Service    │
└───────┬───────┘     └───────┬───────┘
        │                     │
        └──────────┬──────────┘
                   ▼
            ┌──────────┐
            │  Banco   │
            │Compartilhado│
            └──────────┘
```

**Problemas:**
- Banco compartilhado quebra independência
- Acoplamento invisível
- Evolução bloqueada
- Banco como gargalo arquitetural

**Exemplo do que NÃO fazer:**
```sql
--  ERRADO: Foreign key entre serviços
CREATE TABLE notifications (
  payment_id BIGINT NOT NULL,
  CONSTRAINT fk_payment FOREIGN KEY (payment_id) 
    REFERENCES pix_payments(id)  -- Quebra autonomia!
);
```

**No monólito isso é aceitável** (`db/init-monolith.sql` linha 19), mas **nos microsserviços não**.

###  CORRETO: Banco por Serviço

```
┌──────────────┐     ┌──────────────┐
│ Payments      │     │ Notifications│
│ Service       │     │   Service    │
└───────┬───────┘     └───────┬───────┘
        │                     │
        ▼                     ▼
┌──────────┐           ┌──────────┐
│ Payments │           │Notifications│
│   DB     │           │    DB     │
└──────────┘           └──────────┘
```

**Benefícios:**
- Autonomia de dados
- Evolução independente
- Escalabilidade por serviço
- Isolamento de falhas

> **Sem autonomia de dados não existe microsserviço**

**Exemplo correto neste repositório:**

**Payments Service** (`db/init-payments.sql`):
```sql
CREATE TABLE pix_payments (
  id BIGSERIAL PRIMARY KEY,
  amount NUMERIC(18,2) NOT NULL,
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- Sem referências a outros serviços
```

**Notifications Service** (`db/init-notifications.sql`):
```sql
CREATE TABLE notifications (
  id BIGSERIAL PRIMARY KEY,
  payment_id BIGINT NOT NULL,  -- Apenas referência, não FK
  type TEXT NOT NULL,
  recipient TEXT NOT NULL,
  message TEXT NOT NULL,
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- payment_id é apenas uma referência, não foreign key
-- Cada serviço mantém autonomia sobre seus dados
```

## 🔗 Comunicação entre Serviços

### Síncrono vs Assíncrono

#### Síncrono (HTTP/REST)

**Quando usar:**
- Resposta imediata necessária
- Operações simples
- Baixa latência aceitável

**Exemplo - Este Repositório:**

**Cliente HTTP** (`microservices/payments-service/infra/notifications/http_notification_client.go`):
```go
func (c *HttpNotificationClient) SendPaymentCreatedNotification(
    paymentID int64, 
    amount float64,
) error {
    // Chama Notifications Service via HTTP
    url := fmt.Sprintf("%s/notifications", c.baseURL)
    // ... implementação HTTP
}
```

**Uso no Use Case** (`microservices/payments-service/application/create_pix_payment_usecase.go:74`):
```go
// Criar notificação de criação
_ = uc.notificationClient.SendPaymentCreatedNotification(saved.ID, saved.Amount)
```

**Vantagens:**
- Simples de implementar
- Fácil de debugar
- Resposta imediata

**Desvantagens:**
- Acoplamento temporal
- Falhas podem afetar o chamador
- Latência adicionada

#### Assíncrono (Eventos)

**Quando usar:**
- Operações que podem ser processadas depois
- Alta disponibilidade necessária
- Desacoplamento desejado

**Exemplo (não implementado neste repositório, mas seria assim):**
```go
// Payments Service publica evento
eventBus.Publish("payment.created", PaymentCreatedEvent{
    PaymentID: paymentID,
    Amount: amount,
})

// Notifications Service consome evento
eventBus.Subscribe("payment.created", handlePaymentCreated)
```

**Nota:** Este repositório usa comunicação síncrona (HTTP) para simplicidade. Em produção, você pode evoluir para eventos quando necessário.

**Vantagens:**
- Desacoplamento
- Resiliência a falhas
- Escalabilidade

**Desvantagens:**
- Complexidade maior
- Eventual consistency
- Debug mais difícil

### Boas Práticas

1. **Comece simples** - Síncrono primeiro
2. **Assíncrono quando necessário** - Quando precisar de desacoplamento
3. **Evitar coreografia complexa cedo** - Não complique antes do tempo
4. **Webhooks sempre devem ter fila** - Para evitar timeout

##  Eventual Consistency

### O que é?

Em uma arquitetura distribuída, nem sempre todos os serviços terão dados consistentes ao mesmo tempo.

### Exemplo - Este Repositório

**Fluxo Real** (`microservices/payments-service/application/create_pix_payment_usecase.go`):

```
1. Payments Service cria pagamento
   └─▶ Salva no banco próprio (linha 45)
   
2. Payments Service notifica Notifications Service via HTTP (linha 74)
   └─▶ Se falhar, o pagamento já foi criado
   └─▶ Erro é ignorado (linha 74: `_ = uc.notificationClient...`)
   
3. Notifications Service pode processar depois
   └─▶ Eventual consistency
```

**Teste prático:**
```bash
# Pare o notifications-service
docker compose stop notifications-service

# Crie um pagamento
curl -X POST http://localhost:8081/pix \
  -H 'Content-Type: application/json' \
  -d '{"amount": 123.45}'

# O pagamento foi criado mesmo com notificação falhando
# Isso demonstra eventual consistency
```

### Quando é Aceitável?

- Notificações podem ser enviadas depois
- Relatórios podem ter dados ligeiramente desatualizados
- Analytics não precisa ser em tempo real

### Quando NÃO é Aceitável?

- Saldo de conta bancária
- Status de pagamento crítico
- Operações financeiras que requerem consistência forte

## Complexidade Distribuída

### Desafios

1. **Latência** - Chamadas de rede são mais lentas que chamadas locais
2. **Falhas parciais** - Um serviço pode falhar enquanto outros funcionam
3. **Retries e timeouts** - Necessário lidar com falhas de rede
4. **Observabilidade obrigatória** - Difícil debugar sistemas distribuídos
5. **Arquitetura restringe o design** - Ex: Rails segue o Rails

### Soluções

1. **Circuit Breaker** - Evitar cascata de falhas
2. **Retry com backoff** - Tentar novamente com intervalo crescente
3. **Timeouts apropriados** - Não esperar indefinidamente
4. **Logging distribuído** - Correlação de requisições
5. **Monitoring e alerting** - Detectar problemas rapidamente

##  Checklist de Decisão

Antes de criar um microsserviço, pergunte:

- [ ] **Existe fronteira clara de domínio?**
  - O domínio tem responsabilidades bem definidas?
  - Os dados são independentes?

- [ ] **O time consegue operar o serviço?**
  - Tem conhecimento necessário?
  - Tem capacidade de deploy e monitoramento?

- [ ] **O ganho compensa o custo?**
  - O problema justifica a complexidade adicional?
  - Os benefícios superam os custos?

- [ ] **Observabilidade está pronta?**
  - Logging distribuído configurado?
  - Monitoring e alerting funcionando?
  - Tracing de requisições implementado?

- [ ] **Microsserviço ou apenas um bounded context mal definido?**
  - É realmente necessário separar?
  - Não seria melhor apenas melhorar o monólito?

## Lições Aprendidas

### O que Funciona

-  Migração incremental
-  Banco por serviço
-  Começar simples (síncrono)
-  Observabilidade desde o início
-  Times com ownership claro

### O que NÃO Funciona

-  Reescrita total
-  Banco compartilhado
-  Serviços pequenos demais
-  Falta de observabilidade
-  Times sem ownership

##  Referências

- "Building Microservices" - Sam Newman
- "Monolith to Microservices" - Sam Newman
- "Domain-Driven Design" - Eric Evans
- "Patterns of Enterprise Application Architecture" - Martin Fowler
