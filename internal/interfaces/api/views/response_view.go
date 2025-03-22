package views

import "github.com/gin-gonic/gin"

// ResponseView encapsula a lógica de formatação de respostas da API
type ResponseView struct{}

// Response representa a estrutura padrão de resposta da API
// swagger:model
type Response struct {
	// Indica se a requisição foi bem-sucedida
	// example: true
	Success bool `json:"success"`

	// Dados de resposta (opcional)
	Data interface{} `json:"data,omitempty"`

	// Mensagem de erro (apenas quando Success = false)
	// example: Credenciais inválidas
	Error string `json:"error,omitempty"`
}

// NewResponseView cria uma nova instância de ResponseView
func NewResponseView() *ResponseView {
	return &ResponseView{}
}

// Success retorna uma resposta de sucesso formatada
func (rv *ResponseView) Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

// Error retorna uma resposta de erro formatada
func (rv *ResponseView) Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   message,
	})
}

// Download prepara o contexto para um download de arquivo
func (rv *ResponseView) Download(c *gin.Context, filename string, mimeType string, data []byte) {
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", mimeType)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	c.Data(200, mimeType, data)
}
