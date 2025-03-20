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
type PixRequest struct {
	Nome          string   `json:"nome" binding:"required"`
	Chave         string   `json:"chave" binding:"required"`
	Cidade        string   `json:"cidade" binding:"required"`
	Valor         *float64 `json:"valor,omitempty"`
	Identificador *string  `json:"identificador,omitempty"`
	Descricao     *string  `json:"descricao,omitempty"`
}

// PixResponse representa a resposta após a geração de um PIX
type PixResponse struct {
	CodigoPix string `json:"codigo_pix"`
	QRCodeSVG string `json:"qrcode_svg"`
	QRCodePNG string `json:"qrcode_png"`
}