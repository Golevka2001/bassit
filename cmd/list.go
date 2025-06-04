package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	C "github.com/Golevka2001/bassit/constant"

	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:       "list",
		Aliases:   []string{"ls"},
		Short:     "List available resources",
		Long:      ``,
		ValidArgs: []string{"themes", "samples"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "themes":
				listThemesCmd.Run(cmd, args)
			case "samples":
				listSamplesCmd.Run(cmd, args)
			}
		},
	}

	listThemesCmd = &cobra.Command{
		Use:   "themes",
		Short: "List available themes",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			listThemes()
		},
	}

	listSamplesCmd = &cobra.Command{
		Use:   "samples",
		Short: "List available samples",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			listSamples()
		},
	}
)

func init() {
	listCmd.AddCommand(listThemesCmd)
	listCmd.AddCommand(listSamplesCmd)
}

// listThemes prints all available themes in the themes directory
func listThemes() {
	files, err := os.ReadDir(C.ThemeDir())
	cobra.CheckErr(err)

	fmt.Println("Available themes:")
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml") {
			theme := strings.TrimSuffix(file.Name(), ".yaml")
			theme = strings.TrimSuffix(theme, ".yml")
			fmt.Printf("- %s\n", theme)
		}
	}
	if len(files) == 0 {
		fmt.Println("No themes found")
	}
}

// listSamples prints all available samples in the samples directory
func listSamples() {
	// pluck
	pluckDir := filepath.Join(C.SampleDir(), "pluck")
	pluckFiles, err := os.ReadDir(pluckDir)
	cobra.CheckErr(err)
	fmt.Println("Available pluck samples:")
	for _, file := range pluckFiles {
		fmt.Printf("- %s\n", file.Name())
	}
	if len(pluckFiles) == 0 {
		fmt.Println("No pluck samples found")
	}

	// slap
	fmt.Println("Available slap samples:")
	slapDir := filepath.Join(C.SampleDir(), "slap")
	slapFiles, err := os.ReadDir(slapDir)
	cobra.CheckErr(err)
	for _, file := range slapFiles {
		fmt.Printf("- %s\n", file.Name())
	}
	if len(pluckFiles) == 0 {
		fmt.Println("No slap samples found")
	}

	// mute
	fmt.Println("Available mute samples:")
	muteDir := filepath.Join(C.SampleDir(), "mute")
	muteFiles, err := os.ReadDir(muteDir)
	cobra.CheckErr(err)
	for _, file := range muteFiles {
		fmt.Printf("- %s\n", file.Name())
	}
	if len(pluckFiles) == 0 {
		fmt.Println("No mute samples found")
	}
}
