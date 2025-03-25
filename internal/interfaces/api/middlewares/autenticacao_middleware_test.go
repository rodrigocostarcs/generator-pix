package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rodrigocostarcs/pix-generator/internal/domain/models"
	"github.com/rodrigocostarcs/pix-generator/internal/domain/services"
	"github.com/rodrigocostarcs/pix-generator/internal/interfaces/api/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestAutenticacaoMiddleware(t *testing.T) {
	// Configurar Gin para modo de teste
	gin.SetMode(gin.TestMode)

	// Configurar chave JWT para testes
	jwtSecret := "chave_secreta_para_testes"
	os.Setenv("JWT_SECRET", jwtSecret)
	defer os.Unsetenv("JWT_SECRET")

	// Inicializar serviço e middleware
	authService := services.NewAutenticacaoService()
	middleware := middlewares.NewAutenticacaoMiddleware(authService)

	// Criar um estabelecimento de teste
	estabelecimento := models.Estabelecimento{
		ID:           "123e4567-e89b-12d3-a456-426614174000",
		Nome:         "Loja Teste",
		Email:        "loja@teste.com",
		Ativo:        true,
		CriadoEm:     time.Now(),
		AtualizadoEm: time.Now(),
	}

	t.Run("TokenValido", func(t *testing.T) {
		// Gerar token válido
		token, err := authService.GerarToken(estabelecimento)
		assert.NoError(t, err)

		// Criar request com token
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest("GET", "/api/protected", nil)
		c.Request.Header.Set("Authorization", "Bearer "+token)

		// Criar uma rota de teste para verificar se o middleware passa o controle
		var handlerCalled bool
		handler := func(c *gin.Context) {
			handlerCalled = true

			// Verificar se as claims estão no contexto
			id, exists := c.Get("usuarioID")
			assert.True(t, exists)
			assert.Equal(t, estabelecimento.ID, id)

			email, exists := c.Get("usuarioEmail")
			assert.True(t, exists)
			assert.Equal(t, estabelecimento.Email, email)

			nome, exists := c.Get("usuarioNome")
			assert.True(t, exists)
			assert.Equal(t, estabelecimento.Nome, nome)

			c.Status(http.StatusOK)
		}

		// Executar middleware e handler
		mw := middleware.RequererAutenticacao()
		mw(c)

		// Se o middleware não interrompeu o fluxo, chamar o handler
		if !c.IsAborted() {
			handler(c)
		}

		// Verificações
		assert.True(t, handlerCalled, "O handler deveria ter sido chamado")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("SemToken", func(t *testing.T) {
		// Criar request sem token de autorização
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest("GET", "/api/protected", nil)

		// Criar uma rota de teste
		var handlerCalled bool
		handler := func(c *gin.Context) {
			handlerCalled = true
			c.Status(http.StatusOK)
		}

		// Executar middleware e handler
		mw := middleware.RequererAutenticacao()
		mw(c)

		// Se o middleware não interrompeu o fluxo, chamar o handler
		if !c.IsAborted() {
			handler(c)
		}

		// Verificações
		assert.False(t, handlerCalled, "O handler não deveria ter sido chamado")
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Verificar resposta JSON
		assert.Contains(t, w.Body.String(), "Autorização necessária")
	})

	t.Run("FormatoTokenInvalido", func(t *testing.T) {
		// Criar request com formato de token inválido
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest("GET", "/api/protected", nil)
		c.Request.Header.Set("Authorization", "InvalidTokenFormat")

		// Criar uma rota de teste
		var handlerCalled bool
		handler := func(c *gin.Context) {
			handlerCalled = true
			c.Status(http.StatusOK)
		}

		// Executar middleware e handler
		mw := middleware.RequererAutenticacao()
		mw(c)

		// Se o middleware não interrompeu o fluxo, chamar o handler
		if !c.IsAborted() {
			handler(c)
		}

		// Verificações
		assert.False(t, handlerCalled, "O handler não deveria ter sido chamado")
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Formato de token inválido")
	})

	t.Run("TokenExpirado", func(t *testing.T) {
		// Criar claims com tempo de expiração no passado
		expirado := time.Now().Add(-1 * time.Hour).Unix()

		claims := jwt.MapClaims{
			"id":    estabelecimento.ID,
			"email": estabelecimento.Email,
			"nome":  estabelecimento.Nome,
			"exp":   expirado,
		}

		// Criar token expirado manualmente
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(jwtSecret))

		// Criar request com token expirado
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest("GET", "/api/protected", nil)
		c.Request.Header.Set("Authorization", "Bearer "+tokenString)

		// Criar uma rota de teste
		var handlerCalled bool
		handler := func(c *gin.Context) {
			handlerCalled = true
			c.Status(http.StatusOK)
		}

		// Executar middleware e handler
		mw := middleware.RequererAutenticacao()
		mw(c)

		// Se o middleware não interrompeu o fluxo, chamar o handler
		if !c.IsAborted() {
			handler(c)
		}

		// Verificações
		assert.False(t, handlerCalled, "O handler não deveria ter sido chamado")
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Token inválido")
	})
}
