package util

import (
	"fmt"

	"bassit/config"

	"github.com/go-music-theory/music-theory/note"
)

// GetNoteStepFrom returns the note that is `inc` steps away from the given note
// `music-theory` library only provides a method for `Class`
// Parameters:
//   - from: the starting note
//   - inc: the number of steps to move (can be negative)
func GetNoteStepFrom(from note.Note, inc int) (note.Note, error) {
	if inc == 0 {
		return from, nil
	}

	class, octave := from.Class.Step(inc)
	if class == note.Nil {
		return from, fmt.Errorf("failed to get the note step from %s with inc %d", from.Class.String(config.AdjSymbolType), inc)
	}

	target := from
	target.Class = class
	target.Octave += octave // NOTE: `octave` is how may octaves to shift

	// TODO: Do other properties need to be updated?

	return target, nil
}

// GetNoteNameWithOctave returns the note name with octave (e.g., "A2", "C#3")
func GetNoteNameWithOctave(note note.Note) string {
	return fmt.Sprintf("%s%d", note.Class.String(config.AdjSymbolType), note.Octave)
}
