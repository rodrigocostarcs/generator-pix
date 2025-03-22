package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rodrigocostarcs/pix-generator/internal/domain/models"
	"golang.org/x/crypto/bcrypt"
)

// AutenticacaoService contém a lógica para autenticação
type AutenticacaoService struct {
	jwtChaveSecreta []byte
}

// NewAutenticacaoService cria uma nova instância do serviço de autenticação
func NewAutenticacaoService() *AutenticacaoService {
	// Obter a chave secreta para JWT do ambiente (ou criar uma padrão)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "c8b7a3e5f9d2b6a1c4d7e9f3b2a5c8d9e6f3a2b5c8d7e9f3b6a5c8d9e6f3b2a5"
	}

	return &AutenticacaoService{
		jwtChaveSecreta: []byte(jwtSecret),
	}
}

// VerificarSenha verifica se a senha fornecida corresponde ao hash armazenado
func (s *AutenticacaoService) VerificarSenha(hashSenha, senha string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashSenha), []byte(senha))
}

// GerarToken gera um token JWT para um estabelecimento
func (s *AutenticacaoService) GerarToken(estabelecimento models.Estabelecimento) (string, error) {
	// Definir o tempo de expiração
	tempoExpiracao := time.Now().Add(24 * time.Hour)

	// Criar as claims (payload) do token
	claims := jwt.MapClaims{
		"id":    estabelecimento.ID,
		"email": estabelecimento.Email,
		"nome":  estabelecimento.Nome,
		"exp":   tempoExpiracao.Unix(),
	}

	// Criar o token com as claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Assinar o token com a chave secreta
	tokenString, err := token.SignedString(s.jwtChaveSecreta)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidarToken valida um token JWT e retorna as claims
func (s *AutenticacaoService) ValidarToken(tokenString string) (jwt.MapClaims, error) {
	// Analisar o token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verificar se o método de assinatura é o esperado
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return s.jwtChaveSecreta, nil
	})

	if err != nil {
		return nil, err
	}

	// Verificar se o token é válido
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}
