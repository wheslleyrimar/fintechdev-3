# Monólito - Implementação Inicial

## 📋 Sobre

Esta é a implementação **monolítica** inicial, onde todos os domínios (Payments e Notifications) estão no mesmo serviço e compartilham o mesmo banco de dados.

## 🏗️ Arquitetura

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

## ⚠️ Características

### Banco de Dados Compartilhado

- ✅ Simples de implementar
- ✅ Transações ACID entre domínios
- ✅ Queries que cruzam domínios são fáceis
- ❌ Quebra autonomia de dados
- ❌ Acoplamento invisível
- ❌ Evolução bloqueada

### Comunicação Direta

- ✅ Chamadas locais são rápidas
- ✅ Sem latência de rede
- ✅ Fácil de debugar
- ❌ Acoplamento forte
- ❌ Impossível escalar partes específicas
- ❌ Deploy único para tudo

## 🚀 Como Executar

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
**Monitor em Tempo Real:** `http://localhost:8080/monitor` 🆕

> 📖 **Guia completo de testes:** Veja [`COMO_TESTAR.md`](COMO_TESTAR.md) para instruções detalhadas  
> 📚 **Documentação Swagger:** Veja [`SWAGGER.md`](SWAGGER.md) para informações sobre a documentação da API  
> 🔍 **Observabilidade em Tempo Real:** Veja [`OBSERVABILIDADE.md`](OBSERVABILIDADE.md) para monitorar mudanças de status

## 📝 Endpoints Disponíveis

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
1. ✅ Cria pagamento (status: `CREATED`) - retorna imediatamente
2. ✅ Processa autorização no BACEN em background (~3s) (status: `AUTHORIZED`)
3. ✅ Processa liquidação em background (~6s total) (status: `SETTLED`)
4. ✅ Cria 3 notificações (criação, autorização, liquidação)
5. ✅ Use o monitor em tempo real (`http://localhost:8080/monitor`) para ver mudanças de status

> 📖 **Veja o fluxo completo:** [`FLUXO_PIX.md`](FLUXO_PIX.md)

### Buscar Pagamento por ID (GET)
```bash
# No navegador
http://localhost:8080/payments/pix/1

# Ou com curl
curl http://localhost:8080/payments/pix/1
```

## 🔄 Próximo Passo

Veja como este monólito evolui para microsserviços em `../microservices/`
