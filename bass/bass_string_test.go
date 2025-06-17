package bass

import (
	"testing"

	"github.com/Golevka2001/bassit/config"

	"github.com/go-music-theory/music-theory/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBassStringModel(t *testing.T) {
	baseNote := *note.Named("E1")
	bsm, err := newBassString(baseNote)

	require.NoError(t, err)
	assert.Equal(t, baseNote, bsm.BaseNote)
	assert.Equal(t, baseNote, bsm.FretToNote[0])
	assert.Len(t, bsm.FretPressedStates, config.MaxFretCnt+1)
	assert.False(t, bsm.PluckedState)
	assert.Equal(t, 0, bsm.CurValidFret)

	// Verify fret-to-note mapping
	expectedNoteNames := [config.MaxFretCnt + 1]string{
		"E1", "F1", "F#1", "G1", "G#1", "A1", "A#1", "B1",
		"C2", "C#2", "D2", "D#2", "E2", "F2", "F#2", "G2",
		"G#2", "A2", "A#2", "B2", "C3", "C#3", "D3", "D#3",
		"E3",
	}
	expectedNotes := make(map[int]note.Note)
	for i, name := range expectedNoteNames {
		expectedNotes[i] = *note.Named(name)
	}
	for i := 0; i <= config.MaxFretCnt; i++ {
		assert.Equal(t, expectedNotes[i], bsm.FretToNote[i], "Fret %d should map to %s", i, expectedNoteNames[i])
	}
}

func TestNewBassStringModel_InvalidBaseNote(t *testing.T) {
	// Test with invalid base note (no octave)
	_, err := newBassString(*note.Named("A"))
	assert.Error(t, err, "Expected error for base note without octave")

	// Test with invalid base note (nil note)
	_, err = newBassString(*note.Named("J"))
	assert.Error(t, err, "Expected error for nil base note")

	// Test with invalid base note (negative octave)
	_, err = newBassString(*note.Named("A-1"))
	assert.Error(t, err, "Expected error for base note with negative octave")
}

func TestPressFret(t *testing.T) {
	baseNote := *note.Named("A1")
	bsm, _ := newBassString(baseNote)

	// Test pressing valid fret
	bsm.pressFret(3)
	assert.True(t, bsm.FretPressedStates[3])
	assert.Equal(t, 3, bsm.CurValidFret)

	// Test pressing higher fret
	bsm.pressFret(5)
	assert.True(t, bsm.FretPressedStates[5])
	assert.Equal(t, 5, bsm.CurValidFret)

	// Test pressing lower fret (should not change CurValidFret)
	bsm.pressFret(2)
	assert.True(t, bsm.FretPressedStates[2])
	assert.Equal(t, 5, bsm.CurValidFret)

	// Test pressing invalid frets
	bsm.pressFret(-1)
	bsm.pressFret(config.MaxFretCnt + 1)
	assert.Equal(t, 5, bsm.CurValidFret)
}

func TestReleaseFret(t *testing.T) {
	baseNote := *note.Named("D2")
	bsm, _ := newBassString(baseNote)

	// Press multiple frets
	bsm.pressFret(2)
	bsm.pressFret(5)
	bsm.pressFret(7)
	assert.Equal(t, 7, bsm.CurValidFret)

	// Release middle fret (should not change CurValidFret)
	bsm.releaseFret(5)
	assert.False(t, bsm.FretPressedStates[5])
	assert.Equal(t, 7, bsm.CurValidFret)

	// Release highest fret (should update CurValidFret)
	bsm.releaseFret(7)
	assert.False(t, bsm.FretPressedStates[7])
	assert.Equal(t, 2, bsm.CurValidFret)

	// Release last fret (should reset to open string)
	bsm.releaseFret(2)
	assert.False(t, bsm.FretPressedStates[2])
	assert.Equal(t, 0, bsm.CurValidFret)

	// Test releasing invalid frets
	bsm.pressFret(3)
	originalFret := bsm.CurValidFret
	bsm.releaseFret(-1)
	bsm.releaseFret(config.MaxFretCnt + 1)
	assert.Equal(t, originalFret, bsm.CurValidFret)
}

func TestGetNoteToPlay(t *testing.T) {
	baseNote := *note.Named("D2")
	bsm, _ := newBassString(baseNote)

	// Test open string
	note := bsm.GetNoteToPlay()
	assert.Equal(t, baseNote, note)

	// Test pressed fret
	bsm.pressFret(3)
	note = bsm.GetNoteToPlay()
	assert.Equal(t, bsm.FretToNote[3], note)

	// Test multiple pressed frets (should return highest)
	bsm.pressFret(1)
	bsm.pressFret(5)
	note = bsm.GetNoteToPlay()
	assert.Equal(t, bsm.FretToNote[5], note)
}

func TestComplexFretOperations(t *testing.T) {
	baseNote := *note.Named("B1")
	bsm, _ := newBassString(baseNote)

	// Simulate complex pressing and releasing sequence
	bsm.pressFret(3)
	bsm.pressFret(5)
	bsm.pressFret(7)
	bsm.pressFret(2)

	assert.Equal(t, 7, bsm.CurValidFret)
	assert.Equal(t, bsm.FretToNote[7], bsm.GetNoteToPlay())

	// Release frets in different order
	bsm.releaseFret(3)
	bsm.releaseFret(7)
	assert.Equal(t, 5, bsm.CurValidFret)

	bsm.releaseFret(2)
	assert.Equal(t, 5, bsm.CurValidFret)

	bsm.releaseFret(5)
	assert.Equal(t, 0, bsm.CurValidFret)
	assert.Equal(t, baseNote, bsm.GetNoteToPlay())
}
