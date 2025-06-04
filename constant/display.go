package constant

import (
	"github.com/go-music-theory/music-theory/note"
)

var (
	// AdjSymbolType defines how to display accidental notes
	AdjSymbolType = note.Sharp // or note.Flat
)

/* Spacing properties */
const (
	TitleViewHeight              = 1
	BassViewMarginTop            = 3
	FretboardMarginX             = 5
	StringBaseNoteNameMarginLeft = 2
	PluckedStringSignMarginRight = 2
	NutWidth                     = 2
	BlockInlayMarginY            = 2
	// StringSpacing defines the spacing between strings
	StringSpacing = 2
	// StringMarginY defines the spacing between the string and the fretboard edge
	StringMarginY = 1
)

// /* Color properties */
// var (
// 	TitleFgColor = tcell.ColorBlack
// 	TitleBgColor = tcell.ColorWhite

// 	FretboardBorderColor = tcell.NewHexColor(0xefc08c)
// 	FretboardBgColor     = tcell.NewHexColor(0xefc08c)

// 	NutBorderColor = tcell.NewHexColor(0xefeeeb)
// 	NutBgColor     = tcell.NewHexColor(0xefeeeb)

// 	FretWireColor = tcell.NewHexColor(0x9e9084)

// 	StringColor = tcell.NewHexColor(0xa29285)

// 	PosMarkerColor = tcell.NewHexColor(0x332e30)

// 	PressedFretSignColor   = tcell.ColorRed
// 	PluckedStringSignColor = tcell.ColorRed
// 	BaseNoteNameFgColor    = tcell.ColorGreen
// 	BaseNoteNameBgColor    = tcell.ColorNone
// )

// /* Characters for display */
// const (
// 	FretboardVBorderChar = '│'
// 	FretboardHBorderChar = ' ' // '─'

// 	NutVBorderChar  = '│'
// 	NutHBorderChar  = ' ' // '─'
// 	NutULCornerChar = ' ' // '┌'
// 	NutLLCornerChar = ' ' // '└'
// 	NutURCornerChar = ' ' // '┬'
// 	NutLRCornerChar = ' ' // '┴'

// 	FretWireChar      = '║'
// 	FretWireUpperChar = ' ' // '╥'
// 	FretWireLowerChar = ' ' // '╨'

// 	StringChar            = '━'
// 	StringOverFretChar    = '╫'
// 	StringOverBoarderChar = '┿'

// 	VibratingStringChar            = '═'
// 	VibratingStringOverFretChar    = '╬'
// 	VibratingStringOverBoarderChar = '╪'

// 	PressedFretSignChar   = '●'
// 	NotPluckedStringChar  = '░'
// 	PluckedStringSignChar = '█'
// )
