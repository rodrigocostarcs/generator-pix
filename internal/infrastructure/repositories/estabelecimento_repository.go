package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rodrigocostarcs/pix-generator/internal/domain/models"
	"golang.org/x/crypto/bcrypt"
)

// EstabelecimentoRepository interface para persistência de estabelecimentos
type EstabelecimentoRepository interface {
	Salvar(estabelecimento models.EstabelecimentoRequest) (models.Estabelecimento, error)
	BuscarPorID(id string) (models.Estabelecimento, error)
	BuscarPorEmail(email string) (models.Estabelecimento, error)
	Listar() ([]models.Estabelecimento, error)
	Atualizar(id string, estabelecimento models.EstabelecimentoRequest) (models.Estabelecimento, error)
	Excluir(id string) error
}

// MysqlEstabelecimentoRepository implementação MySQL do repositório de estabelecimentos
type MysqlEstabelecimentoRepository struct {
	db *sql.DB
}

// NewMysqlEstabelecimentoRepository cria uma nova instância do repositório MySQL para estabelecimentos
func NewMysqlEstabelecimentoRepository(db *sql.DB) *MysqlEstabelecimentoRepository {
	return &MysqlEstabelecimentoRepository{db: db}
}

// Salvar salva um estabelecimento no banco de dados
func (r *MysqlEstabelecimentoRepository) Salvar(req models.EstabelecimentoRequest) (models.Estabelecimento, error) {
	// Gerar UUID para o novo estabelecimento
	id := uuid.New().String()

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Senha), bcrypt.DefaultCost)
	if err != nil {
		return models.Estabelecimento{}, err
	}

	// Data atual
	now := time.Now()

	// Executar query de inserção
	query := `
        INSERT INTO estabelecimentos (id, nome, descricao, email, senha, ativo, criado_em, atualizado_em)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `

	_, err = r.db.Exec(
		query,
		id,
		req.Nome,
		req.Descricao,
		req.Email,
		string(hashedPassword),
		true,
		now,
		now,
	)

	if err != nil {
		return models.Estabelecimento{}, err
	}

	// Criar e retornar o objeto estabelecimento
	estabelecimento := models.Estabelecimento{
		ID:           id,
		Nome:         req.Nome,
		Descricao:    req.Descricao,
		Email:        req.Email,
		Senha:        string(hashedPassword),
		Ativo:        true,
		CriadoEm:     now,
		AtualizadoEm: now,
	}

	return estabelecimento, nil
}

// BuscarPorID busca um estabelecimento pelo ID
func (r *MysqlEstabelecimentoRepository) BuscarPorID(id string) (models.Estabelecimento, error) {
	var estabelecimento models.Estabelecimento
	var descricao sql.NullString

	query := `
        SELECT id, nome, descricao, email, senha, ativo, criado_em, atualizado_em
        FROM estabelecimentos
        WHERE id = ?
    `

	err := r.db.QueryRow(query, id).Scan(
		&estabelecimento.ID,
		&estabelecimento.Nome,
		&descricao,
		&estabelecimento.Email,
		&estabelecimento.Senha,
		&estabelecimento.Ativo,
		&estabelecimento.CriadoEm,
		&estabelecimento.AtualizadoEm,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Estabelecimento{}, fmt.Errorf("estabelecimento não encontrado")
		}
		return models.Estabelecimento{}, err
	}

	if descricao.Valid {
		estabelecimento.Descricao = &descricao.String
	}

	return estabelecimento, nil
}

// BuscarPorEmail busca um estabelecimento pelo email
func (r *MysqlEstabelecimentoRepository) BuscarPorEmail(email string) (models.Estabelecimento, error) {
	var estabelecimento models.Estabelecimento
	var descricao sql.NullString

	query := `
        SELECT id, nome, descricao, email, senha, ativo, criado_em, atualizado_em
        FROM estabelecimentos
        WHERE email = ?
    `

	err := r.db.QueryRow(query, email).Scan(
		&estabelecimento.ID,
		&estabelecimento.Nome,
		&descricao,
		&estabelecimento.Email,
		&estabelecimento.Senha,
		&estabelecimento.Ativo,
		&estabelecimento.CriadoEm,
		&estabelecimento.AtualizadoEm,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Estabelecimento{}, fmt.Errorf("estabelecimento não encontrado")
		}
		return models.Estabelecimento{}, err
	}

	if descricao.Valid {
		estabelecimento.Descricao = &descricao.String
	}

	return estabelecimento, nil
}

// Listar lista todos os estabelecimentos
func (r *MysqlEstabelecimentoRepository) Listar() ([]models.Estabelecimento, error) {
	var estabelecimentos []models.Estabelecimento

	query := `
        SELECT id, nome, descricao, email, senha, ativo, criado_em, atualizado_em
        FROM estabelecimentos
        ORDER BY nome
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var estabelecimento models.Estabelecimento
		var descricao sql.NullString

		err := rows.Scan(
			&estabelecimento.ID,
			&estabelecimento.Nome,
			&descricao,
			&estabelecimento.Email,
			&estabelecimento.Senha,
			&estabelecimento.Ativo,
			&estabelecimento.CriadoEm,
			&estabelecimento.AtualizadoEm,
		)

		if err != nil {
			return nil, err
		}

		if descricao.Valid {
			estabelecimento.Descricao = &descricao.String
		}

		estabelecimentos = append(estabelecimentos, estabelecimento)
	}

	return estabelecimentos, nil
}

// Atualizar atualiza um estabelecimento
func (r *MysqlEstabelecimentoRepository) Atualizar(id string, req models.EstabelecimentoRequest) (models.Estabelecimento, error) {
	// Buscar estabelecimento atual
	estabelecimento, err := r.BuscarPorID(id)
	if err != nil {
		return models.Estabelecimento{}, err
	}

	// Atualizar campos
	estabelecimento.Nome = req.Nome
	estabelecimento.Descricao = req.Descricao
	estabelecimento.Email = req.Email
	estabelecimento.AtualizadoEm = time.Now()

	// Se a senha foi fornecida, atualizar
	if req.Senha != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Senha), bcrypt.DefaultCost)
		if err != nil {
			return models.Estabelecimento{}, err
		}
		estabelecimento.Senha = string(hashedPassword)
	}

	// Executar query de atualização
	query := `
        UPDATE estabelecimentos 
        SET nome = ?, descricao = ?, email = ?, senha = ?, atualizado_em = ?
        WHERE id = ?
    `

	_, err = r.db.Exec(
		query,
		estabelecimento.Nome,
		estabelecimento.Descricao,
		estabelecimento.Email,
		estabelecimento.Senha,
		estabelecimento.AtualizadoEm,
		id,
	)

	if err != nil {
		return models.Estabelecimento{}, err
	}

	return estabelecimento, nil
}

// Excluir remove um estabelecimento (marcando como inativo)
func (r *MysqlEstabelecimentoRepository) Excluir(id string) error {
	query := `
        UPDATE estabelecimentos 
        SET ativo = false, atualizado_em = ?
        WHERE id = ?
    `

	_, err := r.db.Exec(query, time.Now(), id)
	return err
}
