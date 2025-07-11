package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-music-theory/music-theory/note"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

var (
	Tuning             [StringCnt]*note.Note
	DisplayedFretCount int
	AccidentalStyle    note.AdjSymbol
	ThemeName          string
	SoundpackName      string
)

var (
	defaultCfg = Config{
		Tuning:             [StringCnt]string{"G2", "D2", "A1", "E1"},
		DisplayedFretCount: 12,
		AccidentalStyle:    "sharp",
		Theme:              "default",
		Soundpack:          "default",
	}
)

type Config struct {
	// Tuning defines the tuning for a 4-string bass guitar
	// Octave number is required
	// Strings are ordered from string 1 (highest) to string 4 (lowest)
	Tuning [StringCnt]string `yaml:"tuning"`

	// DisplayedFretCount defines how many frets to display on the fretboard
	DisplayedFretCount int `yaml:"displayed_fret_count"`

	// AccidentalStyle determines how to display accidental notes
	// Valid values are "sharp" or "flat"
	AccidentalStyle string `yaml:"accidental_style"`

	// Theme defines the name of the UI theme to use
	// The theme file will be loaded from `${BASSIT_BASE_DIR}/themes/`
	// This should match the filename (without `.yaml`) in that directory
	Theme string `yaml:"theme"`

	// Soundpack defines the name of the sound pack to use
	// The sound pack will be loaded from `${BASSIT_BASE_DIR}/assets/soundpacks/`
	// This should match the folder name in that directory
	Soundpack string `yaml:"soundpack"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	if path != "" {
		v.SetConfigFile(path)
	} else {
		// Search default config in `${BASSIT_BASE_DIR}``
		v.AddConfigPath(BaseDir)
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	}

	v.AutomaticEnv()

	// Load config
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	cfg := defaultCfg
	if err := v.Unmarshal(&cfg, func(c *mapstructure.DecoderConfig) {
		c.TagName = "yaml"
	}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate the config
	if err := validateConfig(&cfg, true); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &cfg, nil
}

// validateConfig validates the config
func validateConfig(cfg *Config, strict bool) error {
	// `Tuning`
	if len(cfg.Tuning) != StringCnt {
		return fmt.Errorf("`tuning` must be an array of length %d", StringCnt)
	}
	invalid := false
	for i, noteName := range cfg.Tuning {
		n := note.Named(noteName)
		if n.Class == note.Nil || n.Octave < 0 {
			invalid = true
		}
		Tuning[i] = n
	}
	if invalid {
		if strict {
			return fmt.Errorf("invalid tuning: %s", cfg.Tuning)
		}
		cfg.Tuning = defaultCfg.Tuning
		for i := range cfg.Tuning {
			Tuning[i] = note.Named(cfg.Tuning[i])
		}
	}

	// `DisplayedFretCount`
	if cfg.DisplayedFretCount <= 0 || cfg.DisplayedFretCount > MaxFretCnt {
		if strict {
			return fmt.Errorf("`displayed_fret_count` must be between 1 and %d", MaxFretCnt)
		}
		DisplayedFretCount = defaultCfg.DisplayedFretCount
	} else {
		DisplayedFretCount = cfg.DisplayedFretCount
	}

	// `AccidentalStyle`
	switch strings.ToLower(cfg.AccidentalStyle) {
	case "sharp":
		AccidentalStyle = note.Sharp
	case "flat":
		AccidentalStyle = note.Flat
	default:
		if strict {
			return fmt.Errorf("`accidental_style` must be either `sharp` or `flat`")
		}
		AccidentalStyle = note.Sharp
	}

	// `Theme`
	if cfg.Theme == "" {
		if strict {
			return fmt.Errorf("`theme` is required")
		}
		ThemeName = defaultCfg.Theme
	} else {
		ThemeName = cfg.Theme
	}
	// Check if the theme file exists
	themePath := filepath.Join(ThemeDir(), ThemeName+".yaml")
	if _, err := os.Stat(themePath); err != nil {
		themePath = filepath.Join(ThemeDir(), ThemeName+".yml")
	}
	if _, err := os.Stat(themePath); err != nil {
		return fmt.Errorf("theme file not found: %s", ThemeName)
	}

	// `Soundpack`
	if cfg.Soundpack == "" {
		if strict {
			return fmt.Errorf("`soundpack` is required")
		}
		SoundpackName = defaultCfg.Soundpack
	} else {
		SoundpackName = cfg.Soundpack
	}
	// Check if the soundpack folder exists
	if _, err := os.Stat(filepath.Join(SoundpackDir(), SoundpackName)); err != nil {
		return fmt.Errorf("soundpack folder not found: %s", SoundpackName)
	}

	return nil
}
