package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rodrigocostarcs/pix-generator/internal/interfaces/api/handlers"
	"github.com/rodrigocostarcs/pix-generator/internal/interfaces/api/middlewares"
)

// SetupRoutes configura todas as rotas da API
func SetupRoutes(
	router *gin.Engine,
	pixHandler *handlers.PixHandler,
	autenticacaoHandler *handlers.AutenticacaoHandler,
	autenticacaoMiddleware *middlewares.AutenticacaoMiddleware,
) {
	// Rotas públicas
	api := router.Group("/api")
	{
		// Autenticação
		api.POST("/registrar", autenticacaoHandler.Registrar)
		api.POST("/login", autenticacaoHandler.Login)

		// Rota para download de QR code (mantida pública)
		api.GET("/download-qrcode", pixHandler.DownloadQRCode)
	}

	// Rotas protegidas
	protected := router.Group("/api")
	protected.Use(autenticacaoMiddleware.RequererAutenticacao())
	{
		// Rota para geração de PIX
		protected.POST("/generate", pixHandler.GeneratePix)
	}

	// Rota para página inicial (pode ser utilizada para interface web)
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Gerador de PIX API - Use a rota /api/login para autenticar")
	})
}
