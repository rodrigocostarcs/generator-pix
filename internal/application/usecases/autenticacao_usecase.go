package usecases

import (
	"errors"

	"github.com/rodrigocostarcs/pix-generator/internal/domain/models"
	"github.com/rodrigocostarcs/pix-generator/internal/domain/services"
	"github.com/rodrigocostarcs/pix-generator/internal/infrastructure/repositories"
)

// AutenticacaoUseCase implementa o caso de uso para autenticação
type AutenticacaoUseCase struct {
	autenticacaoService       *services.AutenticacaoService
	estabelecimentoRepository repositories.EstabelecimentoRepository
}

// NewAutenticacaoUseCase cria uma nova instância do caso de uso de autenticação
func NewAutenticacaoUseCase(autenticacaoService *services.AutenticacaoService, estabelecimentoRepository repositories.EstabelecimentoRepository) *AutenticacaoUseCase {
	return &AutenticacaoUseCase{
		autenticacaoService:       autenticacaoService,
		estabelecimentoRepository: estabelecimentoRepository,
	}
}

// Registrar registra um novo estabelecimento
func (uc *AutenticacaoUseCase) Registrar(req models.EstabelecimentoRequest) (models.EstabelecimentoResponse, error) {
	// Verificar se já existe um estabelecimento com o mesmo email
	_, err := uc.estabelecimentoRepository.BuscarPorEmail(req.Email)
	if err == nil {
		return models.EstabelecimentoResponse{}, errors.New("email já cadastrado")
	}

	// Criar o estabelecimento
	estabelecimento, err := uc.estabelecimentoRepository.Salvar(req)
	if err != nil {
		return models.EstabelecimentoResponse{}, err
	}

	// Preparar a resposta sem expor a senha
	response := models.EstabelecimentoResponse{
		ID:           estabelecimento.ID,
		Nome:         estabelecimento.Nome,
		Descricao:    estabelecimento.Descricao,
		Email:        estabelecimento.Email,
		Ativo:        estabelecimento.Ativo,
		CriadoEm:     estabelecimento.CriadoEm,
		AtualizadoEm: estabelecimento.AtualizadoEm,
	}

	return response, nil
}

// Login autentica um estabelecimento e gera um token JWT
func (uc *AutenticacaoUseCase) Login(req models.LoginRequest) (models.LoginResponse, error) {
	// Buscar o estabelecimento pelo email
	estabelecimento, err := uc.estabelecimentoRepository.BuscarPorEmail(req.Email)
	if err != nil {
		return models.LoginResponse{}, errors.New("credenciais inválidas")
	}

	// Verificar se o estabelecimento está ativo
	if !estabelecimento.Ativo {
		return models.LoginResponse{}, errors.New("conta desativada")
	}

	// Verificar a senha
	err = uc.autenticacaoService.VerificarSenha(estabelecimento.Senha, req.Senha)
	if err != nil {
		return models.LoginResponse{}, errors.New("credenciais inválidas")
	}

	// Gerar token JWT
	token, err := uc.autenticacaoService.GerarToken(estabelecimento)
	if err != nil {
		return models.LoginResponse{}, err
	}

	// Preparar a resposta
	response := models.LoginResponse{
		Token: token,
		Estabelecimento: models.EstabelecimentoResponse{
			ID:           estabelecimento.ID,
			Nome:         estabelecimento.Nome,
			Descricao:    estabelecimento.Descricao,
			Email:        estabelecimento.Email,
			Ativo:        estabelecimento.Ativo,
			CriadoEm:     estabelecimento.CriadoEm,
			AtualizadoEm: estabelecimento.AtualizadoEm,
		},
	}

	return response, nil
}
