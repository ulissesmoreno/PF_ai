package filememory

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"pf-ai/src/domain"
)

type Reader struct {
	Root         string
	AllowedPath map[string]bool
}

func NewReader(root string, allowed []string) Reader {
	allowedPath := make(map[string]bool, len(allowed))
	for _, item := range allowed {
		allowedPath[filepath.ToSlash(filepath.Clean(item))] = true
	}

	return Reader{
		Root:         root,
		AllowedPath: allowedPath,
	}
}

func (r Reader) ReadMemory(ctx context.Context, memory domain.MemoryFile) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	if !r.AllowedPath[memory.Path] {
		return "", errors.New("memory path is not allowlisted")
	}

	root, err := filepath.Abs(r.Root)
	if err != nil {
		return "", err
	}

	target, err := filepath.Abs(filepath.Join(root, filepath.FromSlash(memory.Path)))
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(root, target)
	if err != nil {
		return "", err
	}
	if rel == ".." || filepath.IsAbs(rel) || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", errors.New("memory path escapes project root")
	}

	bytes, err := os.ReadFile(target)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
