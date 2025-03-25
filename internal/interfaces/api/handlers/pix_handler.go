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
	"github.com/rodrigocostarcs/pix-generator/internal/domain/services"
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
	templateProcessor  *services.TemplateProcessor // Novo campo
}

// NewPixHandler cria uma nova instância do handler PIX
func NewPixHandler(
	generatePixUseCase *usecases.GeneratePixUseCase,
	pixRepository repositories.PixRepository,
	cacheAdapter cache.CacheAdapter,
	templateProcessor *services.TemplateProcessor,
) *PixHandler {
	return &PixHandler{
		generatePixUseCase: generatePixUseCase,
		pixRepository:      pixRepository,
		responseView:       views.NewResponseView(),
		cacheAdapter:       cacheAdapter,
		templateProcessor:  templateProcessor,
	}
}

// GeneratePix processa a requisição para gerar um código PIX
// @Summary      Gerar código PIX
// @Description  Gera um novo código PIX estático com base nos dados fornecidos
// @Tags         pix
// @Accept       json
// @Produce      json
// @Param        request  body      models.PixRequest  true  "Dados para geração do PIX"
// @Success      200      {object}  views.Response{data=models.PixResponse}  "Código PIX gerado com sucesso"
// @Failure      400      {object}  views.Response     "Erro de requisição"
// @Failure      401      {object}  views.Response     "Não autorizado"
// @Failure      500      {object}  views.Response     "Erro interno do servidor"
// @Security     BearerAuth
// @Router       /generate [post]
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
// @Summary      Download QR Code
// @Description  Faz o download de um QR code para o código PIX gerado, opcionalmente aplicando um template
// @Tags         pix
// @Produce      image/png
// @Produce      application/json
// @Param        codigo_pix  query     string  true   "Código PIX gerado"
// @Param        format      query     string  false  "Formato de resposta (json ou png, padrão é png)"
// @Param        template    query     string  false  "Nome do template a ser aplicado (ex: template_pix_1)"
// @Success      200         {file}    file    "QR Code em formato PNG"
// @Success      200         {object}  views.Response{data=models.PixResponse}  "Detalhes do QR Code em JSON"
// @Failure      400         {object}  views.Response  "Código PIX não fornecido"
// @Failure      404         {object}  views.Response  "Código PIX não encontrado"
// @Failure      500         {object}  views.Response  "Erro interno do servidor"
// @Router       /download-qrcode [get]
func (h *PixHandler) DownloadQRCode(c *gin.Context) {
	codigoPix := c.Query("codigo_pix")
	templateName := c.Query("template")

	if codigoPix == "" {
		h.responseView.Error(c, http.StatusBadRequest, "Código PIX é obrigatório")
		return
	}

	ctx := context.Background()

	// Verificar se os dados do QR code estão em cache
	var cachedData CachedPix
	cacheKey := "pix_qrcode:" + codigoPix
	if templateName != "" {
		cacheKey = cacheKey + ":" + templateName
	}

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

	// Se um template foi especificado, aplicá-lo
	if templateName != "" {
		templateImgData, err := h.templateProcessor.ApplyTemplate(cachedData.Pix.QRCodePNG, templateName)
		if err != nil {
			h.responseView.Error(c, http.StatusInternalServerError, "Erro ao aplicar template: "+err.Error())
			return
		}

		h.responseView.Download(c, "pix_template.png", "image/png", templateImgData)
		return
	}

	// Para o download de dados binários, usamos o método Download da responseView
	h.responseView.Download(c, "pix_qrcode.png", "image/png", cachedData.PngData)
}
