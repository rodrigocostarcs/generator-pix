package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rodrigocostarcs/pix-generator/internal/interfaces/api/handlers"
)

// SetupRoutes configura todas as rotas da API
func SetupRoutes(router *gin.Engine, pixHandler *handlers.PixHandler) {
	// Grupo de API
	api := router.Group("/api")
	{
		// Rota para geração de PIX
		api.POST("/generate", pixHandler.GeneratePix)
	}

	// Rota para download de QR code
	router.GET("/download-qrcode", pixHandler.DownloadQRCode)

	// Rota para página inicial (pode ser utilizada para interface web)
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Gerador de PIX API - Use a rota /api/generate para gerar códigos PIX")
	})
}
