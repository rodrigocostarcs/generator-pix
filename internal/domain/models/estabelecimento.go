package models

import "time"

// Estabelecimento representa a entidade de estabelecimento
type Estabelecimento struct {
	ID           string    `json:"id"`
	Nome         string    `json:"nome"`
	Descricao    *string   `json:"descricao,omitempty"`
	Email        string    `json:"email"`
	Senha        string    `json:"-"` // Nunca retornar a senha no JSON
	Ativo        bool      `json:"ativo"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

// EstabelecimentoRequest representa os dados de entrada para criação de um estabelecimento
// swagger:model
type EstabelecimentoRequest struct {
	// Nome do estabelecimento
	// required: true
	// example: Loja do José
	Nome string `json:"nome" binding:"required"`

	// Descrição do estabelecimento (opcional)
	// example: Loja de produtos diversos
	Descricao *string `json:"descricao,omitempty"`

	// Email do estabelecimento (usado para login)
	// required: true
	// example: contato@lojadojose.com.br
	Email string `json:"email" binding:"required,email"`

	// Senha do estabelecimento (mínimo 6 caracteres)
	// required: true
	// example: senha123
	// min length: 6
	Senha string `json:"senha" binding:"required,min=6"`
}

// EstabelecimentoResponse representa a resposta após a criação de um estabelecimento
// swagger:model
type EstabelecimentoResponse struct {
	// ID único do estabelecimento
	// example: 123e4567-e89b-12d3-a456-426614174000
	ID string `json:"id"`

	// Nome do estabelecimento
	// example: Loja do José
	Nome string `json:"nome"`

	// Descrição do estabelecimento (opcional)
	// example: Loja de produtos diversos
	Descricao *string `json:"descricao,omitempty"`

	// Email do estabelecimento
	// example: contato@lojadojose.com.br
	Email string `json:"email"`

	// Status de ativação do estabelecimento
	// example: true
	Ativo bool `json:"ativo"`

	// Data de criação
	// example: 2023-01-01T12:00:00Z
	CriadoEm time.Time `json:"criado_em"`

	// Data da última atualização
	// example: 2023-01-01T12:00:00Z
	AtualizadoEm time.Time `json:"atualizado_em"`
}

// LoginRequest representa os dados de entrada para login
// swagger:model
type LoginRequest struct {
	// Email do estabelecimento
	// required: true
	// example: contato@lojadojose.com.br
	Email string `json:"email" binding:"required,email"`

	// Senha do estabelecimento
	// required: true
	// example: senha123
	Senha string `json:"senha" binding:"required"`
}

// LoginResponse representa a resposta após o login bem-sucedido
// swagger:model
type LoginResponse struct {
	// Token JWT para autenticação nas rotas protegidas
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	Token string `json:"token"`

	// Informações do estabelecimento autenticado
	Estabelecimento EstabelecimentoResponse `json:"estabelecimento"`
}
