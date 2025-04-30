// bassit/internal/bass/bass.go
package bass

import (
    "maps"

    "github.com/Golevka2001/bassit/internal/config"
    C "github.com/Golevka2001/bassit/internal/constant"
    "github.com/Golevka2001/bassit/internal/util"

    "github.com/faiface/beep"
    "github.com/gdamore/tcell/v2"
)

// BassString represents a single string on the bass
type BassString struct {
    BaseNoteName string
    //BaseNoteFreq float64
    KeyToFret   map[rune]int
    KeyToPluck  map[rune]int
    PressedFret int // 0: open string, 1-n: fret number
    IsPlucked   bool
}

// BassSimulator is the main bass simulator struct
type BassSimulator struct {
    Screen      tcell.Screen
    Strings     [C.StringNum]BassString
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

    //// Initialize speaker system for audio
    //if err := speaker.Init(config.SampleRate, config.SampleRate.N(time.Second/10)); err != nil {
    //    screen.Fini()
    //    util.LogError("Failed to initialize speaker", err)
    //    return nil, err
    //}

    // Create the bass bs instance
    bs := &BassSimulator{
        Screen:     screen,
        SampleRate: C.SampleRate,
    }

    bs.initializeStrings()

    return bs, nil
}

func (bs *BassSimulator) initializeStrings() {
    for stringIdx := range bs.Strings {
        curString := &bs.Strings[stringIdx]
        curString.BaseNoteName = C.StandardTuning[stringIdx]
        // TODO: bassString.BaseNoteFreq =
        curString.KeyToFret = C.Key2Fret[stringIdx]
        curString.KeyToPluck = C.Key2Pluck[stringIdx]
        if config.ShiftedKeyEnabled {
            maps.Copy(curString.KeyToFret, C.ShiftedKey2Fret[stringIdx])
            maps.Copy(curString.KeyToPluck, C.ShiftedKey2Pluck[stringIdx])
        }
        curString.PressedFret = 0
        curString.IsPlucked = false
    }
}

// GetPressedFretFromKey finds which string and fret corresponds to a keyboard key
func (bs *BassSimulator) GetPressedFretFromKey(key rune) (stringIdx, fret int) {
    for stringIdx, bassString := range bs.Strings {
        if fret, ok := bassString.KeyToFret[key]; ok {
            return stringIdx, fret
        }
    }
    return -1, -1
}

func (bs *BassSimulator) GetPluckedStringFromKey(key rune) (stringIdx int, pluck int) {
    for stringIdx, bassString := range bs.Strings {
        if pluck, ok := bassString.KeyToPluck[key]; ok {
            return stringIdx, pluck
        }
    }
    return -1, -1
}

func (bs *BassSimulator) PluckString(stringIdx, pluck int) {
    if stringIdx >= 0 && stringIdx < len(bs.Strings) {
        bs.Strings[stringIdx].IsPlucked = true
        // Refresh display
        bs.Render()
    }
}

// PressString marks a string as being pressed on a specific fret
func (bs *BassSimulator) PressString(stringIdx, fret int) {
    if stringIdx >= 0 && stringIdx < len(bs.Strings) {
        bs.Strings[stringIdx].PressedFret = fret
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
