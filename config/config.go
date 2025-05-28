package config

import (
	_ "embed"
	"fmt"
	"os"

	A "bassit/assets"
)

func WriteDefaultConfig(path string) error {
	// Write the embedded default configuration to the specified path
	data, err := A.Assets.ReadFile("config.yaml")
	if err != nil {
		return fmt.Errorf("failed to read embedded config: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write default config to %s: %w", path, err)
	}
	return nil
}
