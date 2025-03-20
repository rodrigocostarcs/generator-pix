package handlers

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rodrigocostarcs/pix-generator/internal/application/usecases"
	"github.com/rodrigocostarcs/pix-generator/internal/domain/models"
	"github.com/rodrigocostarcs/pix-generator/internal/infrastructure/repositories"
)

// PixHandler manipula as requisições da API relacionadas ao PIX
type PixHandler struct {
	generatePixUseCase *usecases.GeneratePixUseCase
	pixRepository      repositories.PixRepository
}

// NewPixHandler cria uma nova instância do handler PIX
func NewPixHandler(generatePixUseCase *usecases.GeneratePixUseCase, pixRepository repositories.PixRepository) *PixHandler {
	return &PixHandler{
		generatePixUseCase: generatePixUseCase,
		pixRepository:      pixRepository,
	}
}

// GeneratePix processa a requisição para gerar um código PIX
func (h *PixHandler) GeneratePix(c *gin.Context) {
	var req models.PixRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validações básicas
	if req.Nome == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome é obrigatório"})
		return
	}

	if req.Chave == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chave PIX é obrigatória"})
		return
	}

	// Definir cidade padrão se não fornecida
	if req.Cidade == "" {
		req.Cidade = "São Paulo"
	}

	// Executar o caso de uso
	response, err := h.generatePixUseCase.Execute(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DownloadQRCode manipula o download do QR code
func (h *PixHandler) DownloadQRCode(c *gin.Context) {
	codigoPix := c.Query("codigo_pix")

	if codigoPix == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código PIX é obrigatório"})
		return
	}

	pixRepository := h.getPixRepository()

	pix, err := pixRepository.FindByCodigoPix(codigoPix)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Código PIX não encontrado"})
		return
	}

	// Remover o prefixo data:image/png;base64,
	base64Data := strings.TrimPrefix(pix.QRCodePNG, "data:image/png;base64,")

	// Decodificar o base64
	decodedData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao decodificar QR code: " + err.Error()})
		return
	}

	// Configurar os headers para download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=pix_qrcode.png")
	c.Header("Content-Type", "image/png")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	// Enviar o PNG como resposta
	c.Data(http.StatusOK, "image/png", decodedData)
}

// getPixRepository método auxiliar para obter o repositório do caso de uso
func (h *PixHandler) getPixRepository() repositories.PixRepository {
	// Esta é uma solução temporária - corrigir para injetar o repositório diretamente no handler

	return h.generatePixUseCase.GetRepository()
}
