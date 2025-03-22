# Gerador de PIX - Implementação em Go com Arquitetura DDD

Este projeto é uma implementação em Go de um gerador de códigos PIX (Sistema de Pagamento Instantâneo Brasileiro). A aplicação segue os princípios de Domain-Driven Design (DDD) para garantir uma arquitetura limpa, separação de responsabilidades e facilidade de manutenção.

## Visão Geral da Arquitetura

O projeto está estruturado de acordo com os padrões arquiteturais do DDD:

```
pix-generator/
├── cmd/                    # Pontos de entrada da aplicação
│   └── api/                # Ponto de entrada do servidor API
├── docs/                   # Documentação Swagger gerada automaticamente
├── internal/               # Código privado da aplicação
│   ├── domain/             # Camada de domínio (regras de negócio principais)
│   │   ├── models/         # Entidades de domínio e objetos de valor
│   │   └── services/       # Serviços de domínio e lógica de negócio
│   ├── application/        # Camada de aplicação (orquestração)
│   │   └── usecases/       # Casos de uso que coordenam operações de domínio
│   ├── infrastructure/     # Camada de infraestrutura (capacidades técnicas)
│   │   └── repositories/   # Implementações de persistência de dados
│   └── interfaces/         # Camada de interface (adaptadores para sistemas externos)
│       └── api/            # Interface da API REST
│           ├── handlers/   # Manipuladores de requisições HTTP
│           ├── routes/     # Definições de rotas da API
│           └── middlewares/# Middlewares HTTP
├── pkg/                    # Pacotes compartilhados, públicos
│   └── utils/              # Utilitários e auxiliares
├── scripts/                # Scripts para desenvolvimento e implantação
│   ├── dev-entrypoint.sh   # Script de inicialização para desenvolvimento
│   └── gen-swagger.sh      # Script para gerar documentação Swagger
└── web/                    # Recursos web
    ├── static/             # Arquivos estáticos (CSS, JS, imagens)
    └── templates/          # Templates HTML
```

## Camadas Arquiteturais

### Camada de Domínio (internal/domain)

O coração da aplicação, contendo toda a lógica e regras de negócio. Esta camada é independente de outras camadas e frameworks externos.

- **Models**: Entidades principais de negócio como códigos PIX e objetos de valor associados.
- **Services**: Lógica de negócio para geração de códigos PIX, QR codes e validações.

### Camada de Aplicação (internal/application)

Orquestra o fluxo de dados de e para a camada de domínio, e coordena a execução de operações de negócio.

- **Use Cases**: Regras de negócio específicas da aplicação que orquestram operações de domínio.

### Camada de Infraestrutura (internal/infrastructure)

Fornece capacidades técnicas para suportar a aplicação, como acesso ao banco de dados e integrações com APIs externas.

- **Repositories**: Implementações para persistência e recuperação de dados.

### Camada de Interface (internal/interfaces)

Contém os adaptadores que conectam a aplicação ao mundo exterior.

- **API**: Implementação da API REST com handlers, rotas e middlewares.

## Benefícios da Arquitetura DDD

1. **Separação de Responsabilidades**: Cada camada tem uma responsabilidade específica, tornando o código mais fácil de manter.
2. **Foco no Negócio**: A arquitetura se concentra na modelagem do domínio de negócio.
3. **Testabilidade**: A lógica de domínio pode ser testada independentemente das preocupações de infraestrutura.
4. **Flexibilidade**: Detalhes de implementação podem mudar sem afetar a lógica de domínio.
5. **Escalabilidade**: Diferentes camadas podem escalar independentemente conforme necessário.

## Primeiros Passos

### Pré-requisitos

- Go 1.24 ou superior
- MySQL 8.0
- Redis (para cache)
- Docker e Docker Compose (opcional)

### Execução com Docker

1. Clone o repositório
2. Configure as variáveis de ambiente (ou use os valores padrão no arquivo .env.example)
3. Execute com Docker Compose:

```bash
make docker-run
```

A aplicação estará disponível em http://localhost:8080 e a documentação Swagger em http://localhost:8080/swagger/index.html.

### Execução Local

1. Clone o repositório
2. Configure as variáveis de ambiente (ou copie de .env.example para .env)
3. Instale as dependências: `go mod tidy`
4. Gere a documentação Swagger: `make swagger`
5. Execute a aplicação: `make run`

## Documentação da API

A documentação da API é gerada automaticamente usando Swagger. Após iniciar a aplicação, você pode acessá-la em:

http://localhost:8080/swagger/index.html

### Endpoints da API

- `POST /api/registrar` - Registrar um novo estabelecimento
- `POST /api/login` - Autenticar um estabelecimento e obter um token JWT
- `POST /api/generate` - Gerar um código PIX (requer autenticação)
- `GET /api/download-qrcode` - Baixar imagem do QR code

## Comandos úteis

O projeto inclui um Makefile para facilitar tarefas comuns:

- `make build` - Compila a aplicação
- `make run` - Executa a aplicação localmente
- `make test` - Executa os testes
- `make swagger` - Gera a documentação Swagger
- `make docker-build` - Constrói a imagem Docker
- `make docker-run` - Executa o contêiner Docker
- `make docker-stop` - Para o contêiner Docker
- `make clean` - Remove binários e arquivos temporários

## Tecnologias

- Go (Golang) 1.24
- Gin (Framework Web)
- MySQL (Banco de dados)
- Redis (Cache)
- JWT (Autenticação)
- Swagger (Documentação da API)
- Docker e Docker Compose (Containerização)

## Contribuindo

1. Faça um fork do repositório
2. Crie sua branch de feature: `git checkout -b feature/minha-nova-feature`
3. Commit suas alterações: `git commit -am 'Adicionar nova feature'`
4. Push para a branch: `git push origin feature/minha-nova-feature`
5. Envie um pull request
