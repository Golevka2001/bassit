package utils

import (
	"fmt"

	"github.com/Golevka2001/bassit/config"

	"github.com/go-music-theory/music-theory/note"
)

// GetNoteStepFrom returns the note that is `inc` steps away from the given note
// `music-theory` library only provides a method for `Class`
// Parameters:
//   - from: the starting note
//   - inc: the number of steps to move (can be negative)
func GetNoteStepFrom(from note.Note, inc int) note.Note {
	if inc == 0 {
		return from
	}

	class, octave := from.Class.Step(inc)

	target := from
	target.Class = class
	target.Octave += octave // NOTE: `octave` is how may octaves to shift

	// TODO: Do other properties need to be updated?

	return target
}

// GetStepBetween returns the number of steps between two notes
// Returns:
//   - positive if `to` is higher than `from`
//   - negative if `to` is lower than `from`
func GetStepBetween(from, to note.Note) int {
	OctaveDiff := int(to.Octave - from.Octave)
	StepDiff := int(to.Class - from.Class)
	return StepDiff + OctaveDiff*12
}

// GetNoteNameWithOctave returns the note name with octave (e.g., "A2", "C#3")
func GetNoteNameWithOctave(note note.Note) string {
	return fmt.Sprintf("%s%d", note.Class.String(config.AccidentalStyle), note.Octave)
}
