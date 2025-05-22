package cmd

import (
	"bassit/app"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bassit",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app.Run()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
