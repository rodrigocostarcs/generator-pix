FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar o código fonte
COPY . .

# Compilar o aplicativo
RUN CGO_ENABLED=0 GOOS=linux go build -o pix-generator ./cmd/api/main.go

# Imagem final
FROM alpine:latest

WORKDIR /app

# Copiar o binário compilado
COPY --from=builder /app/pix-generator .

# Porta que será exposta
EXPOSE 8080

# Comando para executar o aplicativo
CMD ["./pix-generator"]