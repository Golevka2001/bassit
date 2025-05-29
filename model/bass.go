package model

import (
	C "github.com/Golevka2001/bassit/constant"

	"github.com/go-music-theory/music-theory/note"
)

type BassModel struct {
	Strings []*BassStringModel // from highest to lowest
}

func NewBassModel(tuning [C.StringCnt]string) (*BassModel, error) {
	strings := make([]*BassStringModel, C.StringCnt)
	for i, noteName := range tuning {
		baseNote := *note.Named(noteName)
		bsm, err := NewBassStringModel(baseNote)
		if err != nil {
			return nil, err
		}
		strings[i] = bsm
	}

	return &BassModel{
		Strings: strings,
	}, nil
}

func (bm *BassModel) GetLowestAndHighestNotes() (lowest, highest note.Note) {
	lowest = bm.Strings[C.StringCnt-1].BaseNote
	highest = bm.Strings[0].FretToNote[C.DisplayedFretNum]
	return
}

func (bm *BassModel) Press(pos C.PressedPos) {
	if pos.String < 0 || pos.String >= C.StringCnt {
		return
	}
	bm.Strings[pos.String].PressFret(pos.Fret)
}

func (bm *BassModel) Release(pos C.PressedPos) {
	if pos.String < 0 || pos.String >= C.StringCnt {
		return
	}
	bm.Strings[pos.String].ReleaseFret(pos.Fret)
}

func (bm *BassModel) Pluck(stringIdx int) {
	if stringIdx < 0 || stringIdx >= C.StringCnt {
		return
	}
	bm.Strings[stringIdx].PluckedState = true
}

func (bm *BassModel) ReleasePluck(stringIdx int) {
	if stringIdx < 0 || stringIdx >= C.StringCnt {
		return
	}
	bm.Strings[stringIdx].PluckedState = false
}
