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
type EstabelecimentoRequest struct {
	Nome      string  `json:"nome" binding:"required"`
	Descricao *string `json:"descricao,omitempty"`
	Email     string  `json:"email" binding:"required,email"`
	Senha     string  `json:"senha" binding:"required,min=6"`
}

// EstabelecimentoResponse representa a resposta após a criação de um estabelecimento
type EstabelecimentoResponse struct {
	ID           string    `json:"id"`
	Nome         string    `json:"nome"`
	Descricao    *string   `json:"descricao,omitempty"`
	Email        string    `json:"email"`
	Ativo        bool      `json:"ativo"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

// LoginRequest representa os dados de entrada para login
type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Senha string `json:"senha" binding:"required"`
}

// LoginResponse representa a resposta após o login bem-sucedido
type LoginResponse struct {
	Token           string                  `json:"token"`
	Estabelecimento EstabelecimentoResponse `json:"estabelecimento"`
}
