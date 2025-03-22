FROM golang:1.24-alpine

WORKDIR /app

# Instalar ferramentas necessárias
RUN apk add --no-cache git

# Instalar swag para geração de documentação Swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copiar arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Exposição da porta
EXPOSE 8080

# Script de inicialização para desenvolvimento
COPY scripts/dev-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/dev-entrypoint.sh

# Comando para executar em modo de desenvolvimento
CMD ["/usr/local/bin/dev-entrypoint.sh"]