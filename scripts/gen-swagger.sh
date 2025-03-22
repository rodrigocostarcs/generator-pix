#!/bin/sh

# Verificar se o swag está instalado
if ! command -v swag &> /dev/null; then
    echo "Erro: swag não está instalado. Execute 'go install github.com/swaggo/swag/cmd/swag@latest'"
    exit 1
fi

# Navegar para a raiz do projeto
cd "$(dirname "$0")/.."

# Remover documentos antigos
rm -rf ./docs

# Gerar nova documentação
swag init -g cmd/api/main.go --output ./docs --parseDependency true

echo "Documentação Swagger gerada com sucesso!"
echo "Acesse http://localhost:8080/swagger/index.html após iniciar o servidor para visualizar."