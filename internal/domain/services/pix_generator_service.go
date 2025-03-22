package services

import (
	"encoding/base64"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/rodrigocostarcs/pix-generator/internal/domain/models"
	"github.com/skip2/go-qrcode"
)

// PixGeneratorService contém a lógica para geração de códigos PIX
type PixGeneratorService struct{}

// NewPixGeneratorService cria uma nova instância do serviço de geração de PIX
func NewPixGeneratorService() *PixGeneratorService {
	return &PixGeneratorService{}
}

// GerarPixEstatico gera um código PIX estático com base nos parâmetros fornecidos
func (s *PixGeneratorService) GerarPixEstatico(req models.PixRequest) (models.PixResponse, error) {
	// Validações básicas
	nome := removerCaracteresEspeciais(limitarTamanho(req.Nome, 25))
	cidade := limitarTamanho(removerCaracteresEspeciais(req.Cidade), 15)
	chave := req.Chave

	// Tratar identificador e descrição
	var identificador, descricao string
	if req.Identificador != nil {
		identificador = limitarTamanho(removerCaracteresEspeciais(*req.Identificador), 25)
	}
	if req.Descricao != nil {
		descricao = limitarTamanho(removerCaracteresEspeciais(*req.Descricao), 50)
	}

	// Formatar valor (se fornecido)
	var valorFormatado string
	if req.Valor != nil && *req.Valor > 0 {
		valorFormatado = formatarValor(*req.Valor)
	}

	// Construir o payload PIX
	payload := s.construirPayloadPix(nome, chave, cidade, valorFormatado, identificador, descricao)

	// Adicionar CRC
	codigoPix := payload + "6304" + s.calcularCRC(payload+"6304")

	// Gerar QR Codes
	qrSVG, err := s.gerarQRCodeSVG(codigoPix)
	if err != nil {
		return models.PixResponse{}, err
	}

	qrPNG, err := s.GerarQRCodePNG(codigoPix)
	if err != nil {
		return models.PixResponse{}, err
	}

	return models.PixResponse{
		CodigoPix: codigoPix,
		QRCodeSVG: qrSVG,
		QRCodePNG: qrPNG,
	}, nil
}

// construirPayloadPix monta o payload do PIX conforme especificações do Banco Central
func (s *PixGeneratorService) construirPayloadPix(nome, chave, cidade, valor, identificador, descricao string) string {
	var payload strings.Builder

	// Payload Format Indicator (obrigatório)
	addCampo(&payload, "00", "01")

	// Point of Initiation Method - 11 para estático, 12 para dinâmico
	addCampo(&payload, "01", "11")

	// Merchant Account Information (obrigatório)
	addCampo(&payload, "26", s.pixGUI(chave))

	// Merchant Category Code (obrigatório)
	addCampo(&payload, "52", "0000")

	// Transaction Currency (BRL = 986) (obrigatório)
	addCampo(&payload, "53", "986")

	// Country Code (obrigatório)
	addCampo(&payload, "58", "BR")

	// Merchant Name (obrigatório)
	addCampo(&payload, "59", nome)

	// Merchant City (obrigatório)
	addCampo(&payload, "60", cidade)

	// Transaction Amount (opcional)
	if valor != "" {
		addCampo(&payload, "54", valor)
	}

	// Additional Data Field (obrigatório para algumas implementações)
	campoAdicional := s.adicionarCampoAdicional(identificador, descricao)
	if campoAdicional != "" {
		addCampo(&payload, "62", campoAdicional)
	}

	return payload.String()
}

// pixGUI gera o GUI do PIX
func (s *PixGeneratorService) pixGUI(chave string) string {
	var gui strings.Builder

	// GUI do PIX
	addCampo(&gui, "00", "BR.GOV.BCB.PIX")

	// Chave PIX
	addCampo(&gui, "01", chave)

	return gui.String()
}

// adicionarCampoAdicional adiciona campos adicionais ao PIX
func (s *PixGeneratorService) adicionarCampoAdicional(identificador, descricao string) string {
	var campoAdicional strings.Builder

	// Adiciona Reference Label (05) se fornecido
	if identificador != "" {
		addCampo(&campoAdicional, "05", identificador)
	} else {
		addCampo(&campoAdicional, "05", "***")
	}

	// Adiciona Purpose of Transaction (08) se fornecido
	if descricao != "" {
		addCampo(&campoAdicional, "08", descricao)
	}

	return campoAdicional.String()
}

// calcularCRC calcula o CRC-16/CCITT-FALSE conforme especificação do Bacen
func (s *PixGeneratorService) calcularCRC(payload string) string {
	crc := uint16(0xFFFF)
	polynomial := uint16(0x1021)

	for i := 0; i < len(payload); i++ {
		crc ^= uint16(payload[i]) << 8

		for j := 0; j < 8; j++ {
			if (crc & 0x8000) != 0 {
				crc = (crc << 1) ^ polynomial
			} else {
				crc = crc << 1
			}
		}
	}

	return fmt.Sprintf("%04X", crc)
}

// gerarQRCodeSVG gera um QR code em formato SVG
func (s *PixGeneratorService) gerarQRCodeSVG(codigoPix string) (string, error) {
	// Implementação simplificada - em produção você precisaria de uma biblioteca mais completa para SVG
	// Esta é uma simulação
	svg := fmt.Sprintf("<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 250 250'><text x='10' y='100'>QR Code SVG para: %s</text></svg>", codigoPix)
	return svg, nil
}

// GerarQRCodePNG gera um QR code em formato PNG (base64)
func (s *PixGeneratorService) GerarQRCodePNG(codigoPix string) (string, error) {
	qrCode, err := qrcode.Encode(codigoPix, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	base64PNG := "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrCode)
	return base64PNG, nil
}

// Funções auxiliares

// addCampo adiciona um campo EMV ao payload
func addCampo(builder *strings.Builder, id, valor string) {
	if valor == "" {
		return
	}

	tamanho := fmt.Sprintf("%02d", len(valor))
	builder.WriteString(id)
	builder.WriteString(tamanho)
	builder.WriteString(valor)
}

// removerCaracteresEspeciais remove acentos e caracteres especiais
func removerCaracteresEspeciais(texto string) string {
	// Remover caracteres não alfanuméricos exceto espaços
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	texto = reg.ReplaceAllString(texto, "")

	// Converter para maiúsculas e trim
	return strings.ToUpper(strings.TrimSpace(texto))
}

// limitarTamanho limita o tamanho de uma string
func limitarTamanho(texto string, tamanhoMax int) string {
	if len(texto) <= tamanhoMax {
		return texto
	}
	return texto[:tamanhoMax]
}

// formatarValor formata um valor float para o formato exigido pelo PIX
func formatarValor(valor float64) string {
	// Arredondar para 2 casas decimais
	valorArredondado := math.Round(valor*100) / 100

	// Converter para string com 2 casas decimais
	return strconv.FormatFloat(valorArredondado, 'f', 2, 64)
}
