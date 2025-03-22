#!/bin/sh

# Gerar documentação Swagger
echo "Gerando documentação Swagger..."
swag init -g cmd/api/main.go --output ./docs --parseDependency true

# Compilar a aplicação pela primeira vez
go build -o /tmp/app ./cmd/api/main.go

# Iniciar a aplicação
/tmp/app &
APP_PID=$!

echo "Servidor iniciado na porta 8080"
echo "Documentação Swagger disponível em: http://localhost:8080/swagger/index.html"

# Monitorar alterações nos arquivos .go
while true; do
    if find . -name "*.go" -newer /tmp/app | grep -q .; then
        echo "Alterações detectadas, reiniciando aplicação..."
        kill $APP_PID
        
        # Regenerar documentação Swagger
        echo "Atualizando documentação Swagger..."
        swag init -g cmd/api/main.go --output ./docs --parseDependency true
        
        go build -o /tmp/app ./cmd/api/main.go
        /tmp/app &
        APP_PID=$!
        echo "Aplicação reiniciada!"
    fi
    sleep 2
done