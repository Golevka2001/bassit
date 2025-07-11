package assets

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Golevka2001/bassit/audio"
)

//go:embed *
var Assets embed.FS

// NOTE: The paths of embedded resources must use `/`, even on Windows
const (
	cfgFile = "config.yaml"

	thirdPartyDir = "3rdparty"
	soundpackDir  = "soundpacks"
	themeDir      = "themes"

	rbBinForWindows  = thirdPartyDir + "/rubberband-4.0.0-gpl-executable-windows/rubberband.exe"
	rb3BinForWindows = thirdPartyDir + "/rubberband-4.0.0-gpl-executable-windows/rubberband-r3.exe"
	rbDllForWindows  = thirdPartyDir + "/rubberband-4.0.0-gpl-executable-windows/sndfile.dll"

	rbBinForDarwin  = thirdPartyDir + "/rubberband-4.0.0-gpl-executable-macos/rubberband"
	rb3BinForDarwin = thirdPartyDir + "/rubberband-4.0.0-gpl-executable-macos/rubberband-r3"
)

type FileToExtract struct {
	Src  string
	Dst  string
	Perm os.FileMode
}

// ExtractTo extracts all the embedded resources to the given path
func ExtractTo(path string) error {
	// Configuration file
	files := []FileToExtract{
		{
			Src:  cfgFile,
			Dst:  filepath.Join(path, cfgFile),
			Perm: 0644,
		},
	}

	// Theme files
	themeDstDir := filepath.Join(path, themeDir)
	themeFiles, err := Assets.ReadDir(themeDir)
	if err != nil {
		return fmt.Errorf("failed to read themes directory: %w", err)
	}
	for _, file := range themeFiles {
		files = append(files, FileToExtract{
			Src:  filepath.Join(themeDir, file.Name()),
			Dst:  filepath.Join(themeDstDir, file.Name()),
			Perm: 0644,
		})
	}

	// Rubberband binary for Windows and Darwin
	rbDstDir := filepath.Join(path, thirdPartyDir, "rubberband")
	switch runtime.GOOS {
	case "windows":
		files = append(files, FileToExtract{
			Src:  rbBinForWindows,
			Dst:  filepath.Join(rbDstDir, "rubberband.exe"),
			Perm: 0755,
		}, FileToExtract{
			Src:  rb3BinForWindows,
			Dst:  filepath.Join(rbDstDir, "rubberband-r3.exe"),
			Perm: 0755,
		}, FileToExtract{
			Src:  rbDllForWindows,
			Dst:  filepath.Join(rbDstDir, "sndfile.dll"),
			Perm: 0644,
		})
		// Set the Rubberband command for Windows
		audio.RubberbandCommand = filepath.Join(rbDstDir, "rubberband-r3.exe")

	case "darwin":
		files = append(files, FileToExtract{
			Src:  rbBinForDarwin,
			Dst:  filepath.Join(rbDstDir, "rubberband"),
			Perm: 0755,
		}, FileToExtract{
			Src:  rb3BinForDarwin,
			Dst:  filepath.Join(rbDstDir, "rubberband-r3"),
			Perm: 0755,
		})
		// Set the Rubberband command for Darwin
		audio.RubberbandCommand = filepath.Join(rbDstDir, "rubberband-r3")
	}

	// Soundpacks
	soundpackDstDir := filepath.Join(path, soundpackDir)
	soundpackEntries, err := Assets.ReadDir(soundpackDir)
	if err != nil {
		return fmt.Errorf("failed to read soundpacks directory: %w", err)
	}
	for _, entry := range soundpackEntries {
		if entry.IsDir() {
			subDir := filepath.Join(soundpackDir, entry.Name())
			subFiles, err := Assets.ReadDir(subDir)
			if err != nil {
				return fmt.Errorf("failed to read soundpack subdir %s: %w", subDir, err)
			}
			for _, f := range subFiles {
				files = append(files, FileToExtract{
					Src:  filepath.Join(subDir, f.Name()),
					Dst:  filepath.Join(soundpackDstDir, entry.Name(), f.Name()),
					Perm: 0644,
				})
			}
		}
	}

	// Extract all files
	for _, file := range files {
		if err := extractFile(file); err != nil {
			return err
		}
	}

	return nil
}

// extractFile extracts a file from the embedded assets to the given destination path, with the given permissions
func extractFile(file FileToExtract) error {
	// Check if the file exists
	if _, err := os.Stat(file.Dst); err == nil {
		return nil
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(file.Dst), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := Assets.ReadFile(file.Src)
	if err != nil {
		return fmt.Errorf("failed to read embedded file: %w", err)
	}
	if err := os.WriteFile(file.Dst, data, file.Perm); err != nil {
		return fmt.Errorf("failed to write file to %s: %w", file.Dst, err)
	}
	return nil
}
