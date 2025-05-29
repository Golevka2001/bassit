package main

import (
	"os"
	"runtime"

	"github.com/Golevka2001/bassit/cmd"
	C "github.com/Golevka2001/bassit/constant"
)

func main() {
	// Detect OS
	C.OS = runtime.GOOS

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
