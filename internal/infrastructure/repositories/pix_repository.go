package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/rodrigocostarcs/pix-generator/internal/domain/models"
)

// PixRepository interface para persistência de dados PIX
type PixRepository interface {
	Save(pix models.Pix) (uint, error)
	FindByID(id uint) (models.Pix, error)
	List() ([]models.Pix, error)
	FindByCodigoPix(codigoPix string) (models.Pix, error)
}

// MysqlPixRepository implementação MySQL do repositório PIX
type MysqlPixRepository struct {
	db *sql.DB
}

// NewMysqlPixRepository cria uma nova instância do repositório MySQL
func NewMysqlPixRepository(db *sql.DB) *MysqlPixRepository {
	return &MysqlPixRepository{db: db}
}

// Save salva um código PIX no banco de dados
func (r *MysqlPixRepository) Save(pix models.Pix) (uint, error) {
	query := `
		INSERT INTO pix (nome, chave, cidade, valor, identificador, descricao, codigo_pix, qrcode_svg, qrcode_png, criado_em)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		pix.Nome,
		pix.Chave,
		pix.Cidade,
		pix.Valor,
		pix.Identificador,
		pix.Descricao,
		pix.CodigoPix,
		pix.QRCodeSVG,
		pix.QRCodePNG,
		time.Now(),
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

// FindByID busca um PIX pelo ID
func (r *MysqlPixRepository) FindByID(id uint) (models.Pix, error) {
	var pix models.Pix

	query := `
		SELECT id, nome, chave, cidade, valor, identificador, descricao, codigo_pix, qrcode_svg, qrcode_png, criado_em
		FROM pix
		WHERE id = ?
	`

	var valor sql.NullFloat64
	var identificador, descricao sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&pix.ID,
		&pix.Nome,
		&pix.Chave,
		&pix.Cidade,
		&valor,
		&identificador,
		&descricao,
		&pix.CodigoPix,
		&pix.QRCodeSVG,
		&pix.QRCodePNG,
		&pix.CriadoEm,
	)

	if err != nil {
		return models.Pix{}, err
	}

	if valor.Valid {
		pix.Valor = &valor.Float64
	}

	if identificador.Valid {
		pix.Identificador = &identificador.String
	}

	if descricao.Valid {
		pix.Descricao = &descricao.String
	}

	return pix, nil
}

// List lista todos os PIX gerados
func (r *MysqlPixRepository) List() ([]models.Pix, error) {
	var pixList []models.Pix

	query := `
		SELECT id, nome, chave, cidade, valor, identificador, descricao, codigo_pix, qrcode_svg, qrcode_png, criado_em
		FROM pix
		ORDER BY criado_em DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pix models.Pix
		var valor sql.NullFloat64
		var identificador, descricao sql.NullString

		err := rows.Scan(
			&pix.ID,
			&pix.Nome,
			&pix.Chave,
			&pix.Cidade,
			&valor,
			&identificador,
			&descricao,
			&pix.CodigoPix,
			&pix.QRCodeSVG,
			&pix.QRCodePNG,
			&pix.CriadoEm,
		)

		if err != nil {
			return nil, err
		}

		if valor.Valid {
			pix.Valor = &valor.Float64
		}

		if identificador.Valid {
			pix.Identificador = &identificador.String
		}

		if descricao.Valid {
			pix.Descricao = &descricao.String
		}

		pixList = append(pixList, pix)
	}

	return pixList, nil
}

// FindByCodigoPix busca um PIX pelo código gerado
func (r *MysqlPixRepository) FindByCodigoPix(codigoPix string) (models.Pix, error) {
	var pix models.Pix

	query := `
        SELECT id, nome, chave, cidade, valor, identificador, descricao, codigo_pix, qrcode_svg, qrcode_png, criado_em
        FROM pix
        WHERE codigo_pix = ?
        LIMIT 1
    `

	var valor sql.NullFloat64
	var identificador, descricao sql.NullString

	err := r.db.QueryRow(query, codigoPix).Scan(
		&pix.ID,
		&pix.Nome,
		&pix.Chave,
		&pix.Cidade,
		&valor,
		&identificador,
		&descricao,
		&pix.CodigoPix,
		&pix.QRCodeSVG,
		&pix.QRCodePNG,
		&pix.CriadoEm,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Pix{}, fmt.Errorf("código PIX não encontrado")
		}
		return models.Pix{}, err
	}

	if valor.Valid {
		pix.Valor = &valor.Float64
	}

	if identificador.Valid {
		pix.Identificador = &identificador.String
	}

	if descricao.Valid {
		pix.Descricao = &descricao.String
	}

	return pix, nil
}
