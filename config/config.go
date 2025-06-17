package config

import (
	"fmt"
	"strings"

	"github.com/go-music-theory/music-theory/note"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

var (
	// AccidentalStyle defines how to display accidental notes
	AccidentalStyle note.AdjSymbol

	// DisplayedFretCount defines the number of frets to be displayed
	DisplayedFretCount int
)

type Config struct {
	// Tuning defines the tuning for the 4-string bass guitar
	// Octave number is required
	// From the highest string to the lowest string
	Tuning [StringCnt]string `yaml:"tuning"`

	// Theme defines the name of the theme to be used
	// The theme will be loaded from the `${BASSIT_BASE_DIR}/themes/` directory
	// It should match the yaml file name in that directory, excluding the `.yaml` extension
	Theme string `yaml:"theme"`

	// DisplayedFretCount defines the number of frets to be displayed
	DisplayedFretCount int `yaml:"displayed_fret_count"`

	// AccidentalStyle defines how to display accidental notes
	// It can be "sharp" or "flat"
	AccidentalStyle string `yaml:"accidental_style"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	if path != "" {
		v.SetConfigFile(path)
	} else {
		// Search default config in `${BASSIT_BASE_DIR}``
		v.AddConfigPath(BaseDir())
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	}

	v.AutomaticEnv()

	// Load config
	var cfg Config
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	if err := v.Unmarshal(&cfg, func(c *mapstructure.DecoderConfig) { c.TagName = "yaml" }); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set `AccidentalStyle`
	switch strings.ToLower(cfg.AccidentalStyle) {
	case "sharp":
		AccidentalStyle = note.Sharp
	case "flat":
		AccidentalStyle = note.Flat
	default:
		AccidentalStyle = note.Sharp
	}

	// Set `DisplayedFretCount`
	if cfg.DisplayedFretCount <= 0 || cfg.DisplayedFretCount > MaxFretCnt {
		DisplayedFretCount = DefaultDisplayedFretCount
	} else {
		DisplayedFretCount = cfg.DisplayedFretCount
	}

	return &cfg, nil
}
