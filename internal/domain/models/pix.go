package models

import "time"

// Pix representa a entidade principal do sistema
type Pix struct {
	ID            uint      `json:"id"`
	Nome          string    `json:"nome"`
	Chave         string    `json:"chave"`
	Cidade        string    `json:"cidade"`
	Valor         *float64  `json:"valor,omitempty"`
	Identificador *string   `json:"identificador,omitempty"`
	Descricao     *string   `json:"descricao,omitempty"`
	CodigoPix     string    `json:"codigo_pix"`
	QRCodeSVG     string    `json:"qrcode_svg"`
	QRCodePNG     string    `json:"qrcode_png"`
	CriadoEm      time.Time `json:"criado_em"`
}

// PixRequest representa os dados de entrada para geração de um PIX
// swagger:model
type PixRequest struct {
	// Nome do beneficiário do PIX (obrigatório)
	// required: true
	// example: JOSE DA SILVA
	Nome string `json:"nome" binding:"required"`

	// Chave PIX do beneficiário (obrigatório)
	// required: true
	// example: josesilva@email.com
	Chave string `json:"chave" binding:"required"`

	// Cidade do beneficiário (obrigatório)
	// required: true
	// example: SAO PAULO
	Cidade string `json:"cidade" binding:"required"`

	// Valor da transação (opcional)
	// example: 100.50
	Valor *float64 `json:"valor,omitempty"`

	// Identificador único da transação (opcional)
	// example: FATURA123
	Identificador *string `json:"identificador,omitempty"`

	// Descrição da transação (opcional)
	// example: PAGAMENTO DE SERVICOS
	Descricao *string `json:"descricao,omitempty"`
}

// PixResponse representa a resposta após a geração de um PIX
// swagger:model
type PixResponse struct {
	// Código PIX gerado conforme padrão EMV
	// example: 00020101021126580014BR.GOV.BCB.PIX0136josesilva@email.com5204000053039865802BR5913JOSE DA SILVA6009SAO PAULO62150511FATURA12308103100.506304E5B1
	CodigoPix string `json:"codigo_pix"`

	// QR Code em formato SVG
	// example: <svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 250 250'>...</svg>
	QRCodeSVG string `json:"qrcode_svg"`

	// QR Code em formato PNG (base64)
	// example: data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...
	QRCodePNG string `json:"qrcode_png"`
}
