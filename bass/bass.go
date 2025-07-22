package bass

import (
	"github.com/Golevka2001/bassit/config"

	"github.com/go-music-theory/music-theory/note"
)

type BassModel struct {
	strings []*bassStringModel // from highest to lowest
}

func NewBass(tuning [config.StringCnt]string) (*BassModel, error) {
	strings := make([]*bassStringModel, config.StringCnt)
	for i, noteName := range tuning {
		baseNote := *note.Named(noteName)
		bsm, err := NewBassString(baseNote)
		if err != nil {
			return nil, err
		}
		strings[i] = bsm
	}

	return &BassModel{
		strings: strings,
	}, nil
}

// GetBaseNotes returns the open string notes of each string
func (bm *BassModel) GetBaseNotes() []note.Note {
	notes := make([]note.Note, config.StringCnt)
	for i := range bm.strings {
		notes[i] = *bm.strings[i].GetNoteAt(0)
	}
	return notes
}

// GetRange returns the lowest and highest notes according to the displayed fret count
func (bm *BassModel) GetRange() (lowest, highest note.Note) {
	lowest = *bm.strings[config.StringCnt-1].GetNoteAt(0)
	highest = *bm.strings[0].GetNoteAt(config.DisplayedFretCount)
	return
}

// GetActualRange returns the lowest and highest notes according to total fret count
func (bm *BassModel) GetActualRange() (lowest, highest note.Note) {
	lowest = *bm.strings[config.StringCnt-1].GetNoteAt(0)
	highest = *bm.strings[0].GetNoteAt(config.MaxFretCnt)
	return
}

func (bm *BassModel) IsStringVibrating(stringIdx int) bool {
	if stringIdx < 0 || stringIdx >= config.StringCnt {
		return false
	}
	return bm.strings[stringIdx].isVibrating
}

func (bm *BassModel) GetNoteToPlay(stringIdx int) *note.Note {
	if stringIdx < 0 || stringIdx >= config.StringCnt {
		return nil
	}
	return bm.strings[stringIdx].GetNoteToPlay()
}

func (bm *BassModel) IsFretPressed(pos FretboardPosition) bool {
	if pos.StringIdx < 0 || pos.StringIdx >= config.StringCnt {
		return false
	}
	return bm.strings[pos.StringIdx].IsFretPressed(pos.FretIdx)
}

func (bm *BassModel) GetRightmostPressedFretIdxOfString(stringIdx int) int {
	if stringIdx < 0 || stringIdx >= config.StringCnt {
		return -1
	}
	return bm.strings[stringIdx].GetRightmostPressedFret()
}

func (bm *BassModel) GetNoteAt(pos FretboardPosition) *note.Note {
	if pos.StringIdx < 0 || pos.StringIdx >= config.StringCnt {
		return nil
	}
	return bm.strings[pos.StringIdx].GetNoteAt(pos.FretIdx)
}

func (bm *BassModel) PressFret(pos FretboardPosition) {
	if pos.StringIdx < 0 || pos.StringIdx >= config.StringCnt {
		return
	}
	bm.strings[pos.StringIdx].PressFret(pos.FretIdx)
}

func (bm *BassModel) ReleaseFret(pos FretboardPosition) {
	if pos.StringIdx < 0 || pos.StringIdx >= config.StringCnt {
		return
	}
	bm.strings[pos.StringIdx].ReleaseFret(pos.FretIdx)
}

// PluckString plucks the string and returns the time of the last pluck
func (bm *BassModel) PluckString(stringIdx int, t PluckType) int64 {
	if stringIdx < 0 || stringIdx >= config.StringCnt {
		return 0
	}
	return bm.strings[stringIdx].Pluck(t)
}

// StopVibratingString stops the string from vibrating if the time of the last pluck is the same as the parameter `t`
func (bm *BassModel) StopVibratingString(stringIdx int, t int64) bool {
	if stringIdx < 0 || stringIdx >= config.StringCnt {
		return false
	}
	return bm.strings[stringIdx].StopVibrating(t)
}

func (bm *BassModel) StopVibratingStringWithoutCheck(stringIdx int) {
	if stringIdx < 0 || stringIdx >= config.StringCnt {
		return
	}
	bm.strings[stringIdx].StopVibratingWithoutCheck()
}
