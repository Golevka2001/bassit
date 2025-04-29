// bassit/internal/bass/bass.go
package bass

import (
	"time"

	"github.com/Golevka2001/bassit/internal/config"
	"github.com/Golevka2001/bassit/internal/util"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/gdamore/tcell/v2"
)

// BassString represents a single string on the bass
type BassString struct {
	Name      string
	BaseNote  float64
	KeyToFret map[rune]int
	Pressed   int // 0: open string, negative: invalid
}

// BassSimulator is the main bass simulator struct
type BassSimulator struct {
	Screen      tcell.Screen
	Strings     [config.StringNum]BassString
	SampleRate  beep.SampleRate
	CurrentNote *beep.Streamer
}

// NewBassSimulator creates a new instance of the bass simulator
func NewBassSimulator() (*BassSimulator, error) {
	// Initialize the terminal screen
	screen, err := tcell.NewScreen()
	if err != nil {
		util.LogError("Failed to create screen", err)
		return nil, err
	}
	if err := screen.Init(); err != nil {
		util.LogError("Failed to initialize screen", err)
		return nil, err
	}

	// Initialize speaker system for audio
	if err := speaker.Init(config.SampleRate, config.SampleRate.N(time.Second/10)); err != nil {
		screen.Fini()
		util.LogError("Failed to initialize speaker", err)
		return nil, err
	}

	// Create the bass bs instance
	bs := &BassSimulator{
		Screen:     screen,
		SampleRate: config.SampleRate,
	}

	// Define the strings and key mappings
	bs.Strings[0] = BassString{
		Name:     "G",
		BaseNote: config.GStringFreq,
		KeyToFret: map[rune]int{
			'1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, '0': 10, '-': 11,
		},
		Pressed: 0,
	}

	bs.Strings[1] = BassString{
		Name:     "D",
		BaseNote: config.DStringFreq,
		KeyToFret: map[rune]int{
			'q': 1, 'w': 2, 'e': 3, 'r': 4, 't': 5, 'y': 6, 'u': 7, 'i': 8, 'o': 9, 'p': 10,
		},
		Pressed: 0,
	}

	bs.Strings[2] = BassString{
		Name:     "A",
		BaseNote: config.AStringFreq,
		KeyToFret: map[rune]int{
			'a': 1, 's': 2, 'd': 3, 'f': 4, 'g': 5, 'h': 6, 'j': 7, 'k': 8, 'l': 9,
		},
		Pressed: 0,
	}

	bs.Strings[3] = BassString{
		Name:     "E",
		BaseNote: config.EStringFreq,
		KeyToFret: map[rune]int{
			'z': 1, 'x': 2, 'c': 3, 'v': 4, 'b': 5, 'n': 6, 'm': 7, ',': 8,
		},
		Pressed: 0,
	}

	return bs, nil
}

// GetFretFromKey finds which string and fret corresponds to a keyboard key
func (bs *BassSimulator) GetFretFromKey(key rune) (stringIdx, fret int) {
	for i, bassString := range bs.Strings {
		if fret, ok := bassString.KeyToFret[key]; ok {
			return i, fret
		}
	}
	return -1, -1
}

// PressString marks a string as being pressed on a specific fret
func (bs *BassSimulator) PressString(stringIdx, fret int) {
	if stringIdx >= 0 && stringIdx < len(bs.Strings) {
		bs.Strings[stringIdx].Pressed = fret
		// Play the corresponding note
		bs.PlayNote(stringIdx, fret)
		// Refresh display
		bs.Render()
	}
}

// PressString marks a string as being pressed on a specific fret
func (bs *BassSimulator) ReleaseString(stringIdx, fret int) {
	if stringIdx >= 0 && stringIdx < len(bs.Strings) {
		bs.Strings[stringIdx].Pressed = fret
		// Refresh display
		bs.Render()
	}
}

// Cleanup properly closes the simulator
func (bs *BassSimulator) Cleanup() {
	if bs.Screen != nil {
		bs.Screen.Fini()
	}
}

// Render refreshes the display
func (bs *BassSimulator) Render() {
	bs.Screen.Clear()
	RenderTitle(bs)
	RenderBassFretboard(bs)
}
