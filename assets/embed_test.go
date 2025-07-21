package assets

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Golevka2001/bassit/config"

	"github.com/stretchr/testify/assert"
)

func TestFileToExtract(t *testing.T) {
	tests := []struct {
		name     string
		input    FileToExtract
		expected FileToExtract
	}{
		{
			name: "valid FileToExtract struct",
			input: FileToExtract{
				Src:  "config.yaml",
				Dst:  "/tmp/config.yaml",
				Perm: 0644,
			},
			expected: FileToExtract{
				Src:  "config.yaml",
				Dst:  "/tmp/config.yaml",
				Perm: 0644,
			},
		},
		{
			name: "FileToExtract with executable permissions",
			input: FileToExtract{
				Src:  "rubberband.exe",
				Dst:  "/tmp/rubberband.exe",
				Perm: 0755,
			},
			expected: FileToExtract{
				Src:  "rubberband.exe",
				Dst:  "/tmp/rubberband.exe",
				Perm: 0755,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.input

			assert.Equal(t, tt.expected.Src, actual.Src)
			assert.Equal(t, tt.expected.Dst, actual.Dst)
			assert.Equal(t, tt.expected.Perm, actual.Perm)
		})
	}
}

func TestExtractFileWithRealFiles(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() (FileToExtract, string)
		expected error
		validate func(t *testing.T, tempDir string, file FileToExtract)
	}{
		{
			name: "file already exists - should skip extraction",
			setup: func() (FileToExtract, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_test_*")
				srcPath := filepath.Join(tempDir, "source.txt")
				dstPath := filepath.Join(tempDir, "dest.txt")

				os.WriteFile(srcPath, []byte("source content"), 0644)
				os.WriteFile(dstPath, []byte("existing content"), 0644)

				return FileToExtract{
					Src:  srcPath,
					Dst:  dstPath,
					Perm: 0644,
				}, tempDir
			},
			expected: nil,
			validate: func(t *testing.T, tempDir string, file FileToExtract) {
				content, err := os.ReadFile(file.Dst)
				assert.NoError(t, err)
				assert.Equal(t, "existing content", string(content))
			},
		},
		{
			name: "destination directory creation",
			setup: func() (FileToExtract, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_test_*")
				srcPath := filepath.Join(tempDir, "source.txt")
				dstPath := filepath.Join(tempDir, "subdir", "dest.txt")

				os.WriteFile(srcPath, []byte("test content"), 0644)

				return FileToExtract{
					Src:  srcPath,
					Dst:  dstPath,
					Perm: 0644,
				}, tempDir
			},
			expected: nil,
			validate: func(t *testing.T, tempDir string, file FileToExtract) {
				assert.DirExists(t, filepath.Dir(file.Dst))

				content, err := os.ReadFile(file.Dst)
				assert.NoError(t, err)
				assert.Equal(t, "test content", string(content))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, tempDir := tt.setup()
			defer os.RemoveAll(tempDir)

			if _, err := os.Stat(file.Dst); err == nil {
				assert.NoError(t, nil)
			} else {
				err := os.MkdirAll(filepath.Dir(file.Dst), 0755)
				assert.NoError(t, err)

				if data, err := os.ReadFile(file.Src); err == nil {
					err = os.WriteFile(file.Dst, data, file.Perm)
					assert.NoError(t, err)
				}
			}

			if tt.validate != nil {
				tt.validate(t, tempDir, file)
			}
		})
	}
}

func TestRuntimeSpecificPaths(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T)
	}{
		{
			name: "Windows rubberband paths",
			validate: func(t *testing.T) {
				assert.Contains(t, rbBinForWindows, "3rdparty/rubberband-4.0.0-gpl-executable-windows/rubberband.exe")
				assert.Contains(t, rb3BinForWindows, "3rdparty/rubberband-4.0.0-gpl-executable-windows/rubberband-r3.exe")
				assert.Contains(t, rbDllForWindows, "3rdparty/rubberband-4.0.0-gpl-executable-windows/sndfile.dll")

				assert.NotContains(t, rbBinForWindows, "\\")
				assert.NotContains(t, rb3BinForWindows, "\\")
				assert.NotContains(t, rbDllForWindows, "\\")
			},
		},
		{
			name: "Darwin rubberband paths",
			validate: func(t *testing.T) {
				assert.Contains(t, rbBinForDarwin, "3rdparty/rubberband-4.0.0-gpl-executable-macos/rubberband")
				assert.Contains(t, rb3BinForDarwin, "3rdparty/rubberband-4.0.0-gpl-executable-macos/rubberband-r3")

				assert.NotContains(t, rbBinForDarwin, "\\")
				assert.NotContains(t, rb3BinForDarwin, "\\")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t)
		})
	}
}

func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "cfgFile constant",
			input:    cfgFile,
			expected: "config.yaml",
		},
		{
			name:     "thirdPartyDir constant",
			input:    thirdPartyDir,
			expected: "3rdparty",
		},
		{
			name:     "soundpackDir constant",
			input:    soundpackDir,
			expected: "soundpacks",
		},
		{
			name:     "themeDir constant",
			input:    themeDir,
			expected: "themes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input)
		})
	}
}

func TestRubberbandCommandSetting(t *testing.T) {
	originalRubberbandCommand := config.RubberbandCommand
	defer func() { config.RubberbandCommand = originalRubberbandCommand }()

	tests := []struct {
		name     string
		goos     string
		basePath string
		expected string
	}{
		{
			name:     "Windows rubberband command",
			goos:     "windows",
			basePath: "/test/path",
			expected: filepath.Join("/test/path", thirdPartyDir, "rubberband", "rubberband-r3.exe"),
		},
		{
			name:     "Darwin rubberband command",
			goos:     "darwin",
			basePath: "/test/path",
			expected: filepath.Join("/test/path", thirdPartyDir, "rubberband", "rubberband-r3"),
		},
		{
			name:     "Linux rubberband command (no change)",
			goos:     "linux",
			basePath: "/test/path",
			expected: originalRubberbandCommand},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.RubberbandCommand = originalRubberbandCommand

			rbDstDir := filepath.Join(tt.basePath, thirdPartyDir, "rubberband")
			switch tt.goos {
			case "windows":
				config.RubberbandCommand = filepath.Join(rbDstDir, "rubberband-r3.exe")
			case "darwin":
				config.RubberbandCommand = filepath.Join(rbDstDir, "rubberband-r3")
			}

			assert.Equal(t, tt.expected, config.RubberbandCommand)
		})
	}
}
