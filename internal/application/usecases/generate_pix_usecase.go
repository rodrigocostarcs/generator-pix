package usecases

import (
	"time"

	"github.com/rodrigocostarcs/pix-generator/internal/domain/models"
	"github.com/rodrigocostarcs/pix-generator/internal/domain/services"
	"github.com/rodrigocostarcs/pix-generator/internal/infrastructure/repositories"
)

// GeneratePixUseCase implementa o caso de uso para geração de PIX
type GeneratePixUseCase struct {
	pixService    *services.PixGeneratorService
	pixRepository repositories.PixRepository
}

// GetRepository retorna o repositório usado pelo caso de uso
func (uc *GeneratePixUseCase) GetRepository() repositories.PixRepository {
	return uc.pixRepository
}

// NewGeneratePixUseCase cria uma nova instância do caso de uso
func NewGeneratePixUseCase(pixService *services.PixGeneratorService, pixRepository repositories.PixRepository) *GeneratePixUseCase {
	return &GeneratePixUseCase{
		pixService:    pixService,
		pixRepository: pixRepository,
	}
}

// Execute executa o caso de uso para geração de PIX
func (uc *GeneratePixUseCase) Execute(req models.PixRequest) (models.PixResponse, error) {
	// Gerar o código PIX através do serviço de domínio
	pixResponse, err := uc.pixService.GerarPixEstatico(req)
	if err != nil {
		return models.PixResponse{}, err
	}

	// Criar a entidade PIX para persistência
	pix := models.Pix{
		Nome:          req.Nome,
		Chave:         req.Chave,
		Cidade:        req.Cidade,
		Valor:         req.Valor,
		Identificador: req.Identificador,
		Descricao:     req.Descricao,
		CodigoPix:     pixResponse.CodigoPix,
		QRCodeSVG:     pixResponse.QRCodeSVG,
		QRCodePNG:     pixResponse.QRCodePNG,
		CriadoEm:      time.Now(),
	}

	// Persistir a entidade no banco de dados
	_, err = uc.pixRepository.Save(pix)
	if err != nil {
		return models.PixResponse{}, err
	}

	return pixResponse, nil
}
