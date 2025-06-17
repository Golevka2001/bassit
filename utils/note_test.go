package utils

import (
	"testing"

	"github.com/Golevka2001/bassit/config"

	"github.com/go-music-theory/music-theory/note"
	"github.com/stretchr/testify/assert"
)

func TestGetNoteStepFrom(t *testing.T) {
	from := note.Note{Class: note.C, Octave: 4}

	// C4 + 2 steps = D4
	got := GetNoteStepFrom(from, 2)
	want := note.Note{Class: note.D, Octave: 4}
	assert.Equal(t, want, got, "expected %v, got %v", want, got)

	// C4 + 12 steps = C5
	got = GetNoteStepFrom(from, 12)
	want = note.Note{Class: note.C, Octave: 5}
	assert.Equal(t, want, got, "expected %v, got %v", want, got)

	// C4 - 1 step = B3
	got = GetNoteStepFrom(from, -1)
	want = note.Note{Class: note.B, Octave: 3}
	assert.Equal(t, want, got, "expected %v, got %v", want, got)

	// inc == 0, should return original
	got = GetNoteStepFrom(from, 0)
	assert.Equal(t, from, got, "expected %v, got %v", from, got)
}

func TestGetStepBetween(t *testing.T) {
	// C4 to D4 = 2
	from := note.Note{Class: note.C, Octave: 4}
	to := note.Note{Class: note.D, Octave: 4}
	assert.Equal(t, 2, GetStepBetween(from, to), "expected 2")

	// C4 to C5 = 12
	to = note.Note{Class: note.C, Octave: 5}
	assert.Equal(t, 12, GetStepBetween(from, to), "expected 12")

	// C4 to E3 = -8
	from = note.Note{Class: note.C, Octave: 4}
	to = note.Note{Class: note.E, Octave: 3}
	assert.Equal(t, -8, GetStepBetween(from, to), "expected -8")

	// C4 to B3 = -1
	from = note.Note{Class: note.C, Octave: 4}
	to = note.Note{Class: note.B, Octave: 3}
	assert.Equal(t, -1, GetStepBetween(from, to), "expected -1")
}

func TestGetNoteNameWithOctave(t *testing.T) {

	// A2
	n := note.Note{Class: note.A, Octave: 2}
	got := GetNoteNameWithOctave(n)
	want := "A2"
	assert.Equal(t, want, got, "expected %s, got %s", want, got)

	// Sharp
	config.AccidentalStyle = note.Sharp
	// C#3
	n = note.Note{Class: note.Cs, Octave: 3}
	got = GetNoteNameWithOctave(n)
	want = "C#3"
	assert.Equal(t, want, got, "expected %s, got %s", want, got)

	// Flat
	config.AccidentalStyle = note.Flat
	// Db3
	n = note.Note{Class: note.Ds, Octave: 3}
	got = GetNoteNameWithOctave(n)
	want = "Eb3"
	assert.Equal(t, want, got, "expected %s, got %s", want, got)
}
