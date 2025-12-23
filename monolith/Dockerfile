FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

# Copiar código fonte
COPY . .

# Gerar documentação Swagger (se necessário)
# RUN go install github.com/swaggo/swag/cmd/swag@latest
# RUN swag init -g apps/monolith-api/main.go -o apps/monolith-api/docs --parseDependency --parseInternal

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/monolith-api ./apps/monolith-api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/monolith-api .

EXPOSE 8080
CMD ["./monolith-api"]
