package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/rodrigocostarcs/pix-generator/docs"
	"github.com/rodrigocostarcs/pix-generator/internal/application/usecases"
	"github.com/rodrigocostarcs/pix-generator/internal/domain/services"
	"github.com/rodrigocostarcs/pix-generator/internal/infrastructure/cache"
	"github.com/rodrigocostarcs/pix-generator/internal/infrastructure/repositories"
	"github.com/rodrigocostarcs/pix-generator/internal/interfaces/api/handlers"
	"github.com/rodrigocostarcs/pix-generator/internal/interfaces/api/middlewares"
	"github.com/rodrigocostarcs/pix-generator/internal/interfaces/api/routes"
)

// @title           Gerador de PIX API
// @version         1.0
// @description     API para geração e gerenciamento de códigos PIX seguindo arquitetura DDD.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Desenvolvedor
// @contact.email  contato@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Digite 'Bearer ' seguido do token JWT

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	// Configuração do banco de dados
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "generator_pix")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Falha ao pingar o banco de dados: %v", err)
	}

	// Configuração do Redis para cache
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")

	// Configurar o diretório de templates
	templatesDir := getEnv("TEMPLATES_DIR", "./templates")

	// Obter o caminho absoluto para o diretório de templates
	absTemplatesDir, err := filepath.Abs(templatesDir)
	if err != nil {
		log.Printf("Aviso: não foi possível obter o caminho absoluto para templates: %v", err)
		absTemplatesDir = templatesDir
	}

	log.Printf("Diretório de templates (caminho absoluto): %s", absTemplatesDir)

	// Criar diretório se não existir
	if _, err := os.Stat(absTemplatesDir); os.IsNotExist(err) {
		log.Printf("Criando diretório de templates em: %s", absTemplatesDir)
		if err := os.MkdirAll(absTemplatesDir, 0755); err != nil {
			log.Printf("Aviso: não foi possível criar o diretório de templates: %v", err)
		}
	}

	// Listar arquivos no diretório
	files, err := os.ReadDir(absTemplatesDir)
	if err != nil {
		log.Printf("Erro ao ler diretório de templates: %v", err)
	} else {
		log.Printf("Arquivos no diretório %s:", absTemplatesDir)
		for _, file := range files {
			log.Printf(" - %s", file.Name())
		}
	}

	// Inicializar o processador de templates com o caminho absoluto
	templateProcessor := services.NewTemplateProcessor(absTemplatesDir)

	// Criar adaptador de cache
	cacheAdapter := cache.NewRedisAdapter(redisHost, redisPort, redisPassword, 0)

	// Serviços
	pixService := services.NewPixGeneratorService()
	autenticacaoService := services.NewAutenticacaoService()

	// Repositórios
	pixRepository := repositories.NewMysqlPixRepository(db)
	estabelecimentoRepository := repositories.NewMysqlEstabelecimentoRepository(db)

	// Casos de uso
	generatePixUseCase := usecases.NewGeneratePixUseCase(pixService, pixRepository)
	autenticacaoUseCase := usecases.NewAutenticacaoUseCase(autenticacaoService, estabelecimentoRepository)

	// Handlers
	pixHandler := handlers.NewPixHandler(generatePixUseCase, pixRepository, cacheAdapter, templateProcessor)
	autenticacaoHandler := handlers.NovaAutenticacaoHandler(autenticacaoUseCase)

	// Middlewares
	autenticacaoMiddleware := middlewares.NewAutenticacaoMiddleware(autenticacaoService)

	// Configurar o router Gin
	router := gin.Default()

	// Servir arquivos estáticos (incluindo templates)
	router.Static("/templates", absTemplatesDir)

	// Configurar Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Configurar as rotas
	routes.SetupRoutes(router, pixHandler, autenticacaoHandler, autenticacaoMiddleware)

	// Iniciar o servidor
	port := getEnv("PORT", "8080")
	log.Printf("Servidor iniciado na porta %s", port)
	log.Printf("Documentação Swagger disponível em: http://localhost:%s/swagger/index.html", port)
	log.Printf("Métricas Prometheus disponíveis em: http://localhost:%s/metrics", port)
	router.Run(":" + port)
}

// getEnv obtém uma variável de ambiente ou usa o valor padrão
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
