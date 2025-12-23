# Microsserviços - Arquitetura Distribuída

## 📋 Sobre

Esta é a implementação **distribuída** com microsserviços, onde cada serviço tem seu próprio banco de dados e se comunica via HTTP.

## 🏗️ Arquitetura

```
┌──────────────┐     HTTP     ┌──────────────┐
│ Payments     │─────────────▶│ Notifications│
│ Service      │              │   Service    │
└───────┬──────┘              └───────┬──────┘
        │                              │
        ▼                              ▼
┌──────────┐                    ┌──────────┐
│ Payments │                    │Notifications│
│   DB     │                    │    DB     │
└──────────┘                    └──────────┘
```

## ✅ Características

### Autonomia de Dados

- ✅ Cada serviço tem seu próprio banco
- ✅ Evolução independente
- ✅ Escalabilidade por serviço
- ✅ Isolamento de falhas

### Comunicação via HTTP

- ✅ Desacoplamento de serviços
- ✅ Deploy independente
- ✅ Escalabilidade independente
- ⚠️ Latência de rede
- ⚠️ Falhas parciais
- ⚠️ Eventual consistency

## 🚀 Como Executar

```bash
docker compose up --build
```

Serviços disponíveis:
- Payments: `http://localhost:8081/health`
- Notifications: `http://localhost:8082/health`

## 📝 Endpoints Disponíveis

### Payments Service (porta 8081)

```bash
# Criar pagamento PIX (inicia fluxo completo: CREATED -> AUTHORIZED -> SETTLED)
curl -X POST http://localhost:8081/pix \
  -H 'Content-Type: application/json' \
  -d '{"amount": 123.45}'

# Listar todos os pagamentos
curl http://localhost:8081/pix

# Buscar pagamento por ID
curl http://localhost:8081/pix/1

# Monitor em tempo real (SSE) - página HTML
# Acesse no navegador: http://localhost:8081/monitor

# Health check
curl http://localhost:8081/health
```

### Notifications Service (porta 8082)

```bash
# Listar todas as notificações
curl http://localhost:8082/notifications

# Buscar notificação por ID
curl http://localhost:8082/notifications/1

# Health check
curl http://localhost:8082/health
```

## 🔄 Fluxo de Comunicação Completo

O fluxo completo de pagamento PIX funciona da seguinte forma:

1. **Cliente** cria pagamento no Payments Service
2. **Payments Service** salva no seu banco com status `CREATED`
3. **Payments Service** chama Notifications Service via HTTP para notificar criação
4. **Payments Service** processa autorização (simula BACEN) e atualiza para `AUTHORIZED`
5. **Payments Service** chama Notifications Service via HTTP para notificar autorização
6. **Payments Service** processa liquidação (simula BACEN) e atualiza para `SETTLED`
7. **Payments Service** chama Notifications Service via HTTP para notificar liquidação

**Nota:** O fluxo completo acontece em background (goroutine), permitindo que a resposta retorne imediatamente com o pagamento criado.

### Tratamento de Falhas

Se o Notifications Service estiver indisponível:
- O pagamento já foi criado (eventual consistency)
- O fluxo de autorização e liquidação continua normalmente
- As notificações falharão silenciosamente (erro é logado mas não interrompe o fluxo)
- Em produção, implementar retry ou fila de eventos para garantir entrega das notificações

## 📊 Comparação com Monólito

| Aspecto | Monólito | Microsserviços |
|---------|----------|----------------|
| Banco de Dados | Compartilhado | Por serviço |
| Comunicação | Chamada direta | HTTP/Eventos |
| Deploy | Único | Independente |
| Escalabilidade | Tudo junto | Por serviço |
| Complexidade | Baixa | Alta |
| Latência | Baixa | Média |
| Consistência | Forte | Eventual |

## 🎯 Lições Aprendidas

1. **Autonomia de dados é essencial** - Sem isso, não há microsserviço real
2. **Comunicação síncrona é simples** - Comece por aqui
3. **Falhas são esperadas** - Projete para elas
4. **Observabilidade é obrigatória** - Sem ela, é impossível debugar

## 🔍 Monitoramento em Tempo Real

O payments-service inclui uma página de monitoramento em tempo real usando Server-Sent Events (SSE):

- **Página HTML:** `http://localhost:8081/monitor`
- **Endpoint SSE:** `http://localhost:8081/pix/monitor/{id}`

A página permite:
- Criar e monitorar pagamentos em tempo real
- Ver mudanças de status (CREATED → AUTHORIZED → SETTLED)
- Visualizar log de eventos com timestamps

## 📚 Próximos Passos

- Implementar comunicação assíncrona (eventos)
- Adicionar circuit breaker
- Implementar retry com backoff
- Adicionar observabilidade (tracing, métricas)
- Implementar API Gateway
