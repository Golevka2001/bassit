package model

import (
	"fmt"

	C "github.com/Golevka2001/bassit/constant"
	"github.com/Golevka2001/bassit/util"

	"github.com/go-music-theory/music-theory/note"
)

type BassStringModel struct {
	// BaseNote is the open string note
	BaseNote note.Note
	// FretToNote maps the fret number to the corresponding note
	FretToNote map[int]note.Note
	// FretPressedStates records which fret is pressed, "0" is meaningless
	FretPressedStates []bool
	// PluckedState records whether the string is plucked
	PluckedState bool
	// CurValidFret records the "valid" fret number among the pressed frets
	// (e.g. if 3rd, 5th, and 7th frets are pressed, `CurValidFret` should be 7)
	CurValidFret int
}

func NewBassStringModel(baseNote note.Note) (*BassStringModel, error) {
	if baseNote.Class == note.Nil || baseNote.Octave <= 0 {
		return nil, fmt.Errorf("invalid base note")
	}

	fretToNote := make(map[int]note.Note)
	fretToNote[0] = baseNote
	for i := 1; i <= C.MaxFretCnt; i++ {
		fretToNote[i] = util.GetNoteStepFrom(fretToNote[i-1], 1)
	}

	return &BassStringModel{
		BaseNote:          baseNote,
		FretToNote:        fretToNote,
		FretPressedStates: make([]bool, C.MaxFretCnt+1),
		PluckedState:      false,
		CurValidFret:      0,
	}, nil
}

func (bsm *BassStringModel) PressFret(fret int) {
	if fret < 0 || fret > C.MaxFretCnt {
		return
	}
	bsm.FretPressedStates[fret] = true
	bsm.CurValidFret = max(bsm.CurValidFret, fret)
}

func (bsm *BassStringModel) ReleaseFret(fret int) {
	if fret < 0 || fret > C.MaxFretCnt {
		return
	}
	bsm.FretPressedStates[fret] = false
	if bsm.CurValidFret <= fret {
		for i := fret - 1; i >= 0; i-- {
			if bsm.FretPressedStates[i] {
				bsm.CurValidFret = i
				return
			}
		}
		// If no fret is pressed, set `CurValidFret` to 0 (open string)
		bsm.CurValidFret = 0
	}
}

func (bsm *BassStringModel) GetNoteToPlay() note.Note {
	if bsm.CurValidFret == 0 {
		return bsm.BaseNote
	}
	return bsm.FretToNote[bsm.CurValidFret]
}
