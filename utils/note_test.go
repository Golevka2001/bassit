package utils

import (
	"testing"

	"github.com/Golevka2001/bassit/config"

	"github.com/go-music-theory/music-theory/note"
	"github.com/stretchr/testify/assert"
)

func TestGetNoteStepFrom(t *testing.T) {
	tests := []struct {
		name     string
		from     note.Note
		inc      int
		expected note.Note
	}{
		{
			name:     "C4 + 2 steps = D4",
			from:     note.Note{Class: note.C, Octave: 4},
			inc:      2,
			expected: note.Note{Class: note.D, Octave: 4},
		},
		{
			name:     "C4 + 12 steps = C5",
			from:     note.Note{Class: note.C, Octave: 4},
			inc:      12,
			expected: note.Note{Class: note.C, Octave: 5},
		},
		{
			name:     "C4 - 1 step = B3",
			from:     note.Note{Class: note.C, Octave: 4},
			inc:      -1,
			expected: note.Note{Class: note.B, Octave: 3},
		},
		{
			name:     "inc == 0, should return original",
			from:     note.Note{Class: note.C, Octave: 4},
			inc:      0,
			expected: note.Note{Class: note.C, Octave: 4},
		},
		{
			name:     "A4 + 5 steps = D5",
			from:     note.Note{Class: note.A, Octave: 4},
			inc:      5,
			expected: note.Note{Class: note.D, Octave: 5},
		},
		{
			name:     "E2 - 7 steps = A1",
			from:     note.Note{Class: note.E, Octave: 2},
			inc:      -7,
			expected: note.Note{Class: note.A, Octave: 1},
		},
		{
			name:     "G3 + 24 steps = G5",
			from:     note.Note{Class: note.G, Octave: 3},
			inc:      24,
			expected: note.Note{Class: note.G, Octave: 5},
		},
		{
			name:     "F#4 - 13 steps = F3",
			from:     note.Note{Class: note.Fs, Octave: 4},
			inc:      -13,
			expected: note.Note{Class: note.F, Octave: 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := GetNoteStepFrom(tt.from, tt.inc)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestGetStepBetween(t *testing.T) {
	tests := []struct {
		name     string
		from     note.Note
		to       note.Note
		expected int
	}{
		{
			name:     "C4 to D4 = 2",
			from:     note.Note{Class: note.C, Octave: 4},
			to:       note.Note{Class: note.D, Octave: 4},
			expected: 2,
		},
		{
			name:     "C4 to C5 = 12",
			from:     note.Note{Class: note.C, Octave: 4},
			to:       note.Note{Class: note.C, Octave: 5},
			expected: 12,
		},
		{
			name:     "C4 to E3 = -8",
			from:     note.Note{Class: note.C, Octave: 4},
			to:       note.Note{Class: note.E, Octave: 3},
			expected: -8,
		},
		{
			name:     "C4 to B3 = -1",
			from:     note.Note{Class: note.C, Octave: 4},
			to:       note.Note{Class: note.B, Octave: 3},
			expected: -1,
		},
		{
			name:     "same note = 0",
			from:     note.Note{Class: note.A, Octave: 4},
			to:       note.Note{Class: note.A, Octave: 4},
			expected: 0,
		},
		{
			name:     "A4 to A5 = 12",
			from:     note.Note{Class: note.A, Octave: 4},
			to:       note.Note{Class: note.A, Octave: 5},
			expected: 12,
		},
		{
			name:     "G3 to F#4 = 11",
			from:     note.Note{Class: note.G, Octave: 3},
			to:       note.Note{Class: note.Fs, Octave: 4},
			expected: 11,
		},
		{
			name:     "E2 to G1 = -9",
			from:     note.Note{Class: note.E, Octave: 2},
			to:       note.Note{Class: note.G, Octave: 1},
			expected: -9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := GetStepBetween(tt.from, tt.to)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestGetNoteNameWithOctave(t *testing.T) {
	originalAccidentalStyle := config.AccidentalStyle
	defer func() { config.AccidentalStyle = originalAccidentalStyle }()

	tests := []struct {
		name            string
		note            note.Note
		accidentalStyle note.AdjSymbol
		expected        string
	}{
		{
			name:            "A2 natural note",
			note:            note.Note{Class: note.A, Octave: 2},
			accidentalStyle: note.Sharp,
			expected:        "A2",
		},
		{
			name:            "C#3 with sharp style",
			note:            note.Note{Class: note.Cs, Octave: 3},
			accidentalStyle: note.Sharp,
			expected:        "C#3",
		},
		{
			name:            "D#3 with flat style (becomes Eb3)",
			note:            note.Note{Class: note.Ds, Octave: 3},
			accidentalStyle: note.Flat,
			expected:        "Eb3",
		},
		{
			name:            "F#4 with sharp style",
			note:            note.Note{Class: note.Fs, Octave: 4},
			accidentalStyle: note.Sharp,
			expected:        "F#4",
		},
		{
			name:            "F#4 with flat style (becomes Gb4)",
			note:            note.Note{Class: note.Fs, Octave: 4},
			accidentalStyle: note.Flat,
			expected:        "Gb4",
		},
		{
			name:            "B0 low octave",
			note:            note.Note{Class: note.B, Octave: 0},
			accidentalStyle: note.Sharp,
			expected:        "B0",
		},
		{
			name:            "C8 high octave",
			note:            note.Note{Class: note.C, Octave: 8},
			accidentalStyle: note.Sharp,
			expected:        "C8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.AccidentalStyle = tt.accidentalStyle
			actual := GetNoteNameWithOctave(tt.note)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
