#!/bin/sh

# Compilar a aplicação pela primeira vez
go build -o /tmp/app ./cmd/api/main.go

# Iniciar a aplicação
/tmp/app &
APP_PID=$!

# Monitorar alterações nos arquivos .go
while true; do
    if find . -name "*.go" -newer /tmp/app | grep -q .; then
        echo "Alterações detectadas, reiniciando aplicação..."
        kill $APP_PID
        go build -o /tmp/app ./cmd/api/main.go
        /tmp/app &
        APP_PID=$!
    fi
    sleep 2
done