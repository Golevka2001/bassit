package model

import (
	"testing"

	C "github.com/Golevka2001/bassit/constant"

	"github.com/go-music-theory/music-theory/note"
	"github.com/stretchr/testify/assert"
)

func TestNewBassModel(t *testing.T) {
	tests := []struct {
		name    string
		tuning  [C.StringCnt]string
		wantErr bool
	}{
		{
			name:    "valid standard bass tuning",
			tuning:  [C.StringCnt]string{"G2", "D2", "A1", "E1"},
			wantErr: false,
		},
		{
			name:    "valid half-step down tuning",
			tuning:  [C.StringCnt]string{"F#2", "C#2", "G#1", "D#1"},
			wantErr: false,
		},
		{
			name:    "no octave specified",
			tuning:  [C.StringCnt]string{"G", "D", "A", "E"},
			wantErr: true,
		},
		{
			name:    "invalid note name",
			tuning:  [C.StringCnt]string{"G2", "D2", "A1", "K1"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bass, err := NewBassModel(tt.tuning)

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

			assert.Equal(t, C.StringCnt, len(bass.Strings), "Expected %d strings in bass model", C.StringCnt)

			// Verify each string has the correct base note
			for i, noteName := range tt.tuning {
				expectedNote := *note.Named(noteName)
				assert.Equal(t, expectedNote, bass.Strings[i].BaseNote, "Expected base note for string %d to be %v", i, expectedNote)
			}
		})
	}
}
