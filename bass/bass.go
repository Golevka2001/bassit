package bass

import (
	"time"

	"github.com/Golevka2001/bassit/config"

	"github.com/go-music-theory/music-theory/note"
)

type BassModel struct {
	Strings []*bassStringModel // from highest to lowest
}

func New(tuning [config.StringCnt]string) (*BassModel, error) {
	strings := make([]*bassStringModel, config.StringCnt)
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
		notes[i] = bm.Strings[i].baseNote
	}
	return notes
}

func (bm *BassModel) GetLowestAndHighestNotes() (lowest, highest note.Note) {
	lowest = bm.Strings[config.StringCnt-1].baseNote
	highest = bm.Strings[0].fretIdxToNote[config.DisplayedFretCount]
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

func (bm *BassModel) PluckString(stringIdx, position int) {
	if stringIdx < 0 || stringIdx >= config.StringCnt || position < 0 || position >= 2 {
		return
	}
	bm.Strings[stringIdx].pluckString(position)
	bm.Strings[stringIdx].LastPluckID = time.Now().UnixNano()
}

func (bm *BassModel) RestoreString(stringIdx int) {
	if stringIdx < 0 || stringIdx >= config.StringCnt {
		return
	}
	bm.Strings[stringIdx].restoreString()
}
