package code_writer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	vaultPath  string
	reInvalid  = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)
	reEspacos  = regexp.MustCompile(`\s+`)
)

// Init define o caminho do vault do Obsidian.
func Init(path string) {
	vaultPath = path
}

// SalvarNota cria um arquivo .md no vault com frontmatter YAML.
// Se a nota já existir, adiciona sufixo de timestamp para não sobrescrever.
func SalvarNota(titulo, source, conteudo string) error {
	if vaultPath == "" {
		return fmt.Errorf("vault não inicializado — chame obsidian.Init()")
	}

	nomeArquivo := sanitizarNome(titulo)
	caminho := filepath.Join(vaultPath, nomeArquivo+".md")

	// Evitar sobrescrever nota existente
	if _, err := os.Stat(caminho); err == nil {
		ts := time.Now().Format("20060102_150405")
		caminho = filepath.Join(vaultPath, fmt.Sprintf("%s_%s.md", nomeArquivo, ts))
	}

	frontmatter := fmt.Sprintf(`---
title: "%s"
date: %s
source: "%s"
tags:
  - gerado-por-ia
---

`, titulo, time.Now().Format("2006-01-02"), source)

	conteudoFinal := frontmatter + conteudo

	if err := os.WriteFile(caminho, []byte(conteudoFinal), 0644); err != nil {
		return fmt.Errorf("salvar nota %s: %w", caminho, err)
	}

	return nil
}

// sanitizarNome converte o título em nome de arquivo seguro para todos os OSes.
func sanitizarNome(titulo string) string {
	// Remover caracteres inválidos
	nome := reInvalid.ReplaceAllString(titulo, "")
	// Normalizar espaços
	nome = reEspacos.ReplaceAllString(nome, "_")
	// Remover underscores no início/fim
	nome = strings.Trim(nome, "_")
	// Limitar tamanho (255 chars é o máximo em maioria dos FSes)
	if len(nome) > 200 {
		nome = nome[:200]
	}
	if nome == "" {
		nome = fmt.Sprintf("nota_%d", time.Now().Unix())
	}
	return nome
}
