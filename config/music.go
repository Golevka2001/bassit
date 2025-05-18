package config

import (
	"github.com/go-music-theory/music-theory/note"
)

var (
	// AdjSymbolType defines how to display accidental notes
	AdjSymbolType = note.Sharp // or note.Flat

	// StandardPitch defines the standard pitch for tuning
	StandardPitch = 440.0 // A4 in Hz
)
