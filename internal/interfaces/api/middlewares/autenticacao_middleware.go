package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rodrigocostarcs/pix-generator/internal/domain/services"
)

// AutenticacaoMiddleware estrutura do middleware de autenticação
type AutenticacaoMiddleware struct {
	autenticacaoService *services.AutenticacaoService
}

// NewAutenticacaoMiddleware cria uma nova instância do middleware de autenticação
func NewAutenticacaoMiddleware(autenticacaoService *services.AutenticacaoService) *AutenticacaoMiddleware {
	return &AutenticacaoMiddleware{
		autenticacaoService: autenticacaoService,
	}
}

// RequererAutenticacao middleware que requer autenticação para acessar rotas protegidas
func (m *AutenticacaoMiddleware) RequererAutenticacao() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obter o token do header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Autorização necessária"})
			c.Abort()
			return
		}

		// O token deve estar no formato "Bearer {token}"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Formato de token inválido"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Validar o token
		claims, err := m.autenticacaoService.ValidarToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Token inválido: " + err.Error()})
			c.Abort()
			return
		}

		// Armazenar as claims no contexto para uso posterior
		c.Set("usuarioID", claims["id"])
		c.Set("usuarioEmail", claims["email"])
		c.Set("usuarioNome", claims["nome"])

		c.Next()
	}
}
