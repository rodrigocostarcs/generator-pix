-- Criar o banco de dados se não existir
CREATE DATABASE IF NOT EXISTS generator_pix;

-- Usar o banco de dados
USE generator_pix;

-- Criar tabela para armazenar os códigos PIX
CREATE TABLE IF NOT EXISTS pix (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    chave VARCHAR(100) NOT NULL,
    cidade VARCHAR(50) NOT NULL,
    valor DECIMAL(10, 2) NULL,
    identificador VARCHAR(100) NULL,
    descricao VARCHAR(200) NULL,
    codigo_pix TEXT NOT NULL,
    qrcode_svg TEXT NOT NULL,
    qrcode_png TEXT NOT NULL,
    criado_em DATETIME NOT NULL
);