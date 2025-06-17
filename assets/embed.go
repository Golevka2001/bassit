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

	themeDir = "themes/"

	rbBinForWindows  = "3rdparty/rubberband-4.0.0-gpl-executable-windows/rubberband.exe"
	rb3BinForWindows = "3rdparty/rubberband-4.0.0-gpl-executable-windows/rubberband-r3.exe"
	rbDllForWindows  = "3rdparty/rubberband-4.0.0-gpl-executable-windows/sndfile.dll"

	rbBinForDarwin  = "3rdparty/rubberband-4.0.0-gpl-executable-macos/rubberband"
	rb3BinForDarwin = "3rdparty/rubberband-4.0.0-gpl-executable-macos/rubberband-r3"

	audioFile = "audio/bass/pluck/default/C2.wav"
)

type FileToExtract struct {
	Src  string
	Dst  string
	Perm os.FileMode
}

func ExtractTo(path string) error {
	// Configuration file
	files := []FileToExtract{
		{
			Src:  cfgFile,
			Dst:  filepath.Join(path, "config.yaml"),
			Perm: 0644,
		},
	}

	// Theme files
	themeDstDir := filepath.Join(path, "themes")
	themeFiles, err := Assets.ReadDir("themes")
	if err != nil {
		return fmt.Errorf("failed to read themes directory: %w", err)
	}
	for _, file := range themeFiles {
		files = append(files, FileToExtract{
			Src:  themeDir + file.Name(),
			Dst:  filepath.Join(themeDstDir, file.Name()),
			Perm: 0644,
		})
	}

	// Rubberband binary for Windows and Darwin
	rbDstDir := filepath.Join(path, "3rdparty/rubberband/")
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

	// Audio files
	audio.NoteSampleDir = filepath.Join(path, "audio/bass/pluck/default/")
	files = append(files, FileToExtract{
		Src:  audioFile,
		Dst:  filepath.Join(audio.NoteSampleDir, "C2.wav"),
		Perm: 0644,
	})

	for _, file := range files {
		if err := extractFile(file.Src, file.Dst, file.Perm); err != nil {
			return err
		}
	}

	return nil
}

func extractFile(src, dst string, perm os.FileMode) error {
	// Check if the file exists
	if _, err := os.Stat(dst); err == nil {
		return nil
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
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
