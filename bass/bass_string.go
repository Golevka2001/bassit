package bass

import (
	"fmt"

	"github.com/Golevka2001/bassit/config"
	"github.com/Golevka2001/bassit/utils"

	"github.com/go-music-theory/music-theory/note"
)

type bassStringModel struct {
	// baseNote is the open string note
	baseNote note.Note
	// fretIdxToNote maps the fret number to the corresponding note
	fretIdxToNote map[int]note.Note
	// fretPressedStates records which fret is pressed, "0" is meaningless
	fretPressedStates []bool
	// IsVibrating indicates whether the string is vibrating
	IsVibrating bool
	// pluckStates records the state of the two pluck keys
	pluckStates [2]bool
	// CurValidFret records the "valid" fret number among the pressed frets
	// (e.g. if 3rd, 5th, and 7th frets are pressed, `CurValidFret` should be 7)
	CurValidFret int

	LastPluckID int64
}

func newBassString(baseNote note.Note) (*bassStringModel, error) {
	if baseNote.Class == note.Nil || baseNote.Octave <= 0 {
		return nil, fmt.Errorf("invalid base note")
	}

	fretIdxToNote := make(map[int]note.Note)
	fretIdxToNote[0] = baseNote
	for i := 1; i <= config.MaxFretCnt; i++ {
		fretIdxToNote[i] = utils.GetNoteStepFrom(fretIdxToNote[i-1], 1)
	}

	return &bassStringModel{
		baseNote:          baseNote,
		fretIdxToNote:     fretIdxToNote,
		fretPressedStates: make([]bool, config.MaxFretCnt+1),
		IsVibrating:       false,
		CurValidFret:      0,
	}, nil
}

func (bsm *bassStringModel) pressFret(fretIdx int) {
	if fretIdx < 0 || fretIdx > config.MaxFretCnt {
		return
	}
	bsm.fretPressedStates[fretIdx] = true
	bsm.CurValidFret = max(bsm.CurValidFret, fretIdx)
}

func (bsm *bassStringModel) releaseFret(fretIdx int) {
	if fretIdx < 0 || fretIdx > config.MaxFretCnt {
		return
	}
	bsm.fretPressedStates[fretIdx] = false
	if bsm.CurValidFret <= fretIdx {
		for i := fretIdx - 1; i >= 0; i-- {
			if bsm.fretPressedStates[i] {
				bsm.CurValidFret = i
				return
			}
		}
		// If no fret is pressed, set `CurValidFret` to 0 (open string)
		bsm.CurValidFret = 0
	}
}

func (bsm *bassStringModel) pluckString(position int) {
	if position < 0 || position >= 2 {
		return
	}
	bsm.pluckStates[position] = true
	bsm.pluckStates[1-position] = false
	bsm.IsVibrating = true
}

func (bsm *bassStringModel) restoreString() {
	bsm.pluckStates = [2]bool{false, false}
	bsm.IsVibrating = false
}

func (bsm *bassStringModel) GetNoteToPlay() note.Note {
	if bsm.CurValidFret == 0 {
		return bsm.baseNote
	}
	return bsm.fretIdxToNote[bsm.CurValidFret]
}
