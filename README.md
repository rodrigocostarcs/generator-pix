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

## Componentes Principais

- **API REST**: Interface para geração e gerenciamento de códigos PIX
- **Cache com Redis**: Armazenamento em cache para melhorar a performance
- **Monitoramento com Prometheus**: Coleta de métricas para monitoramento
- **Visualização com Grafana**: Dashboard para visualização das métricas
- **Templates para QR Codes**: Sistema para aplicar templates visuais aos QR codes
- **Documentação com Swagger**: Documentação interativa da API

## Primeiros Passos

### Pré-requisitos

- Go 1.24 ou superior
- Docker e Docker Compose
- Git

### Clonar o Repositório

```bash
git clone https://github.com/rodrigocostarcs/pix-generator.git
cd pix-generator
```

### Configuração de Ambiente

Copie o arquivo de exemplo de configuração:

```bash
cp .env.example .env
```

Você pode editar o arquivo `.env` para personalizar as configurações.

### Execução com Docker

A maneira mais fácil de executar a aplicação é usando Docker Compose:

```bash
# Construir e iniciar todos os serviços
docker-compose up -d

# Para ver os logs da aplicação
docker-compose logs -f app
```

Isso iniciará:
- A aplicação principal na porta 8080
- Banco de dados MySQL na porta 3307
- Redis para cache na porta 6379
- Prometheus para métricas na porta 9090
- Grafana para visualização na porta 3000

### Execução Local para Desenvolvimento

Se preferir executar localmente para desenvolvimento:

```bash
# Instalar dependências
go mod tidy

# Gerar documentação Swagger
make swagger

# Executar a aplicação
make run
```

## Documentação da API

A documentação da API é gerada automaticamente usando Swagger. Após iniciar a aplicação, você pode acessá-la em:

```
http://localhost:8080/swagger/index.html
```

### Como Gerar a Documentação Swagger

Se fizer alterações na API, precisará regenerar a documentação:

```bash
# Usando Make
make swagger

# Ou manualmente
./scripts/gen-swagger.sh
```

### Endpoints Principais da API

- `POST /api/registrar` - Registrar um novo estabelecimento
- `POST /api/login` - Autenticar um estabelecimento e obter um token JWT
- `POST /api/generate` - Gerar um código PIX (requer autenticação)
- `GET /api/download-qrcode` - Baixar imagem do QR code (suporta templates)

## Monitoramento

### Métricas com Prometheus

As métricas da aplicação são expostas no endpoint:

```
http://localhost:8080/metrics
```

Para acessar o console do Prometheus:

```
http://localhost:9090
```

### Visualização com Grafana

Para acessar o dashboard Grafana:

```
http://localhost:3000
```

Credenciais padrão:
- Usuário: `admin`
- Senha: `admin`

Na primeira vez que acessar, você precisará configurar a fonte de dados do Prometheus:
1. Acesse "Configuration > Data Sources"
2. Clique em "Add data source" e selecione "Prometheus"
3. No campo "URL", insira `http://prometheus:9090`
4. Clique em "Save & Test"

## Templates de QR Code

A aplicação suporta a aplicação de templates visuais aos QR codes gerados. Para usar esta funcionalidade:

1. Adicione arquivos de template PNG na pasta `templates/`
2. Ao solicitar o download de um QR code, especifique o parâmetro `template`:

```
GET /api/download-qrcode?codigo_pix=SEU_CODIGO_PIX&template=template_pix_1
```

## Cache com Redis

O sistema utiliza Redis para cache, melhorando a performance especialmente para operações frequentes como download de QR codes. O cache é configurado automaticamente quando a aplicação é iniciada com Docker Compose.

## Executando Testes

Para executar os testes automatizados:

```bash
# Todos os testes
make test

# Ou usando Go diretamente
go test ./...

# Para testes com cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Comandos Úteis

O projeto inclui um Makefile para facilitar tarefas comuns:

- `make build` - Compila a aplicação
- `make run` - Executa a aplicação localmente
- `make test` - Executa os testes
- `make swagger` - Gera a documentação Swagger
- `make docker-build` - Constrói a imagem Docker
- `make docker-run` - Executa o contêiner Docker
- `make docker-stop` - Para o contêiner Docker
- `make clean` - Remove binários e arquivos temporários
- `make help` - Exibe a lista de comandos disponíveis

## Estrutura do Banco de Dados

O projeto utiliza MySQL como banco de dados principal. O esquema é inicializado automaticamente pelo script `scripts/init.sql` quando o contêiner Docker é iniciado pela primeira vez.

Principais tabelas:
- `estabelecimentos` - Armazena informações dos estabelecimentos
- `pix` - Armazena os códigos PIX gerados

## Troubleshooting

### Problemas Comuns

#### Erro de conexão com o banco de dados
- Verifique se o serviço MySQL está em execução: `docker-compose ps`
- Verifique as credenciais no arquivo `.env`

#### Templates não são carregados
- Verifique se os arquivos de template estão na pasta `templates/`
- Verifique as permissões dos arquivos

#### Swagger não está acessível
- Verifique se a documentação foi gerada: `make swagger`
- Verifique se a aplicação está em execução na porta 8080

## Contribuindo

1. Faça um fork do repositório
2. Crie sua branch de feature: `git checkout -b feature/minha-nova-feature`
3. Commit suas alterações: `git commit -am 'Adicionar nova feature'`
4. Push para a branch: `git push origin feature/minha-nova-feature`
5. Envie um pull request

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo LICENSE para detalhes.
