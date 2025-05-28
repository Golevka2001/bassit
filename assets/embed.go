package assets

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed *
var Assets embed.FS

func ExtractTo(path string) error {
	// Configuration file
	cfgFile := "config.yaml"
	cfgFileDst := filepath.Join(path, "config.yaml")
	if err := extractFile(cfgFile, cfgFileDst, 0644); err != nil {
		return fmt.Errorf("failed to extract configuration file: %w", err)
	}

	// TODO: Rubberband binary

	// TODO: Audio files

	return nil
}

func extractFile(src, dst string, perm os.FileMode) error {
	// Check if the file exists
	if _, err := os.Stat(dst); err == nil {
		return nil
	}

	data, err := Assets.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read embedded file: %w", err)
	}
	if err := os.WriteFile(dst, data, perm); err != nil {
		return fmt.Errorf("failed to write file to %s: %w", dst, err)
	}
	return nil
}
