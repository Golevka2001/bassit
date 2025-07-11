package fretboard

import "github.com/Golevka2001/bassit/bass"

type PressFretMsg bass.FretboardPosition

type ReleaseFretMsg bass.FretboardPosition

type PluckStringMsg struct {
	bass.FretboardPosition
	Type bass.PluckType
}

type RestorePluckedStringMsg bass.FretboardPosition
