package handlers

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rodrigocostarcs/pix-generator/internal/application/usecases"
	"github.com/rodrigocostarcs/pix-generator/internal/domain/models"
	"github.com/rodrigocostarcs/pix-generator/internal/infrastructure/cache"
	"github.com/rodrigocostarcs/pix-generator/internal/infrastructure/repositories"
	"github.com/rodrigocostarcs/pix-generator/internal/interfaces/api/views"
)

// CachedPix representa a estrutura que será armazenada em cache
type CachedPix struct {
	Pix     models.Pix `json:"pix"`
	PngData []byte     `json:"png_data,omitempty"` // Dados binários do PNG já decodificados
}

// PixHandler manipula as requisições da API relacionadas ao PIX
type PixHandler struct {
	generatePixUseCase *usecases.GeneratePixUseCase
	pixRepository      repositories.PixRepository
	responseView       *views.ResponseView
	cacheAdapter       cache.CacheAdapter
}

// NewPixHandler cria uma nova instância do handler PIX
func NewPixHandler(generatePixUseCase *usecases.GeneratePixUseCase, pixRepository repositories.PixRepository, cacheAdapter cache.CacheAdapter) *PixHandler {
	return &PixHandler{
		generatePixUseCase: generatePixUseCase,
		pixRepository:      pixRepository,
		responseView:       views.NewResponseView(),
		cacheAdapter:       cacheAdapter,
	}
}

// GeneratePix processa a requisição para gerar um código PIX
func (h *PixHandler) GeneratePix(c *gin.Context) {
	var req models.PixRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.responseView.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validações básicas
	if req.Nome == "" {
		h.responseView.Error(c, http.StatusBadRequest, "Nome é obrigatório")
		return
	}

	if req.Chave == "" {
		h.responseView.Error(c, http.StatusBadRequest, "Chave PIX é obrigatória")
		return
	}

	// Definir cidade padrão se não fornecida
	if req.Cidade == "" {
		req.Cidade = "São Paulo"
	}

	// Executar o caso de uso
	response, err := h.generatePixUseCase.Execute(req)
	if err != nil {
		h.responseView.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.responseView.Success(c, http.StatusOK, response)
}

// DownloadQRCode manipula o download do QR code
func (h *PixHandler) DownloadQRCode(c *gin.Context) {
	codigoPix := c.Query("codigo_pix")

	if codigoPix == "" {
		h.responseView.Error(c, http.StatusBadRequest, "Código PIX é obrigatório")
		return
	}

	ctx := context.Background()

	// Verificar se os dados do QR code estão em cache
	var cachedData CachedPix
	cacheKey := "pix_qrcode:" + codigoPix

	// Tentar obter do cache
	err := cache.GetObject(h.cacheAdapter, ctx, cacheKey, &cachedData)

	// Se não estiver em cache ou houver erro, buscar do banco de dados
	if err != nil {
		// Buscar o PIX do repositório
		pix, err := h.pixRepository.FindByCodigoPix(codigoPix)
		if err != nil {
			h.responseView.Error(c, http.StatusNotFound, "Código PIX não encontrado")
			return
		}

		// Processar dados PNG
		base64Data := strings.TrimPrefix(pix.QRCodePNG, "data:image/png;base64,")
		decodedData, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			h.responseView.Error(c, http.StatusInternalServerError, "Erro ao decodificar QR code: "+err.Error())
			return
		}

		// Criar objeto para cache
		cachedData = CachedPix{
			Pix:     pix,
			PngData: decodedData,
		}

		// Armazenar no cache por 24 horas
		_ = cache.SetObject(h.cacheAdapter, ctx, cacheKey, cachedData, 24*time.Hour)
	}

	// Verificar formato solicitado
	format := c.Query("format")
	if format == "json" {
		h.responseView.Success(c, http.StatusOK, gin.H{
			"codigo_pix": cachedData.Pix.CodigoPix,
			"qrcode_png": cachedData.Pix.QRCodePNG,
			"qrcode_svg": cachedData.Pix.QRCodeSVG,
		})
		return
	}

	// Para o download de dados binários, usamos o método Download da responseView
	h.responseView.Download(c, "pix_qrcode.png", "image/png", cachedData.PngData)
}
