package cmd

import (
	"path/filepath"

	"github.com/Golevka2001/bassit/assets"
	"github.com/Golevka2001/bassit/audio"
	"github.com/Golevka2001/bassit/bass"
	"github.com/Golevka2001/bassit/config"
	"github.com/Golevka2001/bassit/ui"
	"github.com/Golevka2001/bassit/ui/common"

	"github.com/spf13/cobra"
)

// Flags
var (
	baseDir   string // default: $HOME/.config/bassit
	cfgPath   string // default: $HOME/.config/bassit/config.yaml
	themeName string // default: specified in `cfgPath`
	skipCheck bool   // default: false
)

var (
	cfg   *config.Config
	theme *config.Theme
)

var (
	rootCmd = &cobra.Command{
		Use:   "bassit",
		Short: "𝄢 bassit - bass in terminal",
		Long: `bassit is a terminal-based bass guitar simulator written in Go.
It allows you to play bass lines using your keyboard.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Create an audio manager
			am, err := audio.New()
			cobra.CheckErr(err)

			// Create a bass model with the specified tuning
			bm, err := bass.New(cfg.Tuning)
			cobra.CheckErr(err)

			ctx := common.UIContext{
				Audio:     am,
				Bass:      bm,
				Config:    cfg,
				Theme:     theme,
				SkipCheck: skipCheck,
			}

			// Create a bubbletea program and run it
			p := ui.NewProgram(&ctx)
			_, err = p.Run()
			cobra.CheckErr(err)
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initApp)

	rootCmd.PersistentFlags().StringVarP(&baseDir, "base-dir", "d", "", "set base directory to store resources")

	rootCmd.Flags().StringVarP(&cfgPath, "config", "c", "", "specify configuration file")
	rootCmd.Flags().StringVarP(&themeName, "theme", "t", "", "specify theme, overrides the one in config file")
	rootCmd.Flags().BoolVar(&skipCheck, "skip-check", false, "skip checking pre-requisites (not recommended)")

	// rootCmd.AddCommand(listCmd)
}

func initApp() {
	// `baseDir`
	if baseDir != "" {
		if !filepath.IsAbs(baseDir) {
			baseDir, _ = filepath.Abs(baseDir)
		}
	}
	if err := config.SetBaseDir(baseDir); err != nil {
		cobra.CheckErr(err)
	}

	// Extract embedded resources
	if err := assets.ExtractTo(config.BaseDir()); err != nil {
		cobra.CheckErr(err)
	}

	// `cfgFile`
	if cfgPath != "" {
		if !filepath.IsAbs(cfgPath) {
			cfgPath, _ = filepath.Abs(cfgPath)
		}
	}
	var err error
	cfg, err = config.LoadConfig(cfgPath)
	if err != nil {
		cobra.CheckErr(err)
	}

	// `theme`
	if themeName == "" {
		themeName = cfg.Theme
	}
	theme, err = config.LoadTheme(themeName)
	if err != nil {
		cobra.CheckErr(err)
	}
}
