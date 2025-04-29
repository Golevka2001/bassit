// bassit/internal/input/keyboard.go
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

			// Check key press for bass string interaction
			if ev.Key() == tcell.KeyRune {
				keyRune := ev.Rune()
				stringIdx, fret := k.simulator.GetFretFromKey(keyRune)

				if stringIdx >= 0 {
					// Press the string at the determined fret
					k.simulator.PressString(stringIdx, fret)
				}
			}
		case *tcell.EventResize:
			k.simulator.Screen.Sync()
			k.simulator.Render()
		}
	}
}
