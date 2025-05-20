package util

import (
    "testing"

    "github.com/go-music-theory/music-theory/note"
)

func TestGetNoteStepFrom(t *testing.T) {
    // C4 + 2 steps = D4
    from := note.Note{Class: note.C, Octave: 4}
    got := GetNoteStepFrom(from, 2)
    want := note.Note{Class: note.D, Octave: 4}
    if got.Class != want.Class || got.Octave != want.Octave {
        t.Errorf("expected %v, got %v", want, got)
    }

    // C4 + 12 steps = C5
    got = GetNoteStepFrom(from, 12)
    want = note.Note{Class: note.C, Octave: 5}
    if got.Class != want.Class || got.Octave != want.Octave {
        t.Errorf("expected %v, got %v", want, got)
    }

    // C4 - 1 step = B3
    got = GetNoteStepFrom(from, -1)
    want = note.Note{Class: note.B, Octave: 3}
    if got.Class != want.Class || got.Octave != want.Octave {
        t.Errorf("expected %v, got %v", want, got)
    }

    // inc == 0, should return original
    got = GetNoteStepFrom(from, 0)
    if got != from {
        t.Errorf("expected %v, got %v", from, got)
    }
}

func TestGetStepBetween(t *testing.T) {
    // C4 to D4 = 2
    from := note.Note{Class: note.C, Octave: 4}
    to := note.Note{Class: note.D, Octave: 4}
    if steps := GetStepBetween(from, to); steps != 2 {
        t.Errorf("expected 2, got %d", steps)
    }

    // C4 to C5 = 12
    to = note.Note{Class: note.C, Octave: 5}
    if steps := GetStepBetween(from, to); steps != 12 {
        t.Errorf("expected 12, got %d", steps)
    }

    // C4 to E3 = -8
    from = note.Note{Class: note.C, Octave: 4}
    to = note.Note{Class: note.E, Octave: 3}
    if steps := GetStepBetween(from, to); steps != -8 {
        t.Errorf("expected -8, got %d", steps)
    }

    // C4 to B3 = 1
    from = note.Note{Class: note.C, Octave: 4}
    to = note.Note{Class: note.B, Octave: 3}
    if steps := GetStepBetween(from, to); steps != -1 {
        t.Errorf("expected -1, got %d", steps)
    }
}

func TestGetNoteNameWithOctave(t *testing.T) {
    // C#3
    n := note.Note{Class: note.Cs, Octave: 3}
    got := GetNoteNameWithOctave(n)
    want := "C#3"
    if got != want {
        t.Errorf("expected %s, got %s", want, got)
    }

    // A2
    n = note.Note{Class: note.A, Octave: 2}
    got = GetNoteNameWithOctave(n)
    want = "A2"
    if got != want {
        t.Errorf("expected %s, got %s", want, got)
    }
}
