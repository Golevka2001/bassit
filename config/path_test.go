package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetBaseDir(t *testing.T) {
	originalBaseDir := BaseDir
	defer func() { BaseDir = originalBaseDir }()

	tests := []struct {
		name     string
		input    string
		setup    func() string
		expected error
		validate func(t *testing.T, input string)
	}{
		{
			name:  "empty path - should use home directory",
			input: "",
			setup: func() string {
				return ""
			},
			expected: nil,
			validate: func(t *testing.T, input string) {
				homeDir, _ := os.UserHomeDir()
				expectedPath := filepath.Join(homeDir, ".config", "bassit")
				assert.Equal(t, expectedPath, BaseDir)

				assert.DirExists(t, BaseDir)
			},
		},
		{
			name:  "valid absolute path",
			input: "", setup: func() string {
				tempDir, _ := os.MkdirTemp("", "bassit_test_*")
				return tempDir
			},
			expected: nil,
			validate: func(t *testing.T, input string) {
				assert.Equal(t, input, BaseDir)
			},
		},
		{
			name:  "valid relative path",
			input: "test_config",
			setup: func() string {
				cwd, _ := os.Getwd()
				testDir := filepath.Join(cwd, "test_config")
				os.MkdirAll(testDir, 0755)
				return testDir
			},
			expected: nil,
			validate: func(t *testing.T, input string) {
				expectedPath, _ := filepath.Abs(input)
				assert.Equal(t, expectedPath, BaseDir)
			},
		},
		{
			name:  "nonexistent path",
			input: "/nonexistent/path/that/does/not/exist",
			setup: func() string {
				return ""
			},
			expected: os.ErrNotExist,
			validate: func(t *testing.T, input string) {
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := tt.setup()
			if tempDir != "" {
				defer os.RemoveAll(tempDir)
			}

			input := tt.input
			if tt.input == "" && tempDir != "" {
				input = tempDir
			}

			actual := SetBaseDir(input)

			if tt.expected != nil {
				assert.Error(t, actual)
				if tt.expected == os.ErrNotExist {
					assert.Contains(t, actual.Error(), "path does not exist")
				}
			} else {
				assert.NoError(t, actual)
			}

			if tt.validate != nil {
				tt.validate(t, input)
			}

			if tt.name == "empty path - should use home directory" && actual == nil {
				homeDir, _ := os.UserHomeDir()
				configDir := filepath.Join(homeDir, ".config", "bassit")
				os.RemoveAll(configDir)
			}
		})
	}
}

func TestThemeDir(t *testing.T) {
	originalBaseDir := BaseDir
	defer func() { BaseDir = originalBaseDir }()

	tests := []struct {
		name     string
		setup    func()
		expected string
	}{
		{
			name: "default theme directory",
			setup: func() {
				BaseDir = "/test/base"
			},
			expected: filepath.Join("/test/base", "themes"),
		},
		{
			name: "theme directory with different base",
			setup: func() {
				BaseDir = "/another/path"
			},
			expected: filepath.Join("/another/path", "themes"),
		},
		{
			name: "theme directory with empty base",
			setup: func() {
				BaseDir = ""
			},
			expected: filepath.Join("", "themes"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			actual := ThemeDir()

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestSoundpackDir(t *testing.T) {
	originalBaseDir := BaseDir
	defer func() { BaseDir = originalBaseDir }()

	tests := []struct {
		name     string
		setup    func()
		expected string
	}{
		{
			name: "default soundpack directory",
			setup: func() {
				BaseDir = "/test/base"
			},
			expected: filepath.Join("/test/base", "soundpacks"),
		},
		{
			name: "soundpack directory with different base",
			setup: func() {
				BaseDir = "/another/path"
			},
			expected: filepath.Join("/another/path", "soundpacks"),
		},
		{
			name: "soundpack directory with empty base",
			setup: func() {
				BaseDir = ""
			},
			expected: filepath.Join("", "soundpacks"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			actual := SoundpackDir()

			assert.Equal(t, tt.expected, actual)
		})
	}
}
