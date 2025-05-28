package assets

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	C "bassit/constant"
)

//go:embed *
var Assets embed.FS

const (
	cfgFile = "config.yaml"

	rbBinForWindows  = "3rdparty/rubberband-4.0.0-gpl-executable-windows/rubberband.exe"
	rb3BinForWindows = "3rdparty/rubberband-4.0.0-gpl-executable-windows/rubberband-r3.exe"
	rbDllForWindows  = "3rdparty/rubberband-4.0.0-gpl-executable-windows/sndfile.dll"

	rbBinForDarwin  = "3rdparty/rubberband-4.0.0-gpl-executable-macos/rubberband"
	rb3BinForDarwin = "3rdparty/rubberband-4.0.0-gpl-executable-macos/rubberband-r3"
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

	// Rubberband binary for Windows and Darwin
	rbDstPrefix := filepath.Join(path, "3rdparty/rubberband/")
	switch C.OS {
	case "windows":
		files = append(files, FileToExtract{
			Src:  rbBinForWindows,
			Dst:  filepath.Join(rbDstPrefix, "rubberband.exe"),
			Perm: 0755,
		}, FileToExtract{
			Src:  rb3BinForWindows,
			Dst:  filepath.Join(rbDstPrefix, "rubberband-r3.exe"),
			Perm: 0755,
		}, FileToExtract{
			Src:  rbDllForWindows,
			Dst:  filepath.Join(rbDstPrefix, "sndfile.dll"),
			Perm: 0644,
		})

		// Set the Rubberband command for Windows
		C.RubberbandCommand = filepath.Join(rbDstPrefix, "rubberband-r3.exe")
	case "darwin":
		files = append(files, FileToExtract{
			Src:  rbBinForDarwin,
			Dst:  filepath.Join(rbDstPrefix, "rubberband"),
			Perm: 0755,
		}, FileToExtract{
			Src:  rb3BinForDarwin,
			Dst:  filepath.Join(rbDstPrefix, "rubberband-r3"),
			Perm: 0755,
		})

		// Set the Rubberband command for Darwin
		C.RubberbandCommand = filepath.Join(rbDstPrefix, "rubberband-r3")
	}

	// Audio files
	C.NoteSampleDir = filepath.Join(path, "audio/note/")
	files = append(files, FileToExtract{
		Src:  "audio/note/C2.wav",
		Dst:  filepath.Join(C.NoteSampleDir, "C2.wav"),
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
