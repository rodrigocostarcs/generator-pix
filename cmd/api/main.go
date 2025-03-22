package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/rodrigocostarcs/pix-generator/internal/application/usecases"
	"github.com/rodrigocostarcs/pix-generator/internal/domain/services"
	"github.com/rodrigocostarcs/pix-generator/internal/infrastructure/cache"
	"github.com/rodrigocostarcs/pix-generator/internal/infrastructure/repositories"
	"github.com/rodrigocostarcs/pix-generator/internal/interfaces/api/handlers"
	"github.com/rodrigocostarcs/pix-generator/internal/interfaces/api/middlewares"
	"github.com/rodrigocostarcs/pix-generator/internal/interfaces/api/routes"
)

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
	pixHandler := handlers.NewPixHandler(generatePixUseCase, pixRepository, cacheAdapter)
	autenticacaoHandler := handlers.NovaAutenticacaoHandler(autenticacaoUseCase)

	// Middlewares
	autenticacaoMiddleware := middlewares.NewAutenticacaoMiddleware(autenticacaoService)

	// Configurar o router Gin
	router := gin.Default()

	// Configurar as rotas
	routes.SetupRoutes(router, pixHandler, autenticacaoHandler, autenticacaoMiddleware)

	// Iniciar o servidor
	port := getEnv("PORT", "8080")
	log.Printf("Servidor iniciado na porta %s", port)
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
