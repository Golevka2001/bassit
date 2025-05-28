package constant

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	// baseDir is used to store configuration files and other resources. It defaults to `$HOME/.config/bassit`
	baseDir = ""
	BaseDir = func() string {
		if baseDir != "" {
			return baseDir
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			// Fallback to current working directory if home directory cannot be determined
			homeDir, _ = os.Getwd()
		}
		baseDir = filepath.Join(homeDir, ".config", "bassit")
		return baseDir
	}()
)

func SetBaseDir(path string) error {
	if path == "" {
		baseDir = BaseDir
		// Check if the base directory exists
		_, err := os.Stat(baseDir)
		if err != nil {
			if os.IsNotExist(err) {
				// Create the base directory if it does not exist
				err = os.MkdirAll(baseDir, 0755)
				if err != nil {
					return fmt.Errorf("error creating base directory: %s, %v", baseDir, err)
				}
				return nil
			} else if os.IsPermission(err) {
				return fmt.Errorf("permission denied for base directory: %s", baseDir)
			}
			return fmt.Errorf("error checking base directory: %s, %v", baseDir, err)
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
		baseDir = path
	}

	return nil
}
