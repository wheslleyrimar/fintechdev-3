# Fintech Dev - Aula 3: Estratégias de Evolução para Microsserviços

## Sobre a Aula

Esta aula demonstra a evolução de um **monólito saudável** para uma **arquitetura distribuída**, focando em estratégias práticas e seguras de migração.

> **Para conteúdo completo da aula, consulte:** [`aula-3/README.md`](aula-3/README.md)

## Estrutura do Repositório

Este repositório contém **duas implementações** que demonstram a evolução:

### 1. Monólito (`monolith/`)
Implementação monolítica inicial com todos os domínios (Payments, Notifications) no mesmo serviço e banco de dados compartilhado.

> **Documentação completa:** [`monolith/README.md`](monolith/README.md)

### 2. Microsserviços (`microservices/`)
Implementação distribuída com payments-service e notifications-service, cada um com seu próprio banco de dados.

> **Documentação completa:** [`microservices/README.md`](microservices/README.md)

> **Comparação detalhada:** [`COMPARACAO.md`](COMPARACAO.md)

## Quick Start

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

> **Para exemplos completos de testes, consulte:**
> - [`monolith/COMO_TESTAR.md`](monolith/COMO_TESTAR.md) - Guia completo de testes do monólito
> - [`monolith/README.md`](monolith/README.md) - Endpoints do monólito
> - [`microservices/README.md`](microservices/README.md) - Endpoints dos microsserviços

## Documentação

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

## Conceitos-Chave

> **Para conceitos detalhados, estratégias e exemplos práticos, consulte:**
> - [`aula-3/README.md`](aula-3/README.md) - Conceitos-chave completos e agenda da aula
> - [`aula-3/estrategias-migracao.md`](aula-3/estrategias-migracao.md) - Estratégias detalhadas de migração
> - [`COMPARACAO.md`](COMPARACAO.md) - Comparação prática Monólito vs Microsserviços

## Tecnologias

- **Go** - Linguagem principal (pode ser estendido para outras)
- **PostgreSQL** - Banco de dados
- **Docker Compose** - Orquestração
- **HTTP/REST** - Comunicação síncrona
- **Eventos** - Comunicação assíncrona (exemplo)

## Licença

Este é um repositório educacional para fins de ensino.
