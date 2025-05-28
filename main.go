package main

import (
	"os"
	"runtime"

	"bassit/cmd"
	C "bassit/constant"
)

func main() {
	// Detect OS
	C.OS = runtime.GOOS

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
