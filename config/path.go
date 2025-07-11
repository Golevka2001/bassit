package config

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	// BaseDir is used to store configuration files and other resources.
	// It defaults to `$HOME/.config/bassit`
	BaseDir = ""

	// ThemeDir is used to store theme files.
	// It defaults to `${BASSIT_BASE_DIR}/themes`
	ThemeDir = func() string {
		return filepath.Join(BaseDir, "themes")
	}

	// SoundpackDir is used to store soundpack files.
	// It defaults to `${BASSIT_BASE_DIR}/soundpacks`
	SoundpackDir = func() string {
		return filepath.Join(BaseDir, "soundpacks")
	}
)

func SetBaseDir(path string) error {
	if path == "" {
		// If no path is provided, use the home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting home directory: %v", err)
		}
		BaseDir = filepath.Join(homeDir, ".config", "bassit")

		// Check if the base directory exists
		_, err = os.Stat(BaseDir)
		if err != nil {
			if os.IsNotExist(err) {
				// Create the base directory if it does not exist
				err = os.MkdirAll(BaseDir, 0755)
				if err != nil {
					return fmt.Errorf("error creating base directory: %s, %v", BaseDir, err)
				}
				return nil
			} else if os.IsPermission(err) {
				return fmt.Errorf("permission denied for base directory: %s", BaseDir)
			}
			return fmt.Errorf("error checking base directory: %s, %v", BaseDir, err)
		}
	} else {
		// Convert to absolute path if not already
		if !filepath.IsAbs(path) {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return fmt.Errorf("error converting to absolute path: %s, %v", path, err)
			}
			path = absPath
		}

		// Check if the path exists
		_, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("path does not exist: %s", path)
			} else if os.IsPermission(err) {
				return fmt.Errorf("permission denied for path: %s", path)
			}
			return fmt.Errorf("error checking path: %s, %v", path, err)
		}
		BaseDir = path
	}

	return nil
}
