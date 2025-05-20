package constant

import (
	"github.com/gdamore/tcell/v2"
	"github.com/go-music-theory/music-theory/note"
)

var (
	// DisplayedFretNum defines the number of frets to be displayed
	DisplayedFretNum = 12

	// AdjSymbolType defines how to display accidental notes
	AdjSymbolType = note.Sharp // or note.Flat
)

/* Spacing properties */
var (
	TitleMarginTop               = 1
	FretboardMarginLeft          = 5
	FretboardMarginRight         = 6
	FretboardMarginTop           = 2
	StringBaseNoteNameMarginLeft = 2
	PluckedStringSignMarginRight = 3
	NutWidth                     = 1
	// StringSpacing defines the spacing between strings
	StringSpacing = 2
	// StringMarginTop defines the spacing between the string and the fretboard edge
	StringMarginTop = 1
)

/* Color properties */
var (
	TitleFgColor           = tcell.ColorWhite
	TitleBgColor           = tcell.ColorNone
	FretboardBorderColor   = tcell.ColorWhite
	FretWireColor          = tcell.ColorSilver
	PosMarkerColor         = tcell.ColorWhiteSmoke
	FretboardBgColor       = tcell.ColorNone
	NutBorderColor         = tcell.ColorWhite
	NutBgColor             = tcell.ColorNone
	BaseNoteNameFgColor    = tcell.ColorGreen
	BaseNoteNameBgColor    = tcell.ColorNone
	StringColor            = tcell.ColorSilver
	PressedFretSignColor   = tcell.ColorRed
	PluckedStringSignColor = tcell.ColorRed
)

/* Characters for display */
var (
	FretboardVBorderChar        = '│'
	FretboardHBorderChar        = '─'
	FretboardULCornerChar       = '┌'
	FretboardLLCornerChar       = '└'
	FretWireChar                = '│'
	FretWireUpperChar           = '┬'
	FretWireLowerChar           = '┴'
	PosMarkerChar               = '◍'
	StringChar                  = '═'
	StringOverFretChar          = '╪'
	VibratingStringChar         = '━'
	VibratingStringOverFretChar = '┿'
	PressedFretSignChar         = '●'
	PluckedStringSignChar       = '█'
)
