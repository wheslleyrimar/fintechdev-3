# Guia do Instrutor - Aula 3

##  Visão Geral

**Duração:** 2-3 horas  
**Objetivo:** Ensinar estratégias práticas de evolução de monólito para microsserviços

##  Timing Sugerido

### Parte 1: Fundamentos (45 min)

- **O Mito dos Microsserviços** (10 min)
  - Microsserviços não são objetivo
  - São resposta a problemas específicos
  - Aumentam complexidade

- **Por que Microsserviços Falham** (10 min)
  - Reescrita total
  - Serviços pequenos demais
  - Falta de observabilidade

- **Quando Migrar / Quando NÃO Migrar** (15 min)
  - Critérios claros
  - Sinais de alerta
  - Exercício prático

- **O que são Microsserviços** (10 min)
  - Princípios fundamentais
  - Bounded Context vs Microsserviço

### Parte 2: Estratégias Práticas (60 min)

- **Migração Incremental** (15 min)
  - Strangler Fig Pattern
  - O que extrair primeiro

- **Autonomia de Dados** (20 min)
  - O maior erro: banco compartilhado
  - Banco por serviço
  - Demonstração prática

- **Comunicação entre Serviços** (15 min)
  - Síncrono vs Assíncrono
  - Quando usar cada um

- **Complexidade Distribuída** (10 min)
  - Desafios e soluções
  - Observabilidade

### Parte 3: Prática (45 min)

- **Demonstração do Código** (20 min)
  - Monólito vs Microsserviços
  - Comparação lado a lado

- **Exercícios Práticos** (25 min)
  - Identificar candidatos a extração
  - Design de comunicação
  - Implementação básica

##  Pontos-Chave

### 1. Microsserviços não são objetivo

**Enfatize:**
- Microsserviços são uma **solução**, não um objetivo
- Aumentam complexidade operacional
- Exigem maturidade técnica e organizacional

**Exemplo:**
> "Se você tem um time de 5 pessoas e um produto instável, microsserviços vão atrapalhar mais do que ajudar"

### 2. Autonomia de dados é essencial

**Enfatize:**
- Sem autonomia de dados, não há microsserviço real
- Banco compartilhado quebra independência
- É o erro mais comum e mais caro

**Demonstração:**
- Mostre o código do monólito (banco compartilhado)
- Mostre o código dos microsserviços (banco por serviço)
- Explique as diferenças práticas

### 3. Migração incremental

**Enfatize:**
- Evitar big bang
- Evolução contínua
- Decisões reversíveis

**Exemplo:**
- Mostre o Strangler Fig Pattern
- Explique como extrair funcionalidades gradualmente
- Demonstre com o código de exemplo

### 4. Começar simples

**Enfatize:**
- Síncrono primeiro
- Assíncrono quando necessário
- Não complicar antes do tempo

**Demonstração:**
- Mostre comunicação HTTP simples
- Explique quando considerar eventos
- Não implemente coreografia complexa cedo

##  Dicas de Ensino

### Use Comparações Visuais

```
Monólito:
┌─────────────────┐
│  Tudo junto     │
└─────────────────┘

Microsserviços:
┌──────┐  ┌──────┐  ┌──────┐
│  A   │  │  B   │  │  C   │
└──────┘  └──────┘  └──────┘
```

### Demonstre com Código

1. **Mostre o monólito primeiro** (`../monolith/`)

   **Banco compartilhado:**
   ```bash
   # Mostre o schema
   cat ../db/init-monolith.sql
   # Destaque: ambas as tabelas no mesmo banco, com foreign key
   ```

   **Chamadas diretas:**
   ```bash
   # Mostre o use case
   cat ../monolith/domains/payments/application/create_pix_payment_usecase.go
   # Linha 75-81: Cria notificação diretamente (in-memory)
   ```

   **Execute e teste:**
   ```bash
   cd ../monolith
   docker compose up --build
   curl -X POST http://localhost:8080/payments/pix \
     -H 'Content-Type: application/json' \
     -d '{"amount": 123.45}'
   ```

2. **Depois mostre microsserviços** (`../microservices/`)

   **Banco por serviço:**
   ```bash
   # Mostre os schemas separados
   cat ../db/init-payments.sql
   cat ../db/init-notifications.sql
   # Destaque: cada serviço tem seu próprio banco, sem foreign keys
   ```

   **Comunicação HTTP:**
   ```bash
   # Mostre o use case
   cat ../microservices/payments-service/application/create_pix_payment_usecase.go
   # Linha 74: Chama notificationClient via HTTP
   
   # Mostre o cliente HTTP
   cat ../microservices/payments-service/infra/notifications/http_notification_client.go
   ```

   **Execute e teste:**
   ```bash
   cd ../microservices
   docker compose up --build
   curl -X POST http://localhost:8081/pix \
     -H 'Content-Type: application/json' \
     -d '{"amount": 123.45}'
   ```

3. **Compare lado a lado**

   **Crie uma tabela comparativa:**
   - O que mudou? (comunicação, banco, deploy)
   - O que ficou igual? (lógica de negócio, domínios)
   - Quais são os trade-offs? (simplicidade vs complexidade)

   **Demonstre falha:**
   ```bash
   # Pare o notifications-service
   docker compose stop notifications-service
   
   # Crie um pagamento
   curl -X POST http://localhost:8081/pix \
     -H 'Content-Type: application/json' \
     -d '{"amount": 456.78}'
   
   # Mostre que o pagamento foi criado mesmo com notificação falhando
   # Isso demonstra eventual consistency
   ```

### Use Casos Reais

- **Quando migrar:** Times grandes, ritmos diferentes, necessidade de escalar
- **Quando NÃO migrar:** Time pequeno, produto instável, falta de testes

### Exercícios Práticos

1. **Identificar candidatos** - Use um sistema real
2. **Design de comunicação** - Desenhe o fluxo
3. **Implementação básica** - Extraia um serviço simples

##  Armadilhas Comuns

### 1. "Microsserviços são sempre melhores"

**Correção:**
- Nem sempre são a melhor solução
- Monólito bem projetado pode ser suficiente
- Aumentam complexidade

### 2. "Posso compartilhar banco de dados"

**Correção:**
- Banco compartilhado quebra autonomia
- É o maior erro possível
- Sem autonomia de dados, não há microsserviço

### 3. "Vou reescrever tudo de uma vez"

**Correção:**
- Big bang é perigoso
- Migração incremental é mais segura
- Strangler Fig Pattern

### 4. "Vou usar eventos desde o início"

**Correção:**
- Comece simples (HTTP)
- Eventos quando necessário
- Não complique antes do tempo

##  Recursos Adicionais

### Durante a Aula

- **Código de exemplo no repositório:**
  - `../monolith/` - Implementação monolítica
  - `../microservices/` - Implementação distribuída
- **Estratégias detalhadas:** `estrategias-migracao.md`
- **Exercícios práticos:** `exercicios.md`
- **Comparação:** `../COMPARACAO.md`

### Após a Aula

- **Exercícios:** `exercicios.md` (com referências específicas ao código)
- **Estratégias detalhadas:** `estrategias-migracao.md`
- **Código para experimentação:**
  - `../monolith/` - Execute e explore o monólito
  - `../microservices/` - Execute e explore os microsserviços
- **Documentação completa:**
  - `../README.md` - Visão geral do repositório
  - `../COMPARACAO.md` - Comparação detalhada
  - `../monolith/README.md` - Documentação do monólito
  - `../microservices/README.md` - Documentação dos microsserviços

##  Checklist de Preparação

Antes da aula, certifique-se de:

- [ ] Ter o repositório clonado
- [ ] Docker instalado e funcionando
- [ ] Código de exemplo testado (monólito e microsserviços)
- [ ] Slides revisados
- [ ] Exercícios preparados
- [ ] Ambiente de demonstração pronto

### Teste os Projetos Antes da Aula

```bash
# Teste o monólito
cd monolith
docker compose up --build
curl http://localhost:8080/health
curl -X POST http://localhost:8080/payments/pix \
  -H 'Content-Type: application/json' \
  -d '{"amount": 123.45}'
docker compose down

# Teste os microsserviços
cd ../microservices
docker compose up --build
curl http://localhost:8081/health
curl http://localhost:8082/health
curl -X POST http://localhost:8081/pix \
  -H 'Content-Type: application/json' \
  -d '{"amount": 123.45}'
docker compose down
```

### Arquivos Importantes para Demonstração

- `../db/init-monolith.sql` - Schema do banco compartilhado
- `../db/init-payments.sql` - Schema do banco do payments-service
- `../db/init-notifications.sql` - Schema do banco do notifications-service
- `../monolith/domains/payments/application/create_pix_payment_usecase.go` - Use case do monólito
- `../microservices/payments-service/application/create_pix_payment_usecase.go` - Use case dos microsserviços
- `../microservices/payments-service/infra/notifications/http_notification_client.go` - Cliente HTTP

##  Avaliação

### Perguntas para Verificar Aprendizado

1. Quando faz sentido migrar para microsserviços?
2. Qual é o maior erro ao criar microsserviços?
3. Como funciona o Strangler Fig Pattern?
4. Por que autonomia de dados é essencial?
5. Quando usar comunicação síncrona vs assíncrona?

### Exercícios de Avaliação

1. Identificar candidatos a extração em um sistema real
2. Desenhar arquitetura de migração
3. Implementar extração de um serviço simples

##  Perguntas Frequentes

### "Posso ter microsserviços com banco compartilhado?"

**Resposta:** Tecnicamente sim, mas você perde os principais benefícios. Sem autonomia de dados, não há microsserviço real.

### "Quantos serviços devo ter?"

**Resposta:** Não há número mágico. Comece com poucos e extraia quando necessário. Serviços pequenos demais são um problema comum.

### "Como lidar com transações distribuídas?"

**Resposta:** Evite quando possível. Use eventual consistency e padrões como Saga quando necessário.

### "Quando usar eventos vs HTTP?"

**Resposta:** Comece com HTTP (síncrono). Use eventos quando precisar de desacoplamento ou alta disponibilidade.

##  Mensagens Finais

Encerre a aula reforçando:

1. **Microsserviços são uma jornada** - Não é um destino
2. **Comece pelo design** - Não pela tecnologia
3. **Evolua incrementalmente** - Evite big bang
4. **Arquitetura serve ao negócio** - Não o contrário
5. **Bounded Context ≠ Microsserviço** - Entenda a diferença
