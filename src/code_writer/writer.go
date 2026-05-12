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
	filePath  string
	reInvalid = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)
	reEspacos = regexp.MustCompile(`\s+`)
)

type CodeFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// Init define o caminho do vault do Obsidian.
func Init(path string) {
	filePath = path
}

// SaveFile cria um arquivo .md no vault com frontmatter YAML.
// Se a nota já existir, adiciona sufixo de timestamp para não sobrescrever.
func SaveFile(titulo, source, conteudo string) error {
	if filePath == "" {
		return fmt.Errorf("vault não inicializado — chame obsidian.Init()")
	}

	nomeArquivo := sanitizarNome(titulo)
	caminho := filepath.Join(filePath, nomeArquivo+".md")

	// Evitar sobrescrever nota existente
	if _, err := os.Stat(caminho); err == nil {
		ts := time.Now().Format("20060102_150405")
		caminho = filepath.Join(filePath, fmt.Sprintf("%s_%s.md", nomeArquivo, ts))
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
func WriteCodeFiles(baseDir string, files []CodeFile) ([]string, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("nenhum arquivo informado para escrita")
	}

	baseAbs, err := filepath.Abs(baseDir)
	if err != nil {
		return nil, fmt.Errorf("resolver workspace: %w", err)
	}

	written := make([]string, 0, len(files))
	for _, file := range files {
		if strings.TrimSpace(file.Path) == "" {
			return written, fmt.Errorf("arquivo sem path")
		}

		target := file.Path
		if !filepath.IsAbs(target) {
			target = filepath.Join(baseAbs, target)
		}

		targetAbs, err := filepath.Abs(filepath.Clean(target))
		if err != nil {
			return written, fmt.Errorf("resolver path %q: %w", file.Path, err)
		}

		rel, err := filepath.Rel(baseAbs, targetAbs)
		if err != nil {
			return written, fmt.Errorf("validar path %q: %w", file.Path, err)
		}
		if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) || filepath.IsAbs(rel) {
			return written, fmt.Errorf("path fora do workspace: %s", file.Path)
		}

		if err := os.MkdirAll(filepath.Dir(targetAbs), 0755); err != nil {
			return written, fmt.Errorf("criar diretório de %s: %w", file.Path, err)
		}
		if err := os.WriteFile(targetAbs, []byte(file.Content), 0644); err != nil {
			return written, fmt.Errorf("escrever %s: %w", file.Path, err)
		}
		written = append(written, targetAbs)
	}

	return written, nil
}

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
