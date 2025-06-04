package constant

import (
	"github.com/go-music-theory/music-theory/note"
)

var (
	// AdjSymbolType defines how to display accidental notes
	AdjSymbolType = note.Sharp // or note.Flat
)

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
