// bassit/main.go
package main

import (
    "os"

    "github.com/Golevka2001/bassit/internal/bass"
    "github.com/Golevka2001/bassit/internal/input"
    "github.com/Golevka2001/bassit/internal/util"
)

func main() {
    // Initialize the bass simulator
    bs, err := bass.NewBassSimulator()
    if err != nil {
        util.LogError("Failed to initialize bass bs", err)
        os.Exit(1)
    }
    bs.Render()
    defer bs.Cleanup()

    // Initialize the keyboard handler
    keyboard, err := input.NewKeyboardHandler(bs)
    if err != nil {
        util.LogError("Failed to initialize keyboard handler", err)
        os.Exit(1)
    }

    keyboard.StartListening()
}
