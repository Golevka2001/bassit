// Package input bassit/internal/input/keyboard.go
package input

import (
    "github.com/Golevka2001/bassit/internal/bass"

    "github.com/gdamore/tcell/v2"
)

// KeyboardHandler manages keyboard input for the bass simulator
type KeyboardHandler struct {
    simulator *bass.BassSimulator
}

// NewKeyboardHandler creates a new keyboard handler
func NewKeyboardHandler(bs *bass.BassSimulator) (*KeyboardHandler, error) {
    return &KeyboardHandler{
        simulator: bs,
    }, nil
}

// StartListening starts the keyboard event loop
func (k *KeyboardHandler) StartListening() {
    for {
        // Poll for events
        ev := k.simulator.Screen.PollEvent()

        // Handle events
        switch ev := ev.(type) {
        case *tcell.EventKey:
            if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
                return
            }

            if ev.Key() == tcell.KeyRune || ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
                keyRune := ev.Rune()
                // NOTE: `backspace` is `Delete` on MacOS
                if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
                    keyRune = '\b'
                }

                // Check key press for plucking
                pluckedStringIdx, pluck := k.simulator.GetPluckedStringFromKey(keyRune)
                if pluckedStringIdx >= 0 {
                    // Pluck the string
                    k.simulator.PluckString(pluckedStringIdx, pluck)
                }

                // Check key press for bass string interaction
                pressedStringIdx, fret := k.simulator.GetPressedFretFromKey(keyRune)
                if pressedStringIdx >= 0 {
                    // Press the string at the determined fret
                    k.simulator.PressString(pressedStringIdx, fret)
                }
            }
        case *tcell.EventResize:
            k.simulator.Screen.Sync()
            k.simulator.Render()
        }
    }
}
