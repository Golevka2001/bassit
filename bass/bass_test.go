package bass

import (
	"testing"

	"github.com/Golevka2001/bassit/config"

	"github.com/go-music-theory/music-theory/note"
	"github.com/stretchr/testify/assert"
)

func TestNewBassModel(t *testing.T) {
	tests := []struct {
		name    string
		tuning  [config.StringCnt]string
		wantErr bool
	}{
		{
			name:    "valid standard bass tuning",
			tuning:  [config.StringCnt]string{"G2", "D2", "A1", "E1"},
			wantErr: false,
		},
		{
			name:    "valid half-step down tuning",
			tuning:  [config.StringCnt]string{"F#2", "C#2", "G#1", "D#1"},
			wantErr: false,
		},
		{
			name:    "no octave specified",
			tuning:  [config.StringCnt]string{"G", "D", "A", "E"},
			wantErr: true,
		},
		{
			name:    "invalid note name",
			tuning:  [config.StringCnt]string{"G2", "D2", "A1", "K1"},
			wantErr: true,
		},
		{
			name:    "negative octave note",
			tuning:  [config.StringCnt]string{"G2", "D2", "A1", "E-1"},
			wantErr: true,
		},
		{
			name:    "5-string bass tuning",
			tuning:  [config.StringCnt]string{"G2", "D2", "A1", "E1"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bass, err := New(tt.tuning)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if bass == nil {
				t.Errorf("Returned nil bass model")
				return
			}

			assert.Equal(t, config.StringCnt, len(bass.Strings), "Expected %d strings in bass model", config.StringCnt)

			// Verify each string has the correct base note
			for i, noteName := range tt.tuning {
				expectedNote := *note.Named(noteName)
				assert.Equal(t, expectedNote, bass.Strings[i].baseNote, "Expected base note for string %d to be %v", i, expectedNote)
			}
		})
	}
}

func TestGetBaseNotes(t *testing.T) {
	tuning := [config.StringCnt]string{"G2", "D2", "A1", "E1"}
	bass, err := New(tuning)
	if err != nil {
		t.Fatalf("Failed to create bass model: %v", err)
	}

	baseNotes := bass.GetBaseNotes()
	assert.Equal(t, config.StringCnt, len(baseNotes), "Expected %d base notes", config.StringCnt)

	for i, noteName := range tuning {
		expectedNote := *note.Named(noteName)
		assert.Equal(t, expectedNote, baseNotes[i], "Expected base note for string %d to be %v", i, expectedNote)
	}
}

func TestGetLowestAndHighestNotes(t *testing.T) {
	tuning := [config.StringCnt]string{"G2", "D2", "A1", "E1"}
	bass, err := New(tuning)
	if err != nil {
		t.Fatalf("Failed to create bass model: %v", err)
	}

	lowest, highest := bass.GetLowestAndHighestNotes()

	expectedLowest := *note.Named("E1")
	expectedHighest := bass.Strings[0].fretIdxToNote[config.DisplayedFretCount]

	assert.Equal(t, expectedLowest, lowest, "Expected lowest note to be E1")
	assert.Equal(t, expectedHighest, highest, "Expected highest note to match the highest fret on the highest string")

	// Test with a different tuning
	altTuning := [config.StringCnt]string{"F#2", "C#2", "G#1", "D#1"}
	altBass, _ := New(altTuning)
	altLowest, altHighest := altBass.GetLowestAndHighestNotes()

	expectedAltLowest := *note.Named("D#1")
	expectedAltHighest := altBass.Strings[0].fretIdxToNote[config.DisplayedFretCount]

	assert.Equal(t, expectedAltLowest, altLowest, "Expected lowest note to be D#1")
	assert.Equal(t, expectedAltHighest, altHighest, "Expected highest note to match the highest fret on the highest string")
}

func TestPressAndRelease(t *testing.T) {
	tuning := [config.StringCnt]string{"G2", "D2", "A1", "E1"}
	bass, err := New(tuning)
	if err != nil {
		t.Fatalf("Failed to create bass model: %v", err)
	}

	// Test pressing a fret
	stringIdx, fretIdx := 0, 5
	bass.Press(stringIdx, fretIdx)
	assert.True(t, bass.Strings[stringIdx].fretPressedStates[fretIdx], "Expected fret to be pressed")
	assert.Equal(t, fretIdx, bass.Strings[stringIdx].CurValidFret, "Expected current valid fret to be updated")

	// Test pressing a higher fret on the same string
	higherFret := 7
	bass.Press(stringIdx, higherFret)
	assert.True(t, bass.Strings[stringIdx].fretPressedStates[fretIdx], "Expected first fret to still be pressed")
	assert.True(t, bass.Strings[stringIdx].fretPressedStates[higherFret], "Expected higher fret to be pressed")
	assert.Equal(t, higherFret, bass.Strings[stringIdx].CurValidFret, "Expected current valid fret to be the higher fret")

	// Test releasing the higher fret
	bass.Release(stringIdx, higherFret)
	assert.False(t, bass.Strings[stringIdx].fretPressedStates[higherFret], "Expected higher fret to be released")
	assert.Equal(t, fretIdx, bass.Strings[stringIdx].CurValidFret, "Expected current valid fret to be the lower fret")

	// Test releasing the lower fret
	bass.Release(stringIdx, fretIdx)
	assert.False(t, bass.Strings[stringIdx].fretPressedStates[fretIdx], "Expected fret to be released")
	assert.Equal(t, 0, bass.Strings[stringIdx].CurValidFret, "Expected current valid fret to be 0 (open string)")

	// Test pressing multiple strings
	bass.Press(0, 3)
	bass.Press(1, 5)
	bass.Press(2, 1)
	assert.Equal(t, 3, bass.Strings[0].CurValidFret)
	assert.Equal(t, 5, bass.Strings[1].CurValidFret)
	assert.Equal(t, 1, bass.Strings[2].CurValidFret)

	// Test invalid indices
	invalidStringIdx := config.StringCnt + 1
	invalidFretIdx := config.MaxFretCnt + 1

	// These should not panic
	bass.Press(invalidStringIdx, fretIdx)
	bass.Press(stringIdx, invalidFretIdx)
	bass.Press(-1, fretIdx)
	bass.Release(invalidStringIdx, fretIdx)
	bass.Release(stringIdx, invalidFretIdx)
	bass.Release(-1, fretIdx)
}

func TestPluckAndRestoreString(t *testing.T) {
	tuning := [config.StringCnt]string{"G2", "D2", "A1", "E1"}
	bass, err := New(tuning)
	if err != nil {
		t.Fatalf("Failed to create bass model: %v", err)
	}

	// Test plucking a string
	stringIdx := 1
	bass.PluckString(stringIdx, 0)
	assert.True(t, bass.Strings[stringIdx].IsVibrating, "Expected string to be vibrating")
	assert.True(t, bass.Strings[stringIdx].pluckStates[0], "Expected pluck state 0 to be true")
	assert.False(t, bass.Strings[stringIdx].pluckStates[1], "Expected pluck state 1 to be false")

	// Test plucking different position
	bass.PluckString(stringIdx, 1)
	assert.True(t, bass.Strings[stringIdx].IsVibrating, "Expected string to be vibrating")
	assert.False(t, bass.Strings[stringIdx].pluckStates[0], "Expected pluck state 0 to be false")
	assert.True(t, bass.Strings[stringIdx].pluckStates[1], "Expected pluck state 1 to be true")

	// Test restoring a string
	bass.RestoreString(stringIdx)
	assert.False(t, bass.Strings[stringIdx].IsVibrating, "Expected string not to be vibrating")
	assert.False(t, bass.Strings[stringIdx].pluckStates[0], "Expected pluck state 0 to be false")
	assert.False(t, bass.Strings[stringIdx].pluckStates[1], "Expected pluck state 1 to be false")

	// Test plucking multiple strings
	bass.PluckString(0, 0)
	bass.PluckString(2, 1)
	assert.True(t, bass.Strings[0].IsVibrating)
	assert.False(t, bass.Strings[1].IsVibrating)
	assert.True(t, bass.Strings[2].IsVibrating)

	// Test invalid indices
	invalidStringIdx := config.StringCnt + 1
	invalidPosition := 2

	// These should not panic
	bass.PluckString(invalidStringIdx, 0)
	bass.PluckString(stringIdx, invalidPosition)
	bass.PluckString(stringIdx, -1)
	bass.RestoreString(invalidStringIdx)
	bass.RestoreString(-1)
}

func TestCombinedOperations(t *testing.T) {
	tuning := [config.StringCnt]string{"G2", "D2", "A1", "E1"}
	bass, err := New(tuning)
	if err != nil {
		t.Fatalf("Failed to create bass model: %v", err)
	}

	// Test combined operations: pressing and plucking
	stringIdx := 2
	bass.Press(stringIdx, 3)
	bass.PluckString(stringIdx, 0)

	assert.Equal(t, 3, bass.Strings[stringIdx].CurValidFret)
	assert.True(t, bass.Strings[stringIdx].IsVibrating)

	// The note played should be the fretted note
	expectedNote := bass.Strings[stringIdx].fretIdxToNote[3]
	actualNote := bass.Strings[stringIdx].GetNoteToPlay()
	assert.Equal(t, expectedNote, actualNote)

	// Restore the string but keep fret pressed
	bass.RestoreString(stringIdx)
	assert.Equal(t, 3, bass.Strings[stringIdx].CurValidFret)
	assert.False(t, bass.Strings[stringIdx].IsVibrating)

	// Releasing the fret should change CurValidFret
	bass.Release(stringIdx, 3)
	assert.Equal(t, 0, bass.Strings[stringIdx].CurValidFret)
}
