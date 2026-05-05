package domain

import (
	"errors"
	"path/filepath"
	"strings"
)

type MemoryFile struct {
	Path  string `json:"path"`
	Label string `json:"label"`
}

func NewMemoryFile(path, label string) (MemoryFile, error) {
	cleanPath := filepath.ToSlash(filepath.Clean(strings.TrimSpace(path)))
	memory := MemoryFile{
		Path:  cleanPath,
		Label: strings.TrimSpace(label),
	}

	if memory.Path == "." || memory.Path == "" {
		return MemoryFile{}, errors.New("memory path is required")
	}
	if filepath.IsAbs(memory.Path) || strings.HasPrefix(memory.Path, "../") || strings.Contains(memory.Path, "/../") {
		return MemoryFile{}, errors.New("memory path must stay inside project root")
	}
	if memory.Label == "" {
		return MemoryFile{}, errors.New("memory label is required")
	}

	return memory, nil
}
