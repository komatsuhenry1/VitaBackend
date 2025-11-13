# Estágio 1: Builder
# Usamos uma imagem Go com Alpine para um build leve
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copia e baixa as dependências PRIMEIRO
# Isso aproveita o cache do Docker se o go.mod/go.sum não mudarem
COPY go.mod go.sum ./
RUN go mod download

# Copia o resto do código-fonte
COPY . .

# Compila o binário estático para Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# ---

# Estágio 2: Final
# Começa de uma imagem Alpine limpa e minúscula
FROM alpine:latest

WORKDIR /app

# Copia APENAS o binário compilado do estágio 'builder'
COPY --from=builder /app/main .

# 💡 MUDANÇA AQUI:
# Expõe a porta 8081, que é a porta do seu .env (SERVER_PORT=8081)
EXPOSE 8081

# Comando para executar a API
CMD ["./main"]