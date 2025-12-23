# Fintech Dev - Aula 3: Estratégias de Evolução para Microsserviços

## 📋 Sobre a Aula

Esta aula demonstra a evolução de um **monólito saudável** para uma **arquitetura distribuída**, focando em estratégias práticas e seguras de migração.

### Objetivos da Aula

- ✅ Entender quando microsserviços fazem sentido
- ✅ Aprender estratégias seguras de evolução
- ✅ Evitar armadilhas comuns de migração
- ✅ Conectar decisões técnicas ao negócio

### Tópicos Abordados

1. **O Mito dos Microsserviços** - Microsserviços não são um objetivo, são uma resposta a problemas específicos
2. **Quando Migrar** - Critérios claros para decisão
3. **Quando NÃO Migrar** - Sinais de alerta
4. **O que são Microsserviços** - Princípios fundamentais
5. **Bounded Context vs Microsserviço** - Diferenças importantes
6. **Estratégias de Migração Incremental** - Strangler Fig Pattern
7. **Dados e Comunicação** - Autonomia de dados e padrões de comunicação
8. **Complexidade Distribuída** - Desafios e boas práticas

## 🏗️ Estrutura do Repositório

Este repositório contém **duas implementações** que demonstram a evolução:

### 1. Monólito (`monolith/`)
Implementação monolítica inicial com todos os domínios (Payments, Notifications) no mesmo serviço e banco de dados compartilhado.

**Características:**
- Banco de dados compartilhado
- Comunicação direta (in-memory)
- Deploy único
- Simples e rápido

### 2. Microsserviços (`microservices/`)
Implementação distribuída com:
- **payments-service**: Serviço de pagamentos
- **notifications-service**: Serviço de notificações
- Cada serviço com seu próprio banco de dados (autonomia de dados)
- Comunicação via HTTP (síncrona)

**Características:**
- Autonomia de dados (banco por serviço)
- Comunicação via HTTP
- Deploy independente
- Isolamento de falhas

> **Nota:** A migração incremental é demonstrada através da comparação entre o monólito e os microsserviços, seguindo o **Strangler Fig Pattern**.

## 🚀 Quick Start

### Executar o Monólito

```bash
cd monolith
docker compose up --build
```

Acesse: `http://localhost:8080/health`

### Executar Microsserviços

```bash
cd microservices
docker compose up --build
```

Serviços disponíveis:
- Payments: `http://localhost:8081/health`
- Notifications: `http://localhost:8082/health`

### Testar Criação de Pagamento

```bash
# Criar pagamento (monólito)
curl -X POST http://localhost:8080/payments/pix \
  -H 'Content-Type: application/json' \
  -d '{"amount": 123.45}'

# Criar pagamento (microsserviços)
curl -X POST http://localhost:8081/pix \
  -H 'Content-Type: application/json' \
  -d '{"amount": 123.45}'

# Listar pagamentos (microsserviços)
curl http://localhost:8081/pix

# Listar notificações (microsserviços)
curl http://localhost:8082/notifications
```

## 📚 Documentação

### Documentação Principal
- [`COMPARACAO.md`](COMPARACAO.md) - Comparação detalhada Monólito vs Microsserviços

### Documentação da Aula
- [`aula-3/README.md`](aula-3/README.md) - Documentação completa da aula
- [`aula-3/guia-instrutor.md`](aula-3/guia-instrutor.md) - Guia para instrutores
- [`aula-3/exercicios.md`](aula-3/exercicios.md) - Exercícios práticos
- [`aula-3/estrategias-migracao.md`](aula-3/estrategias-migracao.md) - Estratégias detalhadas

### Documentação dos Componentes
- [`monolith/README.md`](monolith/README.md) - Documentação do monólito
- [`microservices/README.md`](microservices/README.md) - Documentação dos microsserviços

## 🎯 Conceitos-Chave

### Princípios Fundamentais

- **Single Responsibility** por serviço
- **Independência de deploy**
- **Autonomia de dados**
- **Falhas isoladas**

### Quando Migrar

- Times grandes bloqueando entregas
- Domínios com ritmos diferentes de mudança
- Necessidade de escalar partes específicas
- Autonomia tecnológica como requisito
- Custo de cloud começa a crescer de forma relevante

### Quando NÃO Migrar

- Time pequeno ou pouco maduro
- Produto ainda instável
- Falta de testes automatizados
- Deploy e monitoramento imaturos
- Decisões difíceis no início que serão caras de mudar depois

### O Maior Erro: Dados

> **Sem autonomia de dados não existe microsserviço**

- Banco compartilhado quebra independência
- Acoplamento invisível
- Evolução bloqueada
- Banco como gargalo arquitetural

### Boas Práticas

- Banco por serviço
- Comunicação via API ou eventos
- Eventual consistency como padrão
- Comece simples (síncrono primeiro)
- Observabilidade obrigatória

## 🔄 Estratégias de Migração

### Strangler Fig Pattern

1. Novos serviços ao redor do monólito
2. Funcionalidades migradas aos poucos
3. Sistema legado vai sendo "estrangulado"

### O que Extrair Primeiro

- Alta taxa de mudança
- Fronteira clara de dados
- Baixo risco sistêmico
- Dependências externas
- Áreas de domínio bem definidas (logística, pagamentos, notificações, etc)

## 📖 Checklist de Decisão

Antes de criar um microsserviço, pergunte:

- ✅ Existe fronteira clara de domínio?
- ✅ O time consegue operar o serviço?
- ✅ O ganho compensa o custo?
- ✅ Observabilidade está pronta?
- ✅ Microsserviço ou apenas um bounded context mal definido?

## 💡 Mensagens Finais

- Microsserviços são uma jornada
- Comece pelo design, não pela tecnologia
- Evolua de forma incremental
- Arquitetura serve ao negócio
- Microservices vs Bounded Context não são a mesma coisa

## 🛠️ Tecnologias

- **Go** - Linguagem principal (pode ser estendido para outras)
- **PostgreSQL** - Banco de dados
- **Docker Compose** - Orquestração
- **HTTP/REST** - Comunicação síncrona
- **Eventos** - Comunicação assíncrona (exemplo)

## 📝 Licença

Este é um repositório educacional para fins de ensino.
