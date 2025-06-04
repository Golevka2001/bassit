package config

import (
	"fmt"

	C "github.com/Golevka2001/bassit/constant"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

var (
	// cViper is a viper instance for configuration
	cViper = viper.New()
	Config = config{}
)

type config struct {
	Tuning [C.StringCnt]string `yaml:"tuning"`
	Theme  string              `yaml:"theme"` // Theme name, not path
}

func LoadConfig(cfgFile string) error {
	if cfgFile != "" {
		cViper.SetConfigFile(cfgFile)
	} else {
		// Search default config in baseDir
		cViper.AddConfigPath(C.BaseDir())
		cViper.SetConfigName("config")
		cViper.SetConfigType("yaml")
	}

	cViper.AutomaticEnv()

	// Load config
	if err := cViper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}
	if err := cViper.Unmarshal(&Config, func(c *mapstructure.DecoderConfig) { c.TagName = "yaml" }); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Load theme
	if err := LoadTheme(Config.Theme); err != nil {
		return fmt.Errorf("failed to load theme '%s': %w", Config.Theme, err)
	}

	return nil
}
