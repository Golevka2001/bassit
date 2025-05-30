package cmd

import (
	"path/filepath"

	"github.com/Golevka2001/bassit/app"
	A "github.com/Golevka2001/bassit/assets"
	"github.com/Golevka2001/bassit/config"
	C "github.com/Golevka2001/bassit/constant"
	"github.com/Golevka2001/bassit/view"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	baseDir string // default: $HOME/.config/bassit
	cfgFile string // default: $HOME/.config/bassit/config.yaml

	rootCmd = &cobra.Command{
		Use:   "bassit",
		Short: "𝄢 bassit - bass in terminal",
		Long: `bassit is a terminal-based bass guitar simulator written in Go.
It allows you to play bass lines using your keyboard.`,
		Run: func(cmd *cobra.Command, args []string) {
			app.Run()
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&baseDir, "base-dir", "d", "", "set base directory to store resources")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "specify configuration file")
	rootCmd.PersistentFlags().BoolVar(&view.SkipCheck, "skip-check", false, "skip checking pre-requisites (not recommended)")
}

func initConfig() {
	// `baseDir`
	if baseDir != "" {
		if !filepath.IsAbs(baseDir) {
			baseDir, _ = filepath.Abs(baseDir)
		}
	}
	err := C.SetBaseDir(baseDir)
	if err != nil {
		cobra.CheckErr(err)
	}

	// Extract embedded resources
	err = A.ExtractTo(C.BaseDir)
	if err != nil {
		cobra.CheckErr(err)
	}

	// `cfgFile`
	if cfgFile != "" {
		if !filepath.IsAbs(cfgFile) {
			cfgFile, _ = filepath.Abs(cfgFile)
		}
		viper.SetConfigFile(cfgFile)
	} else {
		// Search default config in baseDir
		viper.AddConfigPath(C.BaseDir)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	// Load config
	if err := viper.ReadInConfig(); err != nil {
		cobra.CheckErr(err)
	}
	err = config.Unmarshal()
	if err != nil {
		cobra.CheckErr(err)
	}
}
