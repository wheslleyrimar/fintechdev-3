# Como Testar a Aplicação Monolítica

## 🚀 Iniciar a Aplicação

### Opção 1: Rodar em Background (Recomendado)

Para não travar o terminal, rode em background:

```bash
cd monolith
docker compose up -d --build
```

O flag `-d` (detached) roda os containers em background.

### Opção 2: Rodar em Foreground (Ver logs)

Se quiser ver os logs em tempo real:

```bash
cd monolith
docker compose up --build
```

Em outro terminal, você pode ver os logs:
```bash
docker compose logs -f monolith-api
```

## 📊 Verificar Status dos Containers

```bash
docker compose ps
```

Deve mostrar algo como:
```
NAME                    STATUS          PORTS
monolith-db-1           Up (healthy)    0.0.0.0:5432->5432/tcp
monolith-monolith-api-1 Up              0.0.0.0:8080->8080/tcp
```

## 🧪 Testar a API

### 1. Health Check

```bash
curl http://localhost:8080/health
```

Resposta esperada:
```json
{"status":"ok","type":"monolith"}
```

### 2. Listar Todos os Pagamentos (GET)

**No navegador:**
```
http://localhost:8080/payments/pix
```

**Com curl:**
```bash
curl http://localhost:8080/payments/pix | jq .
```

Resposta esperada (array de pagamentos):
```json
[
  {
    "id": 1,
    "amount": 123.45,
    "status": "CREATED",
    "created_at": "2024-01-15T10:30:00Z"
  },
  {
    "id": 2,
    "amount": 250.75,
    "status": "CREATED",
    "created_at": "2024-01-15T10:31:00Z"
  }
]
```

### 3. Criar Pagamento PIX (POST)

```bash
curl -X POST http://localhost:8080/payments/pix \
  -H 'Content-Type: application/json' \
  -d '{"amount": 123.45}' \
  | jq .
```

**Se não tiver `jq` instalado:**
```bash
curl -X POST http://localhost:8080/payments/pix \
  -H 'Content-Type: application/json' \
  -d '{"amount": 123.45}'
```

Resposta esperada:
```json
{
  "id": 1,
  "amount": 123.45,
  "status": "CREATED",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### 4. Buscar Pagamento por ID (GET)

**No navegador:**
```
http://localhost:8080/payments/pix/1
```

**Com curl:**
```bash
curl http://localhost:8080/payments/pix/1 | jq .
```

Resposta esperada:
```json
{
  "id": 1,
  "amount": 123.45,
  "status": "CREATED",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### 3. Criar Múltiplos Pagamentos

```bash
# Pagamento 1
curl -X POST http://localhost:8080/payments/pix \
  -H 'Content-Type: application/json' \
  -d '{"amount": 50.00}' | jq .

# Pagamento 2
curl -X POST http://localhost:8080/payments/pix \
  -H 'Content-Type: application/json' \
  -d '{"amount": 250.75}' | jq .

# Pagamento 3
curl -X POST http://localhost:8080/payments/pix \
  -H 'Content-Type: application/json' \
  -d '{"amount": 1000.00}' | jq .
```

## 📝 Ver Logs

### Ver todos os logs
```bash
docker compose logs
```

### Ver logs em tempo real
```bash
docker compose logs -f
```

### Ver apenas logs da API
```bash
docker compose logs -f monolith-api
```

### Ver apenas logs do banco
```bash
docker compose logs -f db
```

### Ver últimas 50 linhas
```bash
docker compose logs --tail=50 monolith-api
```

## 🗄️ Verificar Banco de Dados

### Conectar ao PostgreSQL

```bash
docker exec -it monolith-db-1 psql -U fintech -d fintech
```

### Consultas Úteis

```sql
-- Ver todos os pagamentos
SELECT * FROM pix_payments ORDER BY created_at DESC;

-- Ver todas as notificações
SELECT * FROM notifications ORDER BY created_at DESC;

-- Contar registros
SELECT 
  (SELECT COUNT(*) FROM pix_payments) as total_payments,
  (SELECT COUNT(*) FROM notifications) as total_notifications;

-- Ver pagamentos com detalhes
SELECT 
  id,
  amount,
  status,
  created_at,
  EXTRACT(EPOCH FROM (NOW() - created_at)) as seconds_ago
FROM pix_payments
ORDER BY created_at DESC
LIMIT 10;
```

### Sair do PostgreSQL
```sql
\q
```

## 🛑 Parar a Aplicação

### Parar containers (mantém dados)
```bash
docker compose stop
```

### Parar e remover containers (mantém volumes)
```bash
docker compose down
```

### Parar e remover tudo (incluindo volumes - APAGA DADOS)
```bash
docker compose down -v
```

## 🔍 Debug

### Ver logs de erro
```bash
docker compose logs monolith-api | grep ERROR
```

### Entrar no container da API
```bash
docker exec -it monolith-monolith-api-1 sh
```

### Entrar no container do banco
```bash
docker exec -it monolith-db-1 psql -U fintech -d fintech
```

### Verificar se a API está respondendo
```bash
curl -v http://localhost:8080/health
```

O `-v` mostra detalhes da requisição HTTP.

## 📋 Script de Teste Completo

Crie um arquivo `test.sh`:

```bash
#!/bin/bash

echo "=== Testando Aplicação Monolítica ==="
echo ""

echo "1. Health Check..."
curl -s http://localhost:8080/health | jq .
echo ""

echo "2. Criando pagamento PIX..."
PAYMENT_RESPONSE=$(curl -s -X POST http://localhost:8080/payments/pix \
  -H 'Content-Type: application/json' \
  -d '{"amount": 123.45}')

echo "$PAYMENT_RESPONSE" | jq .
echo ""

echo "3. Verificando banco de dados..."
docker exec monolith-db-1 psql -U fintech -d fintech -c "SELECT COUNT(*) as total_payments FROM pix_payments;"
docker exec monolith-db-1 psql -U fintech -d fintech -c "SELECT COUNT(*) as total_notifications FROM notifications;"

echo ""
echo "=== Testes concluídos ==="
```

Torne executável e execute:
```bash
chmod +x test.sh
./test.sh
```

## 💡 Dicas

1. **Instalar jq** (formatação JSON):
   ```bash
   # macOS
   brew install jq
   
   # Linux
   sudo apt-get install jq
   ```

2. **Usar Postman ou Insomnia** para testes mais visuais

3. **Monitorar logs em tempo real** enquanto testa:
   ```bash
   # Terminal 1
   docker compose logs -f monolith-api
   
   # Terminal 2
   curl -X POST http://localhost:8080/payments/pix ...
   ```

4. **Verificar se containers estão rodando**:
   ```bash
   docker compose ps
   ```

## ❌ Problemas Comuns

### Container não inicia
```bash
# Ver logs de erro
docker compose logs monolith-api

# Verificar se porta 8080 está livre
lsof -i :8080
```

### Erro de conexão com banco
```bash
# Verificar se o banco está saudável
docker compose ps db

# Ver logs do banco
docker compose logs db
```

### Não consegue ver resposta do POST
- Use `jq` para formatar: `curl ... | jq .`
- Verifique logs: `docker compose logs -f monolith-api`
- Use `-v` no curl para ver detalhes: `curl -v ...`

