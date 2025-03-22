package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rodrigocostarcs/pix-generator/internal/application/usecases"
	"github.com/rodrigocostarcs/pix-generator/internal/domain/models"
	"github.com/rodrigocostarcs/pix-generator/internal/interfaces/api/views"
)

// AutenticacaoHandler manipula as requisições da API relacionadas à autenticação
type AutenticacaoHandler struct {
	autenticacaoUseCase *usecases.AutenticacaoUseCase
	responseView        *views.ResponseView
}

// NovaAutenticacaoHandler cria uma nova instância do handler de autenticação
func NovaAutenticacaoHandler(autenticacaoUseCase *usecases.AutenticacaoUseCase) *AutenticacaoHandler {
	return &AutenticacaoHandler{
		autenticacaoUseCase: autenticacaoUseCase,
		responseView:        views.NewResponseView(),
	}
}

// Registrar processa a requisição para registro de um novo estabelecimento
func (h *AutenticacaoHandler) Registrar(c *gin.Context) {
	var req models.EstabelecimentoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.responseView.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Executar o caso de uso
	response, err := h.autenticacaoUseCase.Registrar(req)
	if err != nil {
		h.responseView.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	h.responseView.Success(c, http.StatusCreated, response)
}

// Login processa a requisição de login
func (h *AutenticacaoHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.responseView.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Executar o caso de uso
	response, err := h.autenticacaoUseCase.Login(req)
	if err != nil {
		h.responseView.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	h.responseView.Success(c, http.StatusOK, response)
}
