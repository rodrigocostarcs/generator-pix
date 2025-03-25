package services_test

import (
	"testing"

	"github.com/rodrigocostarcs/pix-generator/internal/domain/models"
	"github.com/rodrigocostarcs/pix-generator/internal/domain/services"
	"github.com/stretchr/testify/assert"
)

// TestGerarPixEstatico testa a geração de PIX estático
func TestGerarPixEstatico(t *testing.T) {
	// Inicializar o serviço
	service := services.NewPixGeneratorService()

	t.Run("DadosCompletos", func(t *testing.T) {
		// Preparar os dados de entrada
		valor := 100.50
		identificador := "FATURA123"
		descricao := "PAGAMENTO DE SERVICOS"

		req := models.PixRequest{
			Nome:          "JOSE DA SILVA",
			Chave:         "josesilva@email.com",
			Cidade:        "SAO PAULO",
			Valor:         &valor,
			Identificador: &identificador,
			Descricao:     &descricao,
		}

		// Executar o método a ser testado
		response, err := service.GerarPixEstatico(req)

		// Verificar resultados
		assert.NoError(t, err)
		assert.NotEmpty(t, response.CodigoPix)
		assert.NotEmpty(t, response.QRCodeSVG)
		assert.NotEmpty(t, response.QRCodePNG)

		// Verificar se o código PIX contém elementos esperados
		assert.Contains(t, response.CodigoPix, "BR.GOV.BCB.PIX")
		assert.Contains(t, response.CodigoPix, "josesilva@email.com")
		assert.Contains(t, response.CodigoPix, "JOSE DA SILVA")
		assert.Contains(t, response.CodigoPix, "SAO PAULO")
		assert.Contains(t, response.CodigoPix, "FATURA123")

		// Verificar se o QR code PNG está no formato base64 correto
		assert.Contains(t, response.QRCodePNG, "data:image/png;base64,")
	})

	t.Run("DadosMinimos", func(t *testing.T) {
		// Testar com dados mínimos obrigatórios
		req := models.PixRequest{
			Nome:   "MARIA OLIVEIRA",
			Chave:  "maria@email.com",
			Cidade: "RIO DE JANEIRO",
		}

		// Executar o método a ser testado
		response, err := service.GerarPixEstatico(req)

		// Verificar resultados
		assert.NoError(t, err)
		assert.NotEmpty(t, response.CodigoPix)
		assert.NotEmpty(t, response.QRCodeSVG)
		assert.NotEmpty(t, response.QRCodePNG)

		// Verificar se o código PIX contém elementos esperados
		assert.Contains(t, response.CodigoPix, "MARIA OLIVEIRA")
		assert.Contains(t, response.CodigoPix, "RIO DE JANEIRO")
		assert.Contains(t, response.CodigoPix, "maria@email.com")

		// Verificar que o código não contém valor (opcional)
		assert.NotContains(t, response.CodigoPix, "54")
	})

	t.Run("NomeMuitoLongo", func(t *testing.T) {
		// Testar com nome muito longo (deve ser truncado)
		req := models.PixRequest{
			Nome:   "NOME MUITO LONGO QUE ULTRAPASSA O LIMITE DE CARACTERES PERMITIDO PELO PIX",
			Chave:  "nome@email.com",
			Cidade: "BELO HORIZONTE",
		}

		// Executar o método a ser testado
		response, err := service.GerarPixEstatico(req)

		// Verificar resultados
		assert.NoError(t, err)
		assert.NotEmpty(t, response.CodigoPix)

		// Verificar que o nome foi truncado (limite 25 caracteres)
		nomeTruncado := "NOME MUITO LONGO QUE ULTR"
		assert.Contains(t, response.CodigoPix, nomeTruncado)
	})
}
