# Aula 3: Estratégias de Evolução para Microsserviços

##  Estrutura da Aula

Esta aula tem duração de **2-3 horas** e aborda como evoluir de um monólito saudável para uma arquitetura distribuída de forma segura e incremental.

###  Sobre Este Repositório

Este repositório contém **dois projetos práticos** que demonstram a evolução de monólito para microsserviços:

1. **`../monolith/`** - Implementação monolítica inicial
   - Todos os domínios (Payments, Notifications) no mesmo serviço
   - Banco de dados compartilhado (PostgreSQL)
   - Comunicação direta (in-memory)
   - Simples e rápido para desenvolver

2. **`../microservices/`** - Implementação distribuída
   - **payments-service**: Serviço de pagamentos PIX
   - **notifications-service**: Serviço de notificações
   - Cada serviço com seu próprio banco de dados (autonomia de dados)
   - Comunicação via HTTP (síncrona)

> ** Dica:** Use estes projetos como **laboratório prático** para entender as diferenças entre monólito e microsserviços.

### Arquivos da Aula

- `guia-instrutor.md` - Guia detalhado para o instrutor com timing e pontos-chave
- `exercicios.md` - Exercícios práticos para os alunos usando o código do repositório
- `estrategias-migracao.md` - Estratégias detalhadas de migração com exemplos do código

### Como Usar

1. **Guia do Instrutor**: Consulte `guia-instrutor.md` antes da aula para entender o timing e pontos-chave
2. **Exercícios**: Use `exercicios.md` para atividades práticas com o código do repositório
3. **Estratégias**: Consulte `estrategias-migracao.md` para entender os conceitos teóricos
4. **Código Prático**: 
   - Execute `../monolith/` para ver o monólito em ação
   - Execute `../microservices/` para ver os microsserviços em ação
   - Compare os dois para entender as diferenças práticas

### Pré-requisitos

- Conhecimento da Aula 2 (Monólito Bem Projetado)
- Familiaridade com conceitos de arquitetura de software
- Entendimento básico de DDD (Domain-Driven Design)
- Docker instalado (para executar os projetos)
- Acesso ao repositório para análise do código

###  Quick Start - Executar os Projetos

Antes de começar a aula, execute os projetos para ver as diferenças na prática:

**Executar o Monólito:**
```bash
cd ../monolith
docker compose up --build
# Acesse: http://localhost:8080/health
```

**Executar Microsserviços:**
```bash
cd ../microservices
docker compose up --build
# Payments: http://localhost:8081/health
# Notifications: http://localhost:8082/health
```

**Testar Criação de Pagamento:**
```bash
# Monólito
curl -X POST http://localhost:8080/payments/pix \
  -H 'Content-Type: application/json' \
  -d '{"amount": 123.45}'

# Microsserviços
curl -X POST http://localhost:8081/pix \
  -H 'Content-Type: application/json' \
  -d '{"amount": 123.45}'
```

>  **Documentação completa:** Veja `../README.md` para mais detalhes sobre como executar e testar os projetos.

### Objetivos de Aprendizado

Ao final desta aula, os alunos serão capazes de:

1. Entender quando microsserviços fazem sentido
2. Aprender estratégias seguras de evolução
3. Evitar armadilhas comuns de migração
4. Conectar decisões técnicas ao negócio
5. Aplicar o Strangler Fig Pattern
6. Entender a importância da autonomia de dados
7. Reconhecer a diferença entre Bounded Context e Microsserviço

##  Agenda da Aula

1. **O Mito dos Microsserviços** - Microsserviços não são um objetivo
2. **Por que microsserviços falham?** - Armadilhas comuns
3. **Quando migrar** - Critérios claros
4. **Quando NÃO migrar** - Sinais de alerta
5. **O que são Microsserviços** - Princípios fundamentais
6. **Princípios Fundamentais** - Single Responsibility, Independência, Autonomia
7. **Bounded Context vs Microsserviço** - Diferenças importantes
8. **Estilos de Implementação** - Domain Model, Transaction Script, Table Module
9. **Migração Incremental** - Estratégias seguras
10. **Strangler Fig Pattern** - Padrão prático
11. **O que extrair primeiro** - Critérios de priorização
12. **Exemplos Comuns de Extração** - Notificações, Antifraude, Relatórios
13. **O Maior Erro: Dados** - Autonomia de dados
14. **Boas Práticas de Dados** - Banco por serviço
15. **Comunicação entre Serviços** - Síncrono vs Assíncrono
16. **Complexidade Distribuída** - Desafios e soluções
17. **Arquitetura e Organização** - Conway's Law
18. **Checklist de Decisão** - Perguntas essenciais
19. **Mensagens Finais** - Resumo e próximos passos

##  Conceitos-Chave

### O Mito dos Microsserviços

- Microsserviços não são um objetivo
- São uma resposta a problemas específicos
- Aumentam complexidade operacional
- Exigem maturidade técnica e organizacional

> **Arquitetura não corrige problemas de processo**

### Por que Microsserviços Falham?

- Reescrita total do sistema
- Serviços pequenos demais
- Falta de observabilidade
- Times sem ownership claro

> **Distribuir código ruim não melhora arquitetura**

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

> **Toda decisão difícil no início é muito difícil de mudar depois**

### Princípios Fundamentais

- Single Responsibility por serviço
- Independência de deploy
- Autonomia de dados
- Falhas isoladas

### Bounded Context vs Microsserviço

- Bounded Context é conceito de domínio (DDD)
- Microsserviço é decisão arquitetural
- Nem todo bounded context vira microsserviço
- Confundir os dois gera serviços artificiais

> **Microsserviços são uma forma de organizar sistemas e times**

### O Maior Erro: Dados

- Banco compartilhado quebra independência
- Acoplamento invisível
- Evolução bloqueada
- Banco como gargalo arquitetural

> **Sem autonomia de dados não existe microsserviço**

### Boas Práticas de Dados

- Banco por serviço
- Comunicação via API ou eventos
- Eventual consistency como padrão

### Comunicação entre Serviços

- Comece simples
- Síncrono primeiro
- Assíncrono quando necessário
- Evitar coreografia complexa cedo
- Webhooks sempre devem ter fila para evitar timeout

### Complexidade Distribuída

- Latência
- Falhas parciais
- Retries e timeouts
- Observabilidade obrigatória
- Arquitetura restringe o design

### Arquitetura e Organização

- Conway's Law
- Times orientados a domínio
- Ownership claro
- Plataforma como habilitadora
- Big Ball of Mud surge quando o sistema se afasta do usuário

### Checklist de Decisão

- Existe fronteira clara de domínio?
- O time consegue operar o serviço?
- O ganho compensa o custo?
- Observabilidade está pronta?
- Microsserviço ou apenas um bounded context mal definido?

##  Mensagens Finais

- Microsserviços são uma jornada
- Comece pelo design, não pela tecnologia
- Evolua de forma incremental
- Arquitetura serve ao negócio
- Microservices vs Bounded Context não são a mesma coisa

##  Análise Prática: Comparando os Projetos

>  **Para análise detalhada e comparação completa, consulte:**
> - [`../COMPARACAO.md`](../COMPARACAO.md) - Comparação detalhada Monólito vs Microsserviços
> - [`../monolith/README.md`](../monolith/README.md) - Documentação do monólito
> - [`../microservices/README.md`](../microservices/README.md) - Documentação dos microsserviços
> - [`estrategias-migracao.md`](estrategias-migracao.md) - Exemplos práticos do código

##  Recursos Adicionais

- **Comparação detalhada:** [`../COMPARACAO.md`](../COMPARACAO.md)
- **Documentação do monólito:** [`../monolith/README.md`](../monolith/README.md)
- **Documentação dos microsserviços:** [`../microservices/README.md`](../microservices/README.md)
- **Como testar:** [`../monolith/COMO_TESTAR.md`](../monolith/COMO_TESTAR.md)
- **Estratégias de migração:** [`estrategias-migracao.md`](estrategias-migracao.md) - Exemplos práticos do código
