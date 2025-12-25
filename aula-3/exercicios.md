# Exercícios Práticos - Aula 3

##  Objetivos

Estes exercícios ajudam a consolidar os conceitos de migração de monólito para microsserviços.

##  Exercício 1: Identificar Candidatos a Extração

### Contexto

Você tem um monólito com os seguintes módulos (baseado no código deste repositório):
- **Pagamentos (PIX)** - Processa pagamentos PIX, autoriza e liquida transações
- **Notificações** - Envia notificações sobre eventos de pagamento (criação, autorização, liquidação)

> **Nota:** Este repositório demonstra a migração do módulo de Notificações do monólito para um microsserviço independente.

### Tarefa

Para cada módulo, avalie:
1. É um bom candidato para extração? Por quê?
2. Qual seria a ordem de prioridade?
3. Quais são os riscos?

### Critérios de Avaliação

- Alta taxa de mudança?
- Fronteira clara de dados?
- Baixo risco sistêmico?
- Dependências externas?
- Área de domínio bem definida?

### Exemplo Prático - Análise do Código

**1. Analise a estrutura do monólito:**
```bash
# Ver estrutura de domínios
ls -la ../monolith/domains/
# payments/ e notifications/ estão no mesmo serviço

# Ver banco de dados compartilhado
cat ../db/init-monolith.sql
# Ambas as tabelas (pix_payments e notifications) no mesmo banco
```

**2. Analise a estrutura dos microsserviços:**
```bash
# Ver serviços separados
ls -la ../microservices/
# payments-service/ e notifications-service/ são independentes

# Ver bancos de dados separados
cat ../db/init-payments.sql
cat ../db/init-notifications.sql
# Cada serviço tem seu próprio banco
```

**3. Por que Notificações foi extraído primeiro?**

No código deste repositório, **Notificações** foi extraído primeiro porque:

-  **Fronteira clara de dados**: Tabela `notifications` já estava separada (veja `db/init-monolith.sql`)
-  **Baixo risco sistêmico**: Falhas em notificações não afetam pagamentos (veja `microservices/payments-service/application/create_pix_payment_usecase.go` linha 74 - erro é ignorado)
-  **Dependências externas**: Em produção, integraria com serviços de email/SMS
-  **Área de domínio bem definida**: Responsabilidade única e clara

**4. Compare a comunicação:**

**Monólito** (`monolith/domains/payments/application/create_pix_payment_usecase.go:75-81`):
```go
// Comunicação direta (in-memory)
notification := notifications.NewNotification(
    saved.ID,
    "PAYMENT_CREATED",
    "user@example.com",
    "Pagamento PIX criado com sucesso",
)
_, _ = uc.notificationRepo.Save(notification)
```

**Microsserviços** (`microservices/payments-service/application/create_pix_payment_usecase.go:74`):
```go
// Comunicação via HTTP (pode falhar)
_ = uc.notificationClient.SendPaymentCreatedNotification(saved.ID, saved.Amount)
```

##  Exercício 2: Design de Comunicação

### Contexto

Você precisa extrair o serviço de Notificações do monólito. Este exercício analisa a implementação atual e propõe melhorias.

### Tarefa

1. **Analise a Implementação Atual (Síncrono HTTP)**

   **Veja o código atual:**
   ```bash
   # Cliente HTTP no payments-service
   cat ../microservices/payments-service/infra/notifications/http_notification_client.go
   
   # Handler no notifications-service
   cat ../microservices/notifications-service/api/notifications_handler.go
   ```

   **Perguntas:**
   - Como funciona o fluxo atual?
   - Onde estão os pontos de falha?
   - O que acontece se o notifications-service estiver offline?
   - Veja linha 74 de `create_pix_payment_usecase.go` - o erro é ignorado. Isso é correto?

2. **Proponha Melhorias**

   **Síncrono (HTTP) com Retry:**
   - Desenhe o fluxo de comunicação com retry
   - Identifique pontos de falha
   - Proponha estratégias de retry (exponential backoff)
   - Implemente circuit breaker

   **Assíncrono (Eventos):**
   - Desenhe o fluxo com eventos
   - Como garantir que a notificação será enviada? (at-least-once delivery)
   - Como lidar com falhas? (dead letter queue)
   - Como evitar duplicação? (idempotência)

3. **Comparação Prática**

   **Execute e teste:**
   ```bash
   # Terminal 1: Iniciar microsserviços
   cd ../microservices
   docker compose up
   
   # Terminal 2: Criar pagamento
   curl -X POST http://localhost:8081/pix \
     -H 'Content-Type: application/json' \
     -d '{"amount": 123.45}'
   
   # Terminal 3: Parar notifications-service (simular falha)
   docker compose stop notifications-service
   
   # Terminal 4: Criar outro pagamento e observar comportamento
   curl -X POST http://localhost:8081/pix \
     -H 'Content-Type: application/json' \
     -d '{"amount": 456.78}'
   
   # Ver logs do payments-service
   docker compose logs payments-service
   ```

   **Perguntas:**
   - Quando usar cada abordagem?
   - Quais são os trade-offs?
   - O que você observou quando o notifications-service estava offline?

##  Exercício 3: Estratégia de Migração

### Contexto

Você tem um monólito com 100.000 linhas de código e precisa migrar para microsserviços.

### Tarefa

Crie um plano de migração usando o Strangler Fig Pattern:

1. **Fase 1** - O que extrair primeiro? Por quê?
2. **Fase 2** - Como manter compatibilidade durante a migração?
3. **Fase 3** - Como garantir que o monólito e os serviços funcionem juntos?
4. **Fase 4** - Como validar que a migração foi bem-sucedida?

### Entregáveis

- Diagrama de arquitetura por fase
- Plano de rollback
- Critérios de sucesso

##  Exercício 4: Autonomia de Dados

### Contexto

Este exercício analisa como a autonomia de dados foi implementada na migração do monólito para microsserviços.

### Tarefa

1. **Análise de Dependências no Monólito**

   **Examine o banco compartilhado:**
   ```bash
   # Ver schema do monólito
   cat ../db/init-monolith.sql
   ```

   **Identifique:**
   - Tabelas usadas por cada módulo
   - Foreign keys entre módulos (linha 19: `CONSTRAINT fk_payment`)
   - Queries que cruzam módulos (se houver)

   **Perguntas:**
   - Por que a foreign key quebra autonomia?
   - O que acontece se você quiser mudar o schema de `pix_payments`?
   - Como isso afeta o deploy independente?

2. **Análise da Separação nos Microsserviços**

   **Examine os bancos separados:**
   ```bash
   # Banco do payments-service
   cat ../db/init-payments.sql
   
   # Banco do notifications-service
   cat ../db/init-notifications.sql
   ```

   **Compare:**
   - Não há foreign keys entre serviços
   - Cada serviço tem seu próprio banco
   - `payment_id` em `notifications` é apenas uma referência (não FK)

   **Perguntas:**
   - Como manter consistência sem foreign keys?
   - O que acontece se um pagamento for deletado?
   - Como lidar com dados compartilhados?

3. **Plano de Migração de Dados**

   **Cenário:** Você tem um monólito em produção com dados. Como migrar?

   **Tarefas:**
   - Como migrar dados existentes?
   - Como garantir integridade durante a migração?
   - Como fazer rollback se necessário?
   - Como lidar com dados que referenciam outros serviços?

   **Dica:** Pesquise sobre "Database Migration Strategies" e "Strangler Fig Pattern"

##  Exercício 5: Observabilidade

### Contexto

Em um sistema distribuído, é difícil rastrear requisições que passam por múltiplos serviços.

### Tarefa

1. **Correlação de Requisições**
   - Como rastrear uma requisição do início ao fim?
   - Como implementar correlation ID?
   - Como visualizar o fluxo completo?

2. **Logging**
   - O que deve ser logado em cada serviço?
   - Como estruturar os logs?
   - Como fazer busca eficiente?

3. **Monitoring**
   - Quais métricas são importantes?
   - Como detectar problemas rapidamente?
   - Como alertar o time responsável?

##  Exercício 6: Caso Real - Análise do Código

### Contexto

Este repositório demonstra uma fintech com um monólito que processa:
- Pagamentos PIX (com fluxo completo: CREATED → AUTHORIZED → SETTLED)
- Notificações sobre eventos de pagamento

O código já implementa a migração do serviço de Notificações do monólito para microsserviços.

### Tarefa

1. **Análise do Código Existente - Comparação Lado a Lado**

   **Use Case de Criação de Pagamento:**
   
   ```bash
   # Monólito: Comunicação direta
   cat ../monolith/domains/payments/application/create_pix_payment_usecase.go
   # Linha 75-81: Cria notificação diretamente
   
   # Microsserviços: Comunicação HTTP
   cat ../microservices/payments-service/application/create_pix_payment_usecase.go
   # Linha 74: Chama notificationClient via HTTP
   ```

   **Diferenças arquiteturais:**
   -  Comunicação: direta (in-memory) → HTTP
   -  Banco: compartilhado → banco por serviço
   -  Deploy: único → independente
   -  Acoplamento: alto → baixo

2. **Análise de Problemas - Quando Escalar**

   **Problemas do monólito em escala:**
   - Times grandes bloqueando entregas
   - Deploy único afeta tudo
   - Escalar tudo junto (custo alto)
   - Tecnologia única

   **Como microsserviços resolvem:**
   - Deploy independente
   - Escalar apenas o que precisa
   - Tecnologias diferentes por serviço
   - Times independentes

   **Novos problemas dos microsserviços:**
   - Complexidade operacional
   - Latência de rede
   - Eventual consistency
   - Debug mais difícil

3. **Proposta de Melhorias - O que Adicionar para Produção**

   **Analise o código atual:**
   ```bash
   # Cliente HTTP atual (sem retry)
   cat ../microservices/payments-service/infra/notifications/http_notification_client.go
   ```

   **Melhorias necessárias:**
   -  Retry com exponential backoff
   -  Circuit breaker
   -  Timeout apropriado
   -  Correlation ID para tracing
   -  Métricas e observabilidade
   -  Fallback strategy

4. **Teste Prático - Comparação de Comportamento**

   **Execute ambos os sistemas:**
   ```bash
   # Terminal 1: Monólito
   cd ../monolith
   docker compose up --build
   
   # Terminal 2: Microsserviços
   cd ../microservices
   docker compose up --build
   ```

   **Teste criação de pagamento:**
   ```bash
   # Monólito
   time curl -X POST http://localhost:8080/payments/pix \
     -H 'Content-Type: application/json' \
     -d '{"amount": 123.45}'
   
   # Microsserviços
   time curl -X POST http://localhost:8081/pix \
     -H 'Content-Type: application/json' \
     -d '{"amount": 123.45}'
   ```

   **Compare:**
   - Tempos de resposta
   - Logs de cada abordagem
   - Comportamento em falha (pare notifications-service)
   - Consistência dos dados

##  Exercício Prático: Analisar e Melhorar a Extração

### Objetivo

Analisar a extração já implementada do serviço de Notificações e propor melhorias.

> **Nota:** A extração já está implementada neste repositório. Este exercício foca em entender e melhorar a implementação existente.

### Passos

1. **Analisar a Implementação Existente**
   ```bash
   # Examinar o código do monólito
   cat monolith/domains/payments/application/create_pix_payment_usecase.go
   
   # Examinar o código dos microsserviços
   cat microservices/payments-service/application/create_pix_payment_usecase.go
   cat microservices/notifications-service/api/notifications_handler.go
   ```

2. **Executar e Comparar**
   ```bash
   # Terminal 1: Executar monólito
   cd monolith
   docker compose up --build
   
   # Terminal 2: Executar microsserviços
   cd microservices
   docker compose up --build
   
   # Terminal 3: Testar monólito
   curl -X POST http://localhost:8080/payments/pix \
     -H 'Content-Type: application/json' \
     -d '{"amount": 123.45}'
   
   # Listar pagamentos no monólito
   curl http://localhost:8080/payments/pix
   
   # Nota: No monólito, as notificações são criadas internamente
   # e não há endpoint público para listá-las (ver banco de dados)
   
   # Terminal 4: Testar microsserviços
   curl -X POST http://localhost:8081/pix \
     -H 'Content-Type: application/json' \
     -d '{"amount": 123.45}'
   
   # Verificar se notificação foi criada no notifications-service
   curl http://localhost:8082/notifications
   
   # Listar todos os pagamentos
   curl http://localhost:8081/pix
   
   # Buscar pagamento por ID
   curl http://localhost:8081/pix/1
   ```

3. **Identificar Diferenças**
   - Como a comunicação mudou?
   - O que mudou na estrutura de dados?
   - Quais são os trade-offs?

4. **Propor Melhorias**
   - Adicionar retry logic na comunicação HTTP
   - Implementar circuit breaker
   - Adicionar correlation ID
   - Melhorar tratamento de erros
   - Adicionar observabilidade (métricas, tracing)

5. **Testar Cenários de Falha**
   ```bash
   # Parar notifications-service e testar
   docker compose stop notifications-service
   
   # Criar pagamento e verificar comportamento
   curl -X POST http://localhost:8081/pix \
     -H 'Content-Type: application/json' \
     -d '{"amount": 123.45}'
   
   # Verificar logs do payments-service
   docker compose logs payments-service
   ```

### Critérios de Sucesso

- [ ] Entendeu as diferenças entre monólito e microsserviços
- [ ] Identificou pontos de melhoria na implementação atual
- [ ] Testou ambos os sistemas e comparou comportamentos
- [ ] Propos melhorias práticas e viáveis
- [ ] Entendeu os trade-offs de cada abordagem

##  Exercícios Adicionais

### Exercício 7: Bounded Context vs Microsserviço

Identifique em um sistema real:
- Quais são os bounded contexts?
- Quais deveriam ser microsserviços?
- Quais podem ficar no monólito?

### Exercício 8: Quando NÃO Migrar

Crie uma lista de critérios para decidir quando **NÃO** migrar para microsserviços.

### Exercício 9: Complexidade Distribuída

Implemente:
- Circuit breaker
- Retry com backoff
- Timeout apropriado
- Fallback strategy

##  Dicas

1. **Comece pequeno** - Extraia funcionalidades simples primeiro
2. **Teste bem** - Testes são ainda mais importantes em sistemas distribuídos
3. **Monitore tudo** - Observabilidade é essencial
4. **Documente decisões** - Por que você fez cada escolha?
5. **Aprenda com erros** - Migração é um processo de aprendizado

## Recursos

- Código de exemplo neste repositório
- Documentação em `estrategias-migracao.md`
- Guia do instrutor em `guia-instrutor.md`
- Documentação do monólito: `../monolith/README.md`
- Documentação dos microsserviços: `../microservices/README.md`
- Comparação detalhada: `../COMPARACAO.md`
