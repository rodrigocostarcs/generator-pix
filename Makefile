.PHONY: build run test docker-build docker-run swagger clean help

# Variáveis
APP_NAME=generator-pix
MAIN_PATH=./cmd/api/main.go
DOCKER_IMAGE_NAME=pix-generator
DOCKER_CONTAINER_NAME=pix-generator-app

help:
	@echo "Comandos disponíveis:"
	@echo "  make build       - Compila a aplicação"
	@echo "  make run         - Executa a aplicação"
	@echo "  make test        - Executa os testes"
	@echo "  make swagger     - Gera a documentação Swagger"
	@echo "  make docker-build - Constrói a imagem Docker"
	@echo "  make docker-run  - Executa o contêiner Docker"
	@echo "  make docker-stop - Para o contêiner Docker"
	@echo "  make clean       - Remove binários e arquivos temporários"

build:
	@echo "Compilando aplicação..."
	go build -o bin/$(APP_NAME) $(MAIN_PATH)
	@echo "Aplicação compilada com sucesso em bin/$(APP_NAME)"

run: swagger
	@echo "Executando aplicação..."
	go run $(MAIN_PATH)

test:
	@echo "Executando testes..."
	go test ./...

swagger:
	@echo "Gerando documentação Swagger..."
	@if ! command -v swag > /dev/null; then \
		echo "Instalando swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	swag init -g $(MAIN_PATH) --output ./docs --parseDependency true
	@echo "Documentação Swagger gerada com sucesso!"

docker-build:
	@echo "Construindo imagem Docker..."
	docker build -t $(DOCKER_IMAGE_NAME) .

docker-run: docker-build
	@echo "Executando contêiner Docker..."
	docker-compose up -d
	@echo "Serviço iniciado em: http://localhost:8080"
	@echo "Documentação Swagger: http://localhost:8080/swagger/index.html"

docker-stop:
	@echo "Parando contêiner Docker..."
	docker-compose down

clean:
	@echo "Limpando arquivos temporários..."
	rm -rf bin/
	rm -rf docs/
	@echo "Limpeza concluída!"