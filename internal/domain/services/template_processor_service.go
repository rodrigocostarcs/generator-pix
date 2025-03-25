package services

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

// TemplatePosition define a posição e tamanho do QR code no template
type TemplatePosition struct {
	X        int    // Posição X do canto superior esquerdo
	Y        int    // Posição Y do canto superior esquerdo
	Size     int    // Tamanho do QR code
	Filename string // Nome do arquivo de template
}

// TemplateProcessor processa os templates para QR code
type TemplateProcessor struct {
	templatesDir string
	templates    map[string]TemplatePosition
}

// NewTemplateProcessor cria uma nova instância do processador de templates
func NewTemplateProcessor(templatesDir string) *TemplateProcessor {
	// Configuração dos templates e posições dos QR codes
	templates := map[string]TemplatePosition{
		"template_pix_1": {
			X:        250, // Posição horizontal mais centralizada
			Y:        400, // Posição vertical melhor
			Size:     200, // Tamanho grande o suficiente para ser legível
			Filename: "template_pix_1.png",
		},
	}

	log.Printf("Inicializando TemplateProcessor com diretório: %s", templatesDir)
	log.Printf("Templates registrados: %v", templates)

	return &TemplateProcessor{
		templatesDir: templatesDir,
		templates:    templates,
	}
}

// ApplyTemplate aplica um template a um código QR existente
func (p *TemplateProcessor) ApplyTemplate(qrCodePNG string, templateName string) ([]byte, error) {
	// Verificar se o template existe
	log.Printf("Solicitação para aplicar template: %s", templateName)

	templateConfig, exists := p.templates[templateName]
	if !exists {
		log.Printf("Template não encontrado no mapa de templates: %s", templateName)
		log.Printf("Templates disponíveis: %v", p.templates)
		return nil, fmt.Errorf("template não encontrado: %s", templateName)
	}

	// Caminho completo para o arquivo de template
	templatePath := filepath.Join(p.templatesDir, templateConfig.Filename)
	log.Printf("Caminho completo do template: %s", templatePath)

	// Verificar se o arquivo existe
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		log.Printf("Arquivo de template não encontrado no caminho: %s", templatePath)

		// Tentar listar arquivos no diretório para debug
		files, err := os.ReadDir(p.templatesDir)
		if err != nil {
			log.Printf("Erro ao listar diretório: %v", err)
		} else {
			log.Printf("Arquivos no diretório %s:", p.templatesDir)
			for _, file := range files {
				log.Printf(" - %s", file.Name())
			}
		}

		return nil, errors.New("arquivo de template não encontrado")
	}

	log.Printf("Arquivo de template encontrado, iniciando processamento...")

	// Decodificar o QR code de base64
	base64Data := strings.TrimPrefix(qrCodePNG, "data:image/png;base64,")
	decodedData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		log.Printf("Erro ao decodificar QR code de base64: %v", err)
		return nil, err
	}

	// Converter para imagem
	qrImg, err := png.Decode(bytes.NewReader(decodedData))
	if err != nil {
		log.Printf("Erro ao decodificar PNG do QR code: %v", err)
		return nil, err
	}

	// Carregar o template
	templateFile, err := os.Open(templatePath)
	if err != nil {
		log.Printf("Erro ao abrir arquivo de template: %v", err)
		return nil, err
	}
	defer templateFile.Close()

	templateImg, err := png.Decode(templateFile)
	if err != nil {
		log.Printf("Erro ao decodificar PNG do template: %v", err)
		return nil, err
	}

	// Criar a imagem resultante
	bounds := templateImg.Bounds()
	resultImg := image.NewRGBA(bounds)

	// Desenhar o template
	draw.Draw(resultImg, bounds, templateImg, image.Point{}, draw.Over)

	// Calcular o tamanho real do QR Code para preservar sua proporção
	qrSize := qrImg.Bounds().Dx()

	// Calcular o fator de escala necessário para ajustar o QR code ao espaço desejado
	// Reduzir um pouco o tamanho para garantir margens
	scaleFactor := float64(templateConfig.Size) / float64(qrSize) * 0.90 // 90% do tamanho para garantir margem

	// Calcular o novo tamanho
	newSize := int(float64(qrSize) * scaleFactor)

	// Calcular as margens adicionais para centralizar dentro do espaço alocado
	marginX := (templateConfig.Size - newSize) / 2
	marginY := (templateConfig.Size - newSize) / 2

	// Definir a área para a zona quieta (margem branca) ao redor do QR code
	quietZone := 20 // pixels de margem branca
	whiteArea := image.Rect(
		templateConfig.X-quietZone,
		templateConfig.Y-quietZone,
		templateConfig.X+templateConfig.Size+quietZone,
		templateConfig.Y+templateConfig.Size+quietZone,
	)

	// Desenhar um fundo branco para a zona quieta
	whiteColor := image.NewUniform(color.White)
	draw.Draw(resultImg, whiteArea, whiteColor, image.Point{}, draw.Src)

	// Redimensionar o QR code para o novo tamanho
	resizedQR := resize.Resize(uint(newSize), uint(newSize), qrImg, resize.Lanczos3)

	// Calcular a posição final para o QR code redimensionado
	qrRect := image.Rect(
		templateConfig.X+marginX,
		templateConfig.Y+marginY,
		templateConfig.X+marginX+newSize,
		templateConfig.Y+marginY+newSize,
	)

	// Desenhar o QR code redimensionado
	draw.Draw(resultImg, qrRect, resizedQR, image.Point{}, draw.Over)

	log.Printf("Composição de imagem concluída, codificando resultado...")

	// Converter para PNG com configuração de qualidade máxima
	var buf bytes.Buffer
	encoder := png.Encoder{
		CompressionLevel: png.NoCompression, // Usar sem compressão para máxima qualidade
	}
	if err := encoder.Encode(&buf, resultImg); err != nil {
		log.Printf("Erro ao codificar imagem resultante: %v", err)
		return nil, err
	}

	log.Printf("Template aplicado com sucesso")
	return buf.Bytes(), nil
}
