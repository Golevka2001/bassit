package main

import (
	"github.com/Golevka2001/bassit/cmd"

	"github.com/spf13/cobra"
)

func main() {
	err := cmd.Execute()
	cobra.CheckErr(err)
}
