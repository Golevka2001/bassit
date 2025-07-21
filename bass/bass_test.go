package bass

import (
	"testing"

	"github.com/Golevka2001/bassit/config"

	"github.com/go-music-theory/music-theory/note"
	"github.com/stretchr/testify/assert"
)

func TestNewBass(t *testing.T) {
	tests := []struct {
		name     string
		tuning   [config.StringCnt]string
		expected error
		validate func(t *testing.T, bass *BassModel, err error)
	}{
		{
			name:     "standard bass tuning",
			tuning:   [config.StringCnt]string{"G2", "D2", "A1", "E1"},
			expected: nil,
			validate: func(t *testing.T, bass *BassModel, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, bass)
				assert.Equal(t, config.StringCnt, len(bass.strings))

				baseNotes := bass.GetBaseNotes()
				assert.Equal(t, config.StringCnt, len(baseNotes))
				assert.Equal(t, note.G, baseNotes[0].Class)
				assert.Equal(t, note.D, baseNotes[1].Class)
				assert.Equal(t, note.A, baseNotes[2].Class)
				assert.Equal(t, note.E, baseNotes[3].Class)
			},
		},
		{
			name:     "alternative tuning",
			tuning:   [config.StringCnt]string{"A2", "E2", "B1", "F#1"},
			expected: nil,
			validate: func(t *testing.T, bass *BassModel, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, bass)

				baseNotes := bass.GetBaseNotes()
				assert.Equal(t, note.A, baseNotes[0].Class)
				assert.Equal(t, note.E, baseNotes[1].Class)
				assert.Equal(t, note.B, baseNotes[2].Class)
				assert.Equal(t, note.Fs, baseNotes[3].Class)
			},
		},
		{
			name:     "invalid tuning - invalid note",
			tuning:   [config.StringCnt]string{"G2", "D2", "A1", "X1"},
			expected: assert.AnError,
			validate: func(t *testing.T, bass *BassModel, err error) {
				assert.Error(t, err)
				assert.Nil(t, bass)
			},
		},
		{
			name:     "invalid tuning - invalid octave",
			tuning:   [config.StringCnt]string{"G2", "D2", "A1", "E0"},
			expected: assert.AnError,
			validate: func(t *testing.T, bass *BassModel, err error) {
				assert.Error(t, err)
				assert.Nil(t, bass)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bass, err := NewBass(tt.tuning)

			if tt.validate != nil {
				tt.validate(t, bass, err)
			}
		})
	}
}

func TestBassModelGetRange(t *testing.T) {
	tuning := [config.StringCnt]string{"G2", "D2", "A1", "E1"}
	bass, err := NewBass(tuning)
	assert.NoError(t, err)
	assert.NotNil(t, bass)

	tests := []struct {
		name     string
		validate func(t *testing.T, bass *BassModel)
	}{
		{
			name: "get displayed range",
			validate: func(t *testing.T, bass *BassModel) {
				lowest, highest := bass.GetRange()

				assert.Equal(t, note.E, lowest.Class)
				assert.Equal(t, note.Octave(1), lowest.Octave)

				assert.Equal(t, note.G, highest.Class)
			},
		},
		{
			name: "get actual range",
			validate: func(t *testing.T, bass *BassModel) {
				lowest, highest := bass.GetActualRange()

				assert.Equal(t, note.E, lowest.Class)
				assert.Equal(t, note.Octave(1), lowest.Octave)

				assert.Equal(t, note.G, highest.Class)
			},
		},
		{
			name: "range consistency",
			validate: func(t *testing.T, bass *BassModel) {
				displayedLow, displayedHigh := bass.GetRange()
				actualLow, actualHigh := bass.GetActualRange()

				assert.Equal(t, displayedLow, actualLow)

				assert.GreaterOrEqual(t, actualHigh.Octave, displayedHigh.Octave)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t, bass)
		})
	}
}

func TestBassModelFretOperations(t *testing.T) {
	tuning := [config.StringCnt]string{"G2", "D2", "A1", "E1"}
	bass, err := NewBass(tuning)
	assert.NoError(t, err)
	assert.NotNil(t, bass)

	tests := []struct {
		name     string
		actions  func(t *testing.T, bass *BassModel)
		validate func(t *testing.T, bass *BassModel)
	}{
		{
			name: "press and check fret",
			actions: func(t *testing.T, bass *BassModel) {
				pos := FretboardPosition{StringIdx: 0, FretIdx: 5}
				bass.PressFret(pos)
			},
			validate: func(t *testing.T, bass *BassModel) {
				pos := FretboardPosition{StringIdx: 0, FretIdx: 5}
				assert.True(t, bass.IsFretPressed(pos))
				assert.Equal(t, 5, bass.GetValidFretIdxOfString(0))
			},
		},
		{
			name: "press multiple frets on same string",
			actions: func(t *testing.T, bass *BassModel) {
				bass.PressFret(FretboardPosition{StringIdx: 1, FretIdx: 3})
				bass.PressFret(FretboardPosition{StringIdx: 1, FretIdx: 7})
				bass.PressFret(FretboardPosition{StringIdx: 1, FretIdx: 5})
			},
			validate: func(t *testing.T, bass *BassModel) {
				assert.True(t, bass.IsFretPressed(FretboardPosition{StringIdx: 1, FretIdx: 3}))
				assert.True(t, bass.IsFretPressed(FretboardPosition{StringIdx: 1, FretIdx: 5}))
				assert.True(t, bass.IsFretPressed(FretboardPosition{StringIdx: 1, FretIdx: 7}))
				assert.Equal(t, 7, bass.GetValidFretIdxOfString(1))
			},
		},
		{
			name: "release fret",
			actions: func(t *testing.T, bass *BassModel) {
				bass.PressFret(FretboardPosition{StringIdx: 2, FretIdx: 3})
				bass.PressFret(FretboardPosition{StringIdx: 2, FretIdx: 7})
				bass.ReleaseFret(FretboardPosition{StringIdx: 2, FretIdx: 7})
			},
			validate: func(t *testing.T, bass *BassModel) {
				assert.True(t, bass.IsFretPressed(FretboardPosition{StringIdx: 2, FretIdx: 3}))
				assert.False(t, bass.IsFretPressed(FretboardPosition{StringIdx: 2, FretIdx: 7}))
				assert.Equal(t, 3, bass.GetValidFretIdxOfString(2))
			},
		},
		{
			name: "operations on invalid string indices",
			actions: func(t *testing.T, bass *BassModel) {
				bass.PressFret(FretboardPosition{StringIdx: -1, FretIdx: 5})
				bass.PressFret(FretboardPosition{StringIdx: config.StringCnt, FretIdx: 5})
			},
			validate: func(t *testing.T, bass *BassModel) {
				assert.False(t, bass.IsFretPressed(FretboardPosition{StringIdx: -1, FretIdx: 5}))
				assert.False(t, bass.IsFretPressed(FretboardPosition{StringIdx: config.StringCnt, FretIdx: 5}))
				assert.Equal(t, -1, bass.GetValidFretIdxOfString(-1))
				assert.Equal(t, -1, bass.GetValidFretIdxOfString(config.StringCnt))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bass, _ = NewBass(tuning)

			if tt.actions != nil {
				tt.actions(t, bass)
			}

			if tt.validate != nil {
				tt.validate(t, bass)
			}
		})
	}
}

func TestBassModelNoteOperations(t *testing.T) {
	tuning := [config.StringCnt]string{"G2", "D2", "A1", "E1"}
	bass, err := NewBass(tuning)
	assert.NoError(t, err)
	assert.NotNil(t, bass)

	tests := []struct {
		name     string
		actions  func(t *testing.T, bass *BassModel)
		validate func(t *testing.T, bass *BassModel)
	}{
		{
			name: "get note at position",
			actions: func(t *testing.T, bass *BassModel) {
			},
			validate: func(t *testing.T, bass *BassModel) {
				note0 := bass.GetNoteAt(FretboardPosition{StringIdx: 0, FretIdx: 0})
				assert.NotNil(t, note0)
				assert.Equal(t, note.G, note0.Class)

				note3 := bass.GetNoteAt(FretboardPosition{StringIdx: 3, FretIdx: 0})
				assert.NotNil(t, note3)
				assert.Equal(t, note.E, note3.Class)

				note1_5 := bass.GetNoteAt(FretboardPosition{StringIdx: 1, FretIdx: 5})
				assert.NotNil(t, note1_5)

				invalidNote := bass.GetNoteAt(FretboardPosition{StringIdx: -1, FretIdx: 0})
				assert.Nil(t, invalidNote)

				invalidNote2 := bass.GetNoteAt(FretboardPosition{StringIdx: config.StringCnt, FretIdx: 0})
				assert.Nil(t, invalidNote2)
			},
		},
		{
			name: "get note to play",
			actions: func(t *testing.T, bass *BassModel) {
				bass.PressFret(FretboardPosition{StringIdx: 0, FretIdx: 3})
				bass.PressFret(FretboardPosition{StringIdx: 1, FretIdx: 5})
			},
			validate: func(t *testing.T, bass *BassModel) {
				note0 := bass.GetNoteToPlay(0)
				expectedNote0 := bass.GetNoteAt(FretboardPosition{StringIdx: 0, FretIdx: 3})
				assert.Equal(t, expectedNote0, note0)

				note1 := bass.GetNoteToPlay(1)
				expectedNote1 := bass.GetNoteAt(FretboardPosition{StringIdx: 1, FretIdx: 5})
				assert.Equal(t, expectedNote1, note1)

				note2 := bass.GetNoteToPlay(2)
				expectedNote2 := bass.GetNoteAt(FretboardPosition{StringIdx: 2, FretIdx: 0})
				assert.Equal(t, expectedNote2, note2)

				invalidNote := bass.GetNoteToPlay(-1)
				assert.Nil(t, invalidNote)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bass, _ = NewBass(tuning)

			if tt.actions != nil {
				tt.actions(t, bass)
			}

			if tt.validate != nil {
				tt.validate(t, bass)
			}
		})
	}
}

func TestBassModelPluckOperations(t *testing.T) {
	tuning := [config.StringCnt]string{"G2", "D2", "A1", "E1"}
	bass, err := NewBass(tuning)
	assert.NoError(t, err)
	assert.NotNil(t, bass)

	tests := []struct {
		name     string
		actions  func(t *testing.T, bass *BassModel) int64
		validate func(t *testing.T, bass *BassModel, pluckTime int64)
	}{
		{
			name: "pluck string",
			actions: func(t *testing.T, bass *BassModel) int64 {
				return bass.PluckString(0, PluckTypeNormal1)
			},
			validate: func(t *testing.T, bass *BassModel, pluckTime int64) {
				assert.True(t, bass.IsStringVibrating(0))
				assert.Greater(t, pluckTime, int64(0))
			},
		},
		{
			name: "stop vibrating string with correct time",
			actions: func(t *testing.T, bass *BassModel) int64 {
				pluckTime := bass.PluckString(1, PluckTypeSlap1)
				stopped := bass.StopVibratingString(1, pluckTime)
				assert.True(t, stopped)
				return pluckTime
			},
			validate: func(t *testing.T, bass *BassModel, pluckTime int64) {
				assert.False(t, bass.IsStringVibrating(1))
			},
		},
		{
			name: "stop vibrating string with incorrect time",
			actions: func(t *testing.T, bass *BassModel) int64 {
				pluckTime := bass.PluckString(2, PluckTypeMute1)
				stopped := bass.StopVibratingString(2, pluckTime+1)
				assert.False(t, stopped)
				return pluckTime
			},
			validate: func(t *testing.T, bass *BassModel, pluckTime int64) {
				assert.True(t, bass.IsStringVibrating(2))
			},
		},
		{
			name: "stop vibrating without check",
			actions: func(t *testing.T, bass *BassModel) int64 {
				pluckTime := bass.PluckString(3, PluckTypeNormal2)
				bass.StopVibratingStringWithoutCheck(3)
				return pluckTime
			},
			validate: func(t *testing.T, bass *BassModel, pluckTime int64) {
				assert.False(t, bass.IsStringVibrating(3))
			},
		},
		{
			name: "operations on invalid string indices",
			actions: func(t *testing.T, bass *BassModel) int64 {
				pluckTime1 := bass.PluckString(-1, PluckTypeNormal1)
				pluckTime2 := bass.PluckString(config.StringCnt, PluckTypeNormal1)
				assert.Equal(t, int64(0), pluckTime1)
				assert.Equal(t, int64(0), pluckTime2)

				stopped1 := bass.StopVibratingString(-1, 123)
				stopped2 := bass.StopVibratingString(config.StringCnt, 123)
				assert.False(t, stopped1)
				assert.False(t, stopped2)

				return 0
			},
			validate: func(t *testing.T, bass *BassModel, pluckTime int64) {
				assert.False(t, bass.IsStringVibrating(-1))
				assert.False(t, bass.IsStringVibrating(config.StringCnt))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bass, _ = NewBass(tuning)

			var pluckTime int64
			if tt.actions != nil {
				pluckTime = tt.actions(t, bass)
			}

			if tt.validate != nil {
				tt.validate(t, bass, pluckTime)
			}
		})
	}
}
