# Monólito - Implementação Inicial

##  Sobre

Esta é a implementação **monolítica** inicial, onde todos os domínios (Payments e Notifications) estão no mesmo serviço e compartilham o mesmo banco de dados.

##  Arquitetura

```
┌─────────────────────────────────┐
│      Monolith API               │
│                                 │
│  ┌──────────┐  ┌─────────────┐ │
│  │ Payments │  │Notifications │ │
│  │ Domain   │  │   Domain     │ │
│  └──────────┘  └─────────────┘ │
│                                 │
└──────────────┬──────────────────┘
               │
               ▼
        ┌──────────┐
        │  Banco   │
        │Compartilhado│
        └──────────┘
```

##  Características

### Banco de Dados Compartilhado

-  Simples de implementar
-  Transações ACID entre domínios
-  Queries que cruzam domínios são fáceis
-  Quebra autonomia de dados
-  Acoplamento invisível
-  Evolução bloqueada

### Comunicação Direta

-  Chamadas locais são rápidas
-  Sem latência de rede
-  Fácil de debugar
-  Acoplamento forte
-  Impossível escalar partes específicas
-  Deploy único para tudo

##  Como Executar

### Rodar em Background (Recomendado - não trava o terminal)

```bash
docker compose up -d --build
```

### Rodar em Foreground (ver logs em tempo real)

```bash
docker compose up --build
```

Acesse: `http://localhost:8080/health`

**Swagger UI:** `http://localhost:8080/swagger/index.html`  
**Monitor em Tempo Real:** `http://localhost:8080/monitor` 

>  **Guia completo de testes:** Veja [`COMO_TESTAR.md`](COMO_TESTAR.md) para instruções detalhadas

##  Endpoints Disponíveis

### Listar Todos os Pagamentos (GET)
```bash
# No navegador
http://localhost:8080/payments/pix

# Ou com curl
curl http://localhost:8080/payments/pix
```

### Criar Pagamento PIX (POST) - Simula Fluxo Completo
```bash
curl -X POST http://localhost:8080/payments/pix \
  -H 'Content-Type: application/json' \
  -d '{"amount": 123.45}'
```

**O que acontece:**
1.  Cria pagamento (status: `CREATED`) - retorna imediatamente
2.  Processa autorização no BACEN em background (~3s) (status: `AUTHORIZED`)
3.  Processa liquidação em background (~6s total) (status: `SETTLED`)
4.  Cria 3 notificações (criação, autorização, liquidação)
5.  Use o monitor em tempo real (`http://localhost:8080/monitor`) para ver mudanças de status

## 📊 Fluxo Completo do PIX

O fluxo simula um pagamento PIX completo desde a criação até a liquidação:

### Timeline do Processo
```
0s     → Pagamento criado (CREATED) - retorna imediatamente
1s     → Delay inicial (para SSE conectar)
3s     → Pagamento autorizado (AUTHORIZED) - após 1s + 2s
6s     → Pagamento liquidado (SETTLED) - após 3s + 3s
```

**Total:** ~6 segundos para completar o fluxo completo (processado em background)

### Status do Pagamento

| Status | Descrição | Quando Ocorre |
|--------|-----------|---------------|
| `CREATED` | Pagamento criado | Imediatamente após criação (POST retorna) |
| `AUTHORIZED` | Autorizado pelo BACEN | Após ~3 segundos (1s delay + 2s) |
| `SETTLED` | Liquidado e finalizado | Após ~6 segundos (3s + 3s) |

### Notificações Criadas

Para cada pagamento, são criadas **3 notificações**:
1. **PAYMENT_CREATED** - "Pagamento PIX criado com sucesso"
2. **PAYMENT_AUTHORIZED** - "Pagamento PIX autorizado pelo BACEN"
3. **PAYMENT_SETTLED** - "Pagamento PIX liquidado com sucesso"

### Buscar Pagamento por ID (GET)
```bash
# No navegador
http://localhost:8080/payments/pix/1

# Ou com curl
curl http://localhost:8080/payments/pix/1
```

## 📖 Documentação Swagger/OpenAPI

A documentação Swagger está disponível em: **`http://localhost:8080/swagger/index.html`**

### Como Usar

1. Inicie a aplicação (veja seção "Como Executar" acima)
2. Acesse o Swagger UI no navegador
3. Teste os endpoints diretamente na interface

### Endpoints Documentados

- **GET** `/health` - Verifica se a API está funcionando
- **GET** `/payments/pix` - Lista todos os pagamentos PIX
- **POST** `/payments/pix` - Cria um novo pagamento PIX
- **GET** `/payments/pix/{id}` - Busca pagamento por ID

### Regenerar Documentação

Se você modificar os endpoints, regenere a documentação:

```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g apps/monolith-api/main.go -o apps/monolith-api/docs --parseDependency --parseInternal
```

## 👁️ Observabilidade em Tempo Real

A aplicação implementa **observabilidade em tempo real** usando **Server-Sent Events (SSE)** para monitorar mudanças de status de pagamentos PIX instantaneamente.

### Como Usar

1. **Acesse a página de monitoramento:**
   ```
   http://localhost:8080/monitor
   ```

2. **Crie um pagamento PIX** (em outro terminal):
   ```bash
   curl -X POST http://localhost:8080/payments/pix \
     -H 'Content-Type: application/json' \
     -d '{"amount": 123.45}'
   ```

3. **No monitor, digite o ID do pagamento** e clique em "Iniciar Monitoramento"

4. **Observe as mudanças em tempo real:**
   - 🟡 **CREATED** - Pagamento criado (imediato)
   - 🔵 **AUTHORIZED** - Autorizado pelo BACEN (após ~3 segundos)
   - 🟢 **SETTLED** - Liquidado (após ~6 segundos)

### Endpoints de Observabilidade

- **Página HTML:** `GET http://localhost:8080/monitor`
- **SSE Stream:** `GET http://localhost:8080/payments/pix/monitor/{id}`

### Formato dos Eventos SSE

```
event: initial
data: {"payment_id":1,"status":"CREATED","amount":123.45,"timestamp":"2024-01-15T10:30:00Z","message":"Status inicial do pagamento"}

event: status_change
data: {"payment_id":1,"status":"AUTHORIZED","amount":123.45,"timestamp":"2024-01-15T10:30:00.5Z","message":"Pagamento PIX autorizado pelo BACEN"}

event: status_change
data: {"payment_id":1,"status":"SETTLED","amount":123.45,"timestamp":"2024-01-15T10:30:01.5Z","message":"Pagamento PIX liquidado com sucesso"}
```

##  Próximo Passo

Veja como este monólito evolui para microsserviços em `../microservices/`
