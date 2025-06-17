package bass

import (
	"github.com/Golevka2001/bassit/config"

	"github.com/go-music-theory/music-theory/note"
)

type BassModel struct {
	Strings []*BassStringModel // from highest to lowest
}

func New(tuning [config.StringCnt]string) (*BassModel, error) {
	strings := make([]*BassStringModel, config.StringCnt)
	for i, noteName := range tuning {
		baseNote := *note.Named(noteName)
		bsm, err := newBassString(baseNote)
		if err != nil {
			return nil, err
		}
		strings[i] = bsm
	}

	return &BassModel{
		Strings: strings,
	}, nil
}

func (bm *BassModel) GetBaseNotes() []note.Note {
	notes := make([]note.Note, config.StringCnt)
	for i := range bm.Strings {
		notes[i] = bm.Strings[i].BaseNote
	}
	return notes
}

func (bm *BassModel) GetLowestAndHighestNotes() (lowest, highest note.Note) {
	lowest = bm.Strings[config.StringCnt-1].BaseNote
	highest = bm.Strings[0].FretToNote[config.DisplayedFretCount]
	return
}

func (bm *BassModel) Press(stringIdx, fretIdx int) {
	if stringIdx < 0 || stringIdx >= config.StringCnt {
		return
	}
	bm.Strings[stringIdx].pressFret(fretIdx)
}

func (bm *BassModel) Release(stringIdx, fretIdx int) {
	if stringIdx < 0 || stringIdx >= config.StringCnt {
		return
	}
	bm.Strings[stringIdx].releaseFret(fretIdx)
}

func (bm *BassModel) PluckString(stringIdx int) {
	if stringIdx < 0 || stringIdx >= config.StringCnt {
		return
	}
	bm.Strings[stringIdx].PluckedState = true
}

func (bm *BassModel) RestoreString(stringIdx int) {
	if stringIdx < 0 || stringIdx >= config.StringCnt {
		return
	}
	bm.Strings[stringIdx].PluckedState = false
}
