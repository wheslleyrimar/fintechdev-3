# Comparação: Monólito vs Microsserviços

##  Visão Geral

Este documento compara as duas implementações disponíveis neste repositório.

##  Arquitetura

### Monólito

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

### Microsserviços

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

##  Comparação Detalhada

| Aspecto | Monólito | Microsserviços |
|---------|----------|----------------|
| **Banco de Dados** | Compartilhado | Por serviço |
| **Comunicação** | Chamada direta in-memory (repository) | HTTP/REST (cliente HTTP) |
| **Deploy** | Único para tudo | Independente por serviço |
| **Escalabilidade** | Tudo junto | Por serviço |
| **Complexidade** | Baixa | Alta |
| **Latência** | Baixa (chamadas locais) | Média (chamadas de rede) |
| **Consistência** | Forte (ACID) | Eventual |
| **Acoplamento** | Alto | Baixo |
| **Observabilidade** | SSE em tempo real | SSE em tempo real (necessária) |
| **Documentação API** | Swagger UI integrado | Endpoints REST básicos |
| **Portas** | 8080 (único serviço) | 8081 (payments), 8082 (notifications) |
| **Testes** | Mais simples | Mais complexos |
| **Debug** | Mais fácil | Mais difícil |

##  Vantagens do Monólito

1. **Simplicidade**
   - Menos complexidade operacional
   - Mais fácil de entender
   - Menos pontos de falha

2. **Performance**
   - Chamadas locais são rápidas
   - Sem latência de rede
   - Transações ACID simples

3. **Desenvolvimento**
   - Mais rápido para desenvolver
   - Fácil de debugar
   - Testes mais simples

4. **Custo**
   - Menos infraestrutura
   - Menos operações
   - Mais eficiente em recursos

##  Vantagens dos Microsserviços

1. **Escalabilidade**
   - Escalar partes específicas
   - Não precisa escalar tudo
   - Otimização de custos

2. **Deploy Independente**
   - Deploy sem afetar outros serviços
   - Rollback independente
   - Releases mais frequentes

3. **Autonomia**
   - Times independentes
   - Tecnologias diferentes
   - Evolução independente

4. **Resiliência**
   - Falhas isoladas
   - Um serviço pode falhar sem afetar outros
   - Melhor isolamento

##  Desvantagens do Monólito

1. **Acoplamento**
   - Mudanças afetam tudo
   - Deploy único
   - Difícil escalar partes específicas

2. **Evolução**
   - Tecnologia única
   - Times grandes bloqueiam uns aos outros
   - Difícil de dividir responsabilidades

3. **Escalabilidade**
   - Precisa escalar tudo junto
   - Não pode otimizar partes específicas
   - Custo cresce linearmente

##  Desvantagens dos Microsserviços

1. **Complexidade**
   - Mais complexidade operacional
   - Mais pontos de falha
   - Mais difícil de debugar

2. **Latência**
   - Chamadas de rede são mais lentas
   - Múltiplas chamadas aumentam latência
   - Timeout e retries necessários

3. **Consistência**
   - Eventual consistency
   - Transações distribuídas são complexas
   - Sincronização de dados

4. **Custo**
   - Mais infraestrutura
   - Mais operações
   - Observabilidade obrigatória

##  Quando Usar Cada Um?

### Use Monólito Quando:

-  Time pequeno (< 10 pessoas)
-  Produto ainda instável
-  Falta de testes automatizados
-  Deploy e monitoramento imaturos
-  Não há necessidade de escalar partes específicas
-  Times não bloqueiam uns aos outros

### Use Microsserviços Quando:

-  Times grandes bloqueando entregas
-  Domínios com ritmos diferentes de mudança
-  Necessidade de escalar partes específicas
-  Autonomia tecnológica como requisito
-  Custo de cloud começa a crescer de forma relevante
-  Maturidade técnica e organizacional

##  Exemplo Prático

### Cenário: Criar Pagamento

#### Monólito

```go
// Use case recebe ambos os repositórios diretamente
createUC := NewCreatePixPaymentUseCase(
    paymentRepo,      // Repositório de pagamentos
    notificationRepo,  // Repositório de notificações (mesmo banco)
    gateway,
    eventBroadcaster,
)

// No use case, chamada direta in-memory
payment := paymentRepo.Save(newPayment)
notification := notificationRepo.Save(newNotification)
// Ambos usam o mesmo pool de conexão (transação ACID simples)
```

**Características:**
- Chamada direta in-memory (rápida, ~nanossegundos)
- Transação única no mesmo banco
- Consistência forte (ACID)
- Simples de debugar (tudo no mesmo processo)
- Swagger UI disponível em `/swagger/`

#### Microsserviços

```go
// Payments Service - use case recebe cliente HTTP
createUC := NewCreatePixPaymentUseCase(
    paymentRepo,           // Repositório próprio (banco próprio)
    notificationClient,    // Cliente HTTP para notificações
    gateway,
    eventBroadcaster,
)

// No use case, comunicação via HTTP
payment := paymentRepo.Save(newPayment) // Salva no banco próprio

// Chama Notifications Service via HTTP (pode falhar)
err := notificationClient.SendPaymentCreatedNotification(
    payment.ID, 
    payment.Amount,
)
// Se falhar, o pagamento já foi criado (eventual consistency)
```

**Características:**
- Chamada HTTP (mais lenta, ~milissegundos)
- Transações separadas em bancos diferentes
- Eventual consistency (notificação pode falhar)
- Mais complexo de debugar (requer logs distribuídos)
- Timeout de 5 segundos configurado
- Falhas são tratadas silenciosamente (erro logado)

##  Migração

### Estratégia Recomendada

1. **Comece com Monólito**
   - Mais simples
   - Mais rápido
   - Menos risco

2. **Evolua quando Necessário**
   - Extraia funcionalidades gradualmente
   - Use Strangler Fig Pattern
   - Mantenha compatibilidade

3. **Não Force Migração**
   - Se o monólito funciona, mantenha
   - Migre apenas quando houver necessidade real
   - Evite migração por moda

##  Conclusão

**Não há resposta única.** A escolha depende de:

- Tamanho do time
- Maturidade técnica
- Necessidades de negócio
- Complexidade do domínio
- Recursos disponíveis

**Regra de ouro:** Comece simples (monólito) e evolua quando necessário (microsserviços).

##  Detalhes de Implementação

### Observabilidade

Ambas as implementações incluem **observabilidade em tempo real** usando **Server-Sent Events (SSE)**:

- **Monólito**: `http://localhost:8080/monitor/{paymentID}`
- **Microsserviços**: `http://localhost:8081/monitor/{paymentID}`

Isso permite monitorar mudanças de status de pagamentos em tempo real sem polling.

### Endpoints Disponíveis

#### Monólito (porta 8080)
- `POST /pix` - Criar pagamento
- `GET /pix` - Listar pagamentos
- `GET /pix/{id}` - Buscar pagamento
- `GET /pix/monitor/{id}` - Monitor SSE
- `GET /swagger/` - Documentação Swagger
- `GET /health` - Health check

#### Microsserviços
**Payments Service (porta 8081)**
- `POST /pix` - Criar pagamento
- `GET /pix` - Listar pagamentos
- `GET /pix/{id}` - Buscar pagamento
- `GET /pix/monitor/{id}` - Monitor SSE
- `GET /health` - Health check

**Notifications Service (porta 8082)**
- `POST /notifications` - Criar notificação (chamado internamente)
- `GET /notifications` - Listar notificações
- `GET /notifications/{id}` - Buscar notificação
- `GET /health` - Health check

##  Próximos Passos

1. Execute ambos os exemplos
2. Compare performance
3. Analise complexidade
4. Teste a observabilidade em tempo real (SSE)
5. Decida qual faz sentido para seu contexto
