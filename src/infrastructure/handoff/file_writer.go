package handoff

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"pf-ai/src/domain"
)

type FileWriter struct {
	Dir string
}

func (w FileWriter) WriteHandoff(ctx context.Context, item domain.Handoff) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	if err := item.Validate(); err != nil {
		return "", err
	}
	if err := os.MkdirAll(w.Dir, 0o755); err != nil {
		return "", err
	}

	filename := safeFilename(item.Header.Timestamp + "_" + string(item.Header.Intent) + "_" + item.Header.TaskRef + ".json")
	path := filepath.Join(w.Dir, filename)

	payload, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(path, payload, 0o600); err != nil {
		return "", err
	}

	return path, nil
}

func safeFilename(value string) string {
	value = strings.ReplaceAll(value, ":", "")
	value = strings.ReplaceAll(value, " ", "_")
	value = strings.ReplaceAll(value, string(filepath.Separator), "_")
	value = strings.ReplaceAll(value, "/", "_")
	value = strings.ReplaceAll(value, "\\", "_")
	return value
}
