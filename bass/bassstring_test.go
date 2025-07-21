package bass

import (
	"testing"
	"time"

	"github.com/Golevka2001/bassit/config"

	"github.com/go-music-theory/music-theory/note"
	"github.com/stretchr/testify/assert"
)

func TestNewBassString(t *testing.T) {
	tests := []struct {
		name     string
		baseNote note.Note
		expected error
		validate func(t *testing.T, bsm *bassStringModel, err error)
	}{
		{
			name:     "valid bass string creation - E1",
			baseNote: note.Note{Class: note.E, Octave: 1},
			expected: nil,
			validate: func(t *testing.T, bsm *bassStringModel, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, bsm)

				baseNote := bsm.GetNoteAt(0)
				assert.NotNil(t, baseNote)
				assert.Equal(t, note.E, baseNote.Class)
				assert.Equal(t, note.Octave(1), baseNote.Octave)

				assert.Equal(t, config.MaxFretCnt+1, len(bsm.fretIdxToNote))

				assert.Equal(t, 0, bsm.rightmostPressedFret)
				assert.False(t, bsm.isVibrating)
			},
		},
		{
			name:     "valid bass string creation - A1",
			baseNote: note.Note{Class: note.A, Octave: 1},
			expected: nil,
			validate: func(t *testing.T, bsm *bassStringModel, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, bsm)

				baseNote := bsm.GetNoteAt(0)
				assert.Equal(t, note.A, baseNote.Class)
				assert.Equal(t, note.Octave(1), baseNote.Octave)
			},
		},
		{
			name:     "invalid bass string - nil class",
			baseNote: note.Note{Class: note.Nil, Octave: 1},
			expected: assert.AnError,
			validate: func(t *testing.T, bsm *bassStringModel, err error) {
				assert.Error(t, err)
				assert.Nil(t, bsm)
				assert.Contains(t, err.Error(), "invalid base note")
			},
		},
		{
			name:     "invalid bass string - zero octave",
			baseNote: note.Note{Class: note.C, Octave: 0},
			expected: assert.AnError,
			validate: func(t *testing.T, bsm *bassStringModel, err error) {
				assert.Error(t, err)
				assert.Nil(t, bsm)
				assert.Contains(t, err.Error(), "invalid base note")
			},
		},
		{
			name:     "invalid bass string - negative octave",
			baseNote: note.Note{Class: note.G, Octave: -1},
			expected: assert.AnError,
			validate: func(t *testing.T, bsm *bassStringModel, err error) {
				assert.Error(t, err)
				assert.Nil(t, bsm)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bsm, err := NewBassString(tt.baseNote)

			if tt.validate != nil {
				tt.validate(t, bsm, err)
			}
		})
	}
}

func TestBassStringGetNoteAt(t *testing.T) {
	baseNote := note.Note{Class: note.E, Octave: 1}
	bsm, err := NewBassString(baseNote)
	assert.NoError(t, err)
	assert.NotNil(t, bsm)

	tests := []struct {
		name     string
		fretIdx  int
		expected *note.Note
	}{
		{
			name:     "open string (fret 0)",
			fretIdx:  0,
			expected: &note.Note{Class: note.E, Octave: 1},
		},
		{
			name:     "first fret",
			fretIdx:  1,
			expected: &note.Note{Class: note.F, Octave: 1},
		},
		{
			name:     "twelfth fret (one octave up)",
			fretIdx:  12,
			expected: &note.Note{Class: note.E, Octave: 2},
		},
		{
			name:     "maximum fret",
			fretIdx:  config.MaxFretCnt,
			expected: nil},
		{
			name:     "invalid fret - negative",
			fretIdx:  -1,
			expected: nil,
		},
		{
			name:     "invalid fret - too high",
			fretIdx:  config.MaxFretCnt + 1,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := bsm.GetNoteAt(tt.fretIdx)

			if tt.expected == nil && tt.fretIdx >= 0 && tt.fretIdx <= config.MaxFretCnt {
				assert.NotNil(t, actual)
			} else if tt.expected == nil {
				assert.Nil(t, actual)
			} else {
				assert.NotNil(t, actual)
				assert.Equal(t, tt.expected.Class, actual.Class)
				assert.Equal(t, tt.expected.Octave, actual.Octave)
			}
		})
	}
}

func TestBassStringFretOperations(t *testing.T) {
	baseNote := note.Note{Class: note.E, Octave: 1}
	bsm, err := NewBassString(baseNote)
	assert.NoError(t, err)
	assert.NotNil(t, bsm)

	tests := []struct {
		name     string
		actions  func(t *testing.T, bsm *bassStringModel)
		validate func(t *testing.T, bsm *bassStringModel)
	}{
		{
			name: "press single fret",
			actions: func(t *testing.T, bsm *bassStringModel) {
				bsm.PressFret(5)
			},
			validate: func(t *testing.T, bsm *bassStringModel) {
				assert.True(t, bsm.IsFretPressed(5))
				assert.Equal(t, 5, bsm.rightmostPressedFret)

				noteToPlay := bsm.GetNoteToPlay()
				expectedNote := bsm.GetNoteAt(5)
				assert.Equal(t, expectedNote, noteToPlay)
			},
		},
		{
			name: "press multiple frets",
			actions: func(t *testing.T, bsm *bassStringModel) {
				bsm.PressFret(3)
				bsm.PressFret(7)
				bsm.PressFret(5)
			},
			validate: func(t *testing.T, bsm *bassStringModel) {
				assert.True(t, bsm.IsFretPressed(3))
				assert.True(t, bsm.IsFretPressed(5))
				assert.True(t, bsm.IsFretPressed(7))
				assert.Equal(t, 7, bsm.rightmostPressedFret)
			},
		},
		{
			name: "release fret - not rightmost",
			actions: func(t *testing.T, bsm *bassStringModel) {
				bsm.PressFret(3)
				bsm.PressFret(7)
				bsm.PressFret(5)
				bsm.ReleaseFret(5)
			},
			validate: func(t *testing.T, bsm *bassStringModel) {
				assert.True(t, bsm.IsFretPressed(3))
				assert.False(t, bsm.IsFretPressed(5))
				assert.True(t, bsm.IsFretPressed(7))
				assert.Equal(t, 7, bsm.rightmostPressedFret)
			},
		},
		{
			name: "release rightmost fret",
			actions: func(t *testing.T, bsm *bassStringModel) {
				bsm.PressFret(3)
				bsm.PressFret(7)
				bsm.ReleaseFret(7)
			},
			validate: func(t *testing.T, bsm *bassStringModel) {
				assert.True(t, bsm.IsFretPressed(3))
				assert.False(t, bsm.IsFretPressed(7))
				assert.Equal(t, 3, bsm.rightmostPressedFret)
			},
		},
		{
			name: "release all frets",
			actions: func(t *testing.T, bsm *bassStringModel) {
				bsm.PressFret(3)
				bsm.PressFret(7)
				bsm.ReleaseFret(7)
				bsm.ReleaseFret(3)
			},
			validate: func(t *testing.T, bsm *bassStringModel) {
				assert.False(t, bsm.IsFretPressed(3))
				assert.False(t, bsm.IsFretPressed(7))
				assert.Equal(t, 0, bsm.rightmostPressedFret)
				noteToPlay := bsm.GetNoteToPlay()
				openNote := bsm.GetNoteAt(0)
				assert.Equal(t, openNote, noteToPlay)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bsm.fretPressed = [config.MaxFretCnt + 1]bool{}
			bsm.rightmostPressedFret = 0

			if tt.actions != nil {
				tt.actions(t, bsm)
			}

			if tt.validate != nil {
				tt.validate(t, bsm)
			}
		})
	}
}

func TestBassStringPluckOperations(t *testing.T) {
	baseNote := note.Note{Class: note.E, Octave: 1}
	bsm, err := NewBassString(baseNote)
	assert.NoError(t, err)
	assert.NotNil(t, bsm)

	tests := []struct {
		name     string
		actions  func(t *testing.T, bsm *bassStringModel) int64
		validate func(t *testing.T, bsm *bassStringModel, pluckTime int64)
	}{
		{
			name: "pluck string",
			actions: func(t *testing.T, bsm *bassStringModel) int64 {
				return bsm.Pluck(PluckTypeNormal1)
			},
			validate: func(t *testing.T, bsm *bassStringModel, pluckTime int64) {
				assert.True(t, bsm.isVibrating)
				assert.Equal(t, PluckTypeNormal1, bsm.pluckType)
				assert.Equal(t, pluckTime, bsm.lastPluckTime)
				assert.Greater(t, pluckTime, int64(0))
			},
		},
		{
			name: "stop vibrating with correct time",
			actions: func(t *testing.T, bsm *bassStringModel) int64 {
				pluckTime := bsm.Pluck(PluckTypeSlap1)
				stopped := bsm.StopVibrating(pluckTime)
				assert.True(t, stopped)
				return pluckTime
			},
			validate: func(t *testing.T, bsm *bassStringModel, pluckTime int64) {
				assert.False(t, bsm.isVibrating)
			},
		},
		{
			name: "stop vibrating with incorrect time",
			actions: func(t *testing.T, bsm *bassStringModel) int64 {
				pluckTime := bsm.Pluck(PluckTypeMute1)
				stopped := bsm.StopVibrating(pluckTime + 1)
				assert.False(t, stopped)
				return pluckTime
			},
			validate: func(t *testing.T, bsm *bassStringModel, pluckTime int64) {
				assert.True(t, bsm.isVibrating)
			},
		},
		{
			name: "stop vibrating without check",
			actions: func(t *testing.T, bsm *bassStringModel) int64 {
				pluckTime := bsm.Pluck(PluckTypeNormal2)
				bsm.StopVibratingWithoutCheck()
				return pluckTime
			},
			validate: func(t *testing.T, bsm *bassStringModel, pluckTime int64) {
				assert.False(t, bsm.isVibrating)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bsm.isVibrating = false
			bsm.lastPluckTime = 0

			var pluckTime int64
			if tt.actions != nil {
				pluckTime = tt.actions(t, bsm)
			}

			if tt.validate != nil {
				tt.validate(t, bsm, pluckTime)
			}
		})
	}
}

func TestBassStringBoundaryConditions(t *testing.T) {
	baseNote := note.Note{Class: note.E, Octave: 1}
	bsm, err := NewBassString(baseNote)
	assert.NoError(t, err)
	assert.NotNil(t, bsm)

	tests := []struct {
		name     string
		validate func(t *testing.T, bsm *bassStringModel)
	}{
		{
			name: "press invalid frets",
			validate: func(t *testing.T, bsm *bassStringModel) {
				bsm.PressFret(-1)
				assert.False(t, bsm.IsFretPressed(-1))
				assert.Equal(t, 0, bsm.rightmostPressedFret)

				bsm.PressFret(config.MaxFretCnt + 1)
				assert.False(t, bsm.IsFretPressed(config.MaxFretCnt+1))
				assert.Equal(t, 0, bsm.rightmostPressedFret)
			},
		},
		{
			name: "release invalid frets",
			validate: func(t *testing.T, bsm *bassStringModel) {
				bsm.PressFret(5)
				assert.Equal(t, 5, bsm.rightmostPressedFret)

				bsm.ReleaseFret(-1)
				bsm.ReleaseFret(config.MaxFretCnt + 1)

				assert.Equal(t, 5, bsm.rightmostPressedFret)
				assert.True(t, bsm.IsFretPressed(5))
			},
		},
		{
			name: "pluck timing consistency",
			validate: func(t *testing.T, bsm *bassStringModel) {
				time1 := bsm.Pluck(PluckTypeNormal1)
				time.Sleep(time.Millisecond)
				time2 := bsm.Pluck(PluckTypeNormal2)

				assert.NotEqual(t, time1, time2)
				assert.Greater(t, time2, time1)
				assert.Equal(t, time2, bsm.lastPluckTime)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bsm.fretPressed = [config.MaxFretCnt + 1]bool{}
			bsm.rightmostPressedFret = 0
			bsm.isVibrating = false
			bsm.lastPluckTime = 0

			if tt.validate != nil {
				tt.validate(t, bsm)
			}
		})
	}
}
