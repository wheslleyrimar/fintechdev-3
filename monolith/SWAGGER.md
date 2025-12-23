# Swagger/OpenAPI Documentation

## 📚 Documentação da API

A documentação Swagger está disponível em:

**Swagger UI:** `http://localhost:8080/swagger/index.html`

## 🚀 Como Usar

### 1. Iniciar a Aplicação

```bash
cd monolith
docker compose up -d --build
```

### 2. Acessar o Swagger UI

Abra no navegador:
```
http://localhost:8080/swagger/index.html
```

### 3. Testar os Endpoints

O Swagger UI permite:
- ✅ Ver todos os endpoints disponíveis
- ✅ Ver exemplos de requisições e respostas
- ✅ Testar os endpoints diretamente na interface
- ✅ Ver os modelos de dados (schemas)

## 📝 Endpoints Documentados

### Health Check
- **GET** `/health` - Verifica se a API está funcionando

### Payments
- **GET** `/payments/pix` - Lista todos os pagamentos PIX
- **POST** `/payments/pix` - Cria um novo pagamento PIX
- **GET** `/payments/pix/{id}` - Busca pagamento por ID

## 🔄 Regenerar Documentação

Se você modificar os endpoints ou adicionar novos, regenere a documentação:

```bash
cd monolith
go run github.com/swaggo/swag/cmd/swag@latest init -g apps/monolith-api/main.go -o apps/monolith-api/docs --parseDependency --parseInternal
```

## 📋 Estrutura dos Arquivos

```
monolith/
├── apps/
│   └── monolith-api/
│       ├── docs/              # Documentação Swagger gerada
│       │   ├── docs.go        # Código Go gerado
│       │   ├── swagger.json   # Especificação OpenAPI (JSON)
│       │   └── swagger.yaml   # Especificação OpenAPI (YAML)
│       ├── http/
│       │   └── payments_facade.go  # Handlers com anotações Swagger
│       └── main.go            # Main com configuração Swagger
```

## 🎯 Anotações Swagger

As anotações Swagger são adicionadas nos handlers usando comentários especiais:

```go
// @Summary      Descrição curta
// @Description  Descrição detalhada
// @Tags         nome-da-tag
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do pagamento"
// @Success      200  {object}  payments.PixPayment
// @Failure      400  {object}  map[string]string
// @Router       /payments/pix/{id} [get]
func handler(w http.ResponseWriter, r *http.Request) {
    // ...
}
```

## 💡 Dicas

1. **Teste direto no Swagger UI**: Você pode executar requisições diretamente na interface
2. **Veja os modelos**: Clique em "Schemas" para ver a estrutura dos objetos
3. **Copie exemplos**: Use os exemplos de requisição/resposta como referência
4. **Exporte a especificação**: Baixe o `swagger.json` ou `swagger.yaml` para usar em outras ferramentas

## 🔗 Links Úteis

- [Swagger UI](http://localhost:8080/swagger/index.html)
- [Swagger JSON](http://localhost:8080/swagger/doc.json)
- [Documentação Swaggo](https://github.com/swaggo/swag)

