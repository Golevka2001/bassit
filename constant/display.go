package constant

import (
	"github.com/gdamore/tcell/v2"
	"github.com/go-music-theory/music-theory/note"
)

var (
	SkipWelcome     = false
	WelcomeDuration = 1500 // in ms
)

var (
	// DisplayedFretNum defines the number of frets to be displayed
	DisplayedFretNum = 12

	// AdjSymbolType defines how to display accidental notes
	AdjSymbolType = note.Sharp // or note.Flat
)

/* Spacing properties */
var (
	TitleViewHeight              = 1
	BassViewMarginTop            = 3
	FretboardMarginLeft          = 5
	FretboardMarginRight         = 5
	StringBaseNoteNameMarginLeft = 2
	PluckedStringSignMarginRight = 2
	NutWidth                     = 1
	// StringSpacing defines the spacing between strings
	StringSpacing = 2
	// StringMarginTop defines the spacing between the string and the fretboard edge
	StringMarginTop = 1
)

/* Color properties */
var (
	TitleFgColor           = tcell.ColorBlack
	TitleBgColor           = tcell.ColorWhite
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
	StringChar                  = '━'
	StringOverFretChar          = '┿'
	VibratingStringChar         = '═'
	VibratingStringOverFretChar = '╪'
	PressedFretSignChar         = '●'
	NotPluckedStringChar        = '░'
	PluckedStringSignChar       = '█'
)
