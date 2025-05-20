package cmd

import (
    "fmt"
    "os"

    "bassit/app"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "bassit",
    Short: "Bassit is bass in terminal",
    Long:  Help,
    Run: func(cmd *cobra.Command, args []string) {
        app.Run()
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func init() {
    rootCmd.PersistentFlags().BoolP("help", "h", false, "Show this help message")
}
