package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-music-theory/music-theory/note"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		validate func(t *testing.T, cfg Config)
	}{
		{
			name: "default config values",
			config: Config{
				Tuning:             [StringCnt]string{"G2", "D2", "A1", "E1"},
				DisplayedFretCount: 12,
				MaxFretGapForHP:    4,
				AccidentalStyle:    "sharp",
				Theme:              "default",
				Soundpack:          "default",
			},
			validate: func(t *testing.T, cfg Config) {
				assert.Equal(t, [StringCnt]string{"G2", "D2", "A1", "E1"}, cfg.Tuning)
				assert.Equal(t, 12, cfg.DisplayedFretCount)
				assert.Equal(t, 4, cfg.MaxFretGapForHP)
				assert.Equal(t, "sharp", cfg.AccidentalStyle)
				assert.Equal(t, "default", cfg.Theme)
				assert.Equal(t, "default", cfg.Soundpack)
			},
		},
		{
			name: "alternative config values",
			config: Config{
				Tuning:             [StringCnt]string{"A2", "E2", "B1", "F#1"},
				DisplayedFretCount: 15,
				MaxFretGapForHP:    6,
				AccidentalStyle:    "flat",
				Theme:              "dark",
				Soundpack:          "custom",
			},
			validate: func(t *testing.T, cfg Config) {
				assert.Equal(t, [StringCnt]string{"A2", "E2", "B1", "F#1"}, cfg.Tuning)
				assert.Equal(t, 15, cfg.DisplayedFretCount)
				assert.Equal(t, 6, cfg.MaxFretGapForHP)
				assert.Equal(t, "flat", cfg.AccidentalStyle)
				assert.Equal(t, "dark", cfg.Theme)
				assert.Equal(t, "custom", cfg.Soundpack)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t, tt.config)
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T)
	}{
		{
			name: "default config structure",
			validate: func(t *testing.T) {
				assert.Equal(t, [StringCnt]string{"G2", "D2", "A1", "E1"}, defaultCfg.Tuning)
				assert.Equal(t, 12, defaultCfg.DisplayedFretCount)
				assert.Equal(t, 4, defaultCfg.MaxFretGapForHP)
				assert.Equal(t, "sharp", defaultCfg.AccidentalStyle)
				assert.Equal(t, "default", defaultCfg.Theme)
				assert.Equal(t, "default", defaultCfg.Soundpack)
			},
		},
		{
			name: "default tuning is valid",
			validate: func(t *testing.T) {
				for _, noteName := range defaultCfg.Tuning {
					n := note.Named(noteName)
					assert.NotEqual(t, note.Nil, n.Class)
					assert.GreaterOrEqual(t, n.Octave, note.Octave(0))
				}
			},
		},
		{
			name: "default displayed fret count is valid",
			validate: func(t *testing.T) {
				assert.Greater(t, defaultCfg.DisplayedFretCount, 0)
				assert.LessOrEqual(t, defaultCfg.DisplayedFretCount, MaxFretCnt)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t)
		})
	}
}

func TestValidateConfig(t *testing.T) {
	originalTuning := Tuning
	originalDisplayedFretCount := DisplayedFretCount
	originalAccidentalStyle := AccidentalStyle
	originalThemeName := ThemeName
	originalSoundpackName := SoundpackName
	originalBaseDir := BaseDir
	defer func() {
		Tuning = originalTuning
		DisplayedFretCount = originalDisplayedFretCount
		AccidentalStyle = originalAccidentalStyle
		ThemeName = originalThemeName
		SoundpackName = originalSoundpackName
		BaseDir = originalBaseDir
	}()

	tests := []struct {
		name     string
		config   Config
		strict   bool
		setup    func() string
		expected error
		validate func(t *testing.T, cfg *Config, tempDir string)
	}{
		{
			name: "valid config - strict mode",
			config: Config{
				Tuning:             [StringCnt]string{"G2", "D2", "A1", "E1"},
				DisplayedFretCount: 12,
				MaxFretGapForHP:    4,
				AccidentalStyle:    "sharp",
				Theme:              "test_theme",
				Soundpack:          "test_soundpack",
			},
			strict: true,
			setup: func() string {
				tempDir, _ := os.MkdirTemp("", "bassit_config_test_*")

				themeDir := filepath.Join(tempDir, "themes")
				soundpackDir := filepath.Join(tempDir, "soundpacks", "test_soundpack")
				os.MkdirAll(themeDir, 0755)
				os.MkdirAll(soundpackDir, 0755)

				themeFile := filepath.Join(themeDir, "test_theme.yaml")
				os.WriteFile(themeFile, []byte("test: theme"), 0644)

				BaseDir = tempDir

				return tempDir
			},
			expected: nil,
			validate: func(t *testing.T, cfg *Config, tempDir string) {
				assert.Equal(t, note.G, Tuning[0].Class)
				assert.Equal(t, note.D, Tuning[1].Class)
				assert.Equal(t, note.A, Tuning[2].Class)
				assert.Equal(t, note.E, Tuning[3].Class)
				assert.Equal(t, 12, DisplayedFretCount)
				assert.Equal(t, 4, MaxFretGapForHP)
				assert.Equal(t, note.Sharp, AccidentalStyle)
				assert.Equal(t, "test_theme", ThemeName)
				assert.Equal(t, "test_soundpack", SoundpackName)
			},
		},
		{
			name: "invalid tuning - strict mode",
			config: Config{
				Tuning: [StringCnt]string{"G2", "D2", "A1", "X1"},
				DisplayedFretCount: 12,
				MaxFretGapForHP: 4,
				AccidentalStyle: "sharp",
				Theme:           "test_theme",
				Soundpack:       "test_soundpack",
			},
			strict: true,
			setup: func() string {
				tempDir, _ := os.MkdirTemp("", "bassit_config_test_*")
				BaseDir = tempDir
				return tempDir
			},
			expected: assert.AnError,
			validate: func(t *testing.T, cfg *Config, tempDir string) {
			},
		},
		{
			name: "invalid tuning - non-strict mode",
			config: Config{
				Tuning: [StringCnt]string{"G2", "D2", "A1", "X1"},
				DisplayedFretCount: 12,
				MaxFretGapForHP: 4,
				AccidentalStyle: "sharp",
				Theme:           "test_theme",
				Soundpack:       "test_soundpack",
			},
			strict: false,
			setup: func() string {
				tempDir, _ := os.MkdirTemp("", "bassit_config_test_*")

				themeDir := filepath.Join(tempDir, "themes")
				soundpackDir := filepath.Join(tempDir, "soundpacks", "test_soundpack")
				os.MkdirAll(themeDir, 0755)
				os.MkdirAll(soundpackDir, 0755)

				themeFile := filepath.Join(themeDir, "test_theme.yaml")
				os.WriteFile(themeFile, []byte("test: theme"), 0644)

				BaseDir = tempDir
				return tempDir
			},
			expected: nil,
			validate: func(t *testing.T, cfg *Config, tempDir string) {
				assert.Equal(t, defaultCfg.Tuning, cfg.Tuning)
				assert.Equal(t, note.G, Tuning[0].Class)
				assert.Equal(t, note.D, Tuning[1].Class)
				assert.Equal(t, note.A, Tuning[2].Class)
				assert.Equal(t, note.E, Tuning[3].Class)
			},
		},
		{
			name: "invalid displayed fret count - strict mode",
			config: Config{
				Tuning:             [StringCnt]string{"G2", "D2", "A1", "E1"},
				DisplayedFretCount: 0,
				MaxFretGapForHP: 4,
				AccidentalStyle: "sharp",
				Theme:     "test_theme",
				Soundpack: "test_soundpack",
			},
			strict: true,
			setup: func() string {
				tempDir, _ := os.MkdirTemp("", "bassit_config_test_*")
				BaseDir = tempDir
				return tempDir
			},
			expected: assert.AnError,
			validate: func(t *testing.T, cfg *Config, tempDir string) {
			},
		},
		{
			name: "invalid displayed fret count - non-strict mode",
			config: Config{
				Tuning:             [StringCnt]string{"G2", "D2", "A1", "E1"},
				DisplayedFretCount: -5,
				MaxFretGapForHP: 4,
				AccidentalStyle: "sharp",
				Theme:     "test_theme",
				Soundpack: "test_soundpack",
			},
			strict: false,
			setup: func() string {
				tempDir, _ := os.MkdirTemp("", "bassit_config_test_*")

				themeDir := filepath.Join(tempDir, "themes")
				soundpackDir := filepath.Join(tempDir, "soundpacks", "test_soundpack")
				os.MkdirAll(themeDir, 0755)
				os.MkdirAll(soundpackDir, 0755)

				themeFile := filepath.Join(themeDir, "test_theme.yaml")
				os.WriteFile(themeFile, []byte("test: theme"), 0644)

				BaseDir = tempDir
				return tempDir
			},
			expected: nil,
			validate: func(t *testing.T, cfg *Config, tempDir string) {
				assert.Equal(t, defaultCfg.DisplayedFretCount, DisplayedFretCount)
			},
		},
		{
			name: "max fret gap for hp validation",
			config: Config{
				Tuning:             [StringCnt]string{"G2", "D2", "A1", "E1"},
				DisplayedFretCount: 12,
				MaxFretGapForHP:    8,
				AccidentalStyle:    "sharp",
				Theme:              "test_theme",
				Soundpack:          "test_soundpack",
			},
			strict: false,
			setup: func() string {
				tempDir, _ := os.MkdirTemp("", "bassit_config_test_*")

				themeDir := filepath.Join(tempDir, "themes")
				soundpackDir := filepath.Join(tempDir, "soundpacks", "test_soundpack")
				os.MkdirAll(themeDir, 0755)
				os.MkdirAll(soundpackDir, 0755)

				themeFile := filepath.Join(themeDir, "test_theme.yaml")
				os.WriteFile(themeFile, []byte("test: theme"), 0644)

				BaseDir = tempDir
				return tempDir
			},
			expected: nil,
			validate: func(t *testing.T, cfg *Config, tempDir string) {
				assert.Equal(t, 8, MaxFretGapForHP)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := tt.setup()
			defer os.RemoveAll(tempDir)

			err := validateConfig(&tt.config, tt.strict)

			if tt.expected != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.validate != nil {
				tt.validate(t, &tt.config, tempDir)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	originalBaseDir := BaseDir
	defer func() { BaseDir = originalBaseDir }()

	tests := []struct {
		name     string
		setup    func() (string, string)
		expected error
		validate func(t *testing.T, cfg *Config, err error, tempDir string)
	}{
		{
			name: "load valid config file",
			setup: func() (string, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_config_test_*")

				themeDir := filepath.Join(tempDir, "themes")
				soundpackDir := filepath.Join(tempDir, "soundpacks", "custom")
				os.MkdirAll(themeDir, 0755)
				os.MkdirAll(soundpackDir, 0755)

				themeFile := filepath.Join(themeDir, "custom.yaml")
				os.WriteFile(themeFile, []byte("test: theme"), 0644)

				configContent := `
tuning: ["A2", "E2", "B1", "F#1"]
displayed_fret_count: 15
max_fret_gap_for_hp: 6
accidental_style: "flat"
theme: "custom"
soundpack: "custom"
`
				configPath := filepath.Join(tempDir, "config.yaml")
				os.WriteFile(configPath, []byte(configContent), 0644)

				BaseDir = tempDir
				return configPath, tempDir
			},
			expected: nil,
			validate: func(t *testing.T, cfg *Config, err error, tempDir string) {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				assert.Equal(t, [StringCnt]string{"A2", "E2", "B1", "F#1"}, cfg.Tuning)
				assert.Equal(t, 15, cfg.DisplayedFretCount)
				assert.Equal(t, 6, cfg.MaxFretGapForHP)
				assert.Equal(t, "flat", cfg.AccidentalStyle)
				assert.Equal(t, "custom", cfg.Theme)
				assert.Equal(t, "custom", cfg.Soundpack)
			},
		},
		{
			name: "load config from default location",
			setup: func() (string, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_config_test_*")

				themeDir := filepath.Join(tempDir, "themes")
				soundpackDir := filepath.Join(tempDir, "soundpacks", "default")
				os.MkdirAll(themeDir, 0755)
				os.MkdirAll(soundpackDir, 0755)

				themeFile := filepath.Join(themeDir, "default.yaml")
				os.WriteFile(themeFile, []byte("test: theme"), 0644)

				configContent := `
tuning: ["G2", "D2", "A1", "E1"]
displayed_fret_count: 12
max_fret_gap_for_hp: 4
accidental_style: "sharp"
theme: "default"
soundpack: "default"
`
				configPath := filepath.Join(tempDir, "config.yaml")
				os.WriteFile(configPath, []byte(configContent), 0644)

				BaseDir = tempDir
				return "", tempDir
			},
			expected: nil,
			validate: func(t *testing.T, cfg *Config, err error, tempDir string) {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				assert.Equal(t, defaultCfg.Tuning, cfg.Tuning)
				assert.Equal(t, 12, cfg.DisplayedFretCount)
				assert.Equal(t, 4, cfg.MaxFretGapForHP)
				assert.Equal(t, "sharp", cfg.AccidentalStyle)
				assert.Equal(t, "default", cfg.Theme)
				assert.Equal(t, "default", cfg.Soundpack)
			},
		},
		{
			name: "config file not found",
			setup: func() (string, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_config_test_*")
				BaseDir = tempDir
				return filepath.Join(tempDir, "nonexistent.yaml"), tempDir
			},
			expected: assert.AnError,
			validate: func(t *testing.T, cfg *Config, err error, tempDir string) {
				assert.Error(t, err)
				assert.Nil(t, cfg)
				assert.Contains(t, err.Error(), "error reading config file")
			},
		},
		{
			name: "invalid YAML syntax",
			setup: func() (string, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_config_test_*")

				configContent := `
tuning: ["G2", "D2", "A1", "E1"
displayed_fret_count: 12
invalid yaml syntax here
`
				configPath := filepath.Join(tempDir, "config.yaml")
				os.WriteFile(configPath, []byte(configContent), 0644)

				BaseDir = tempDir
				return configPath, tempDir
			},
			expected: assert.AnError,
			validate: func(t *testing.T, cfg *Config, err error, tempDir string) {
				assert.Error(t, err)
				assert.Nil(t, cfg)
				assert.Contains(t, err.Error(), "error reading config file")
			},
		},
		{
			name: "config with missing theme file",
			setup: func() (string, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_config_test_*")

				soundpackDir := filepath.Join(tempDir, "soundpacks", "default")
				os.MkdirAll(soundpackDir, 0755)

				configContent := `
tuning: ["G2", "D2", "A1", "E1"]
displayed_fret_count: 12
max_fret_gap_for_hp: 4
accidental_style: "sharp"
theme: "nonexistent"
soundpack: "default"
`
				configPath := filepath.Join(tempDir, "config.yaml")
				os.WriteFile(configPath, []byte(configContent), 0644)

				BaseDir = tempDir
				return configPath, tempDir
			},
			expected: assert.AnError,
			validate: func(t *testing.T, cfg *Config, err error, tempDir string) {
				assert.Error(t, err)
				assert.Nil(t, cfg)
				assert.Contains(t, err.Error(), "theme file not found")
			},
		},
		{
			name: "config with missing soundpack directory",
			setup: func() (string, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_config_test_*")

				themeDir := filepath.Join(tempDir, "themes")
				os.MkdirAll(themeDir, 0755)

				themeFile := filepath.Join(themeDir, "default.yaml")
				os.WriteFile(themeFile, []byte("test: theme"), 0644)

				configContent := `
tuning: ["G2", "D2", "A1", "E1"]
displayed_fret_count: 12
max_fret_gap_for_hp: 4
accidental_style: "sharp"
theme: "default"
soundpack: "nonexistent"
`
				configPath := filepath.Join(tempDir, "config.yaml")
				os.WriteFile(configPath, []byte(configContent), 0644)

				BaseDir = tempDir
				return configPath, tempDir
			},
			expected: assert.AnError,
			validate: func(t *testing.T, cfg *Config, err error, tempDir string) {
				assert.Error(t, err)
				assert.Nil(t, cfg)
				assert.Contains(t, err.Error(), "soundpack folder not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath, tempDir := tt.setup()
			defer os.RemoveAll(tempDir)

			cfg, err := LoadConfig(configPath)

			if tt.validate != nil {
				tt.validate(t, cfg, err, tempDir)
			}
		})
	}
}
