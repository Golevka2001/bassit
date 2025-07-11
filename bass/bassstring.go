package bass

import (
	"fmt"
	"time"

	"github.com/Golevka2001/bassit/config"
	"github.com/Golevka2001/bassit/utils"

	"github.com/go-music-theory/music-theory/note"
)

type bassStringModel struct {
	// fretIdxToNote maps the fret number to the corresponding note
	fretIdxToNote map[int]note.Note
	// fretPressed records which fret is pressed, "0" is meaningless
	fretPressed [config.MaxFretCnt + 1]bool
	// rightmostPressedFret stores the highest-numbered (rightmost) pressed fret,
	// and is updated whenever a fret is pressed or released
	rightmostPressedFret int
	// isVibrating indicates whether the string is vibrating
	isVibrating bool
	// pluckType records the type of last pluck
	pluckType PluckType
	// lastPluckTime records the time of last pluck
	lastPluckTime int64
}

func NewBassString(baseNote note.Note) (*bassStringModel, error) {
	if baseNote.Class == note.Nil || baseNote.Octave <= 0 {
		return nil, fmt.Errorf("invalid base note")
	}

	fretIdxToNote := make(map[int]note.Note)
	fretIdxToNote[0] = baseNote
	for i := 1; i <= config.MaxFretCnt; i++ {
		fretIdxToNote[i] = utils.GetNoteStepFrom(fretIdxToNote[i-1], 1)
	}

	return &bassStringModel{
		fretIdxToNote:        fretIdxToNote,
		fretPressed:          [config.MaxFretCnt + 1]bool{},
		rightmostPressedFret: 0,
	}, nil
}

// GetNoteAt returns the note at the given fret index
// If the fret index is out of range, it returns nil
func (bsm *bassStringModel) GetNoteAt(fretIdx int) *note.Note {
	n, ok := bsm.fretIdxToNote[fretIdx]
	if !ok {
		return nil
	}
	return &n
}

// GetNoteToPlay returns the note corresponding to the rightmost pressed fret
// If no fret is pressed, it returns the base note
func (bsm *bassStringModel) GetNoteToPlay() *note.Note {
	return bsm.GetNoteAt(bsm.rightmostPressedFret)
}

func (bsm *bassStringModel) PressFret(fretIdx int) {
	if fretIdx < 0 || fretIdx > config.MaxFretCnt {
		return
	}
	bsm.fretPressed[fretIdx] = true
	bsm.rightmostPressedFret = max(bsm.rightmostPressedFret, fretIdx)
}

func (bsm *bassStringModel) ReleaseFret(fretIdx int) {
	if fretIdx < 0 || fretIdx > config.MaxFretCnt {
		return
	}
	bsm.fretPressed[fretIdx] = false
	if bsm.rightmostPressedFret <= fretIdx {
		for i := fretIdx - 1; i >= 0; i-- {
			if bsm.fretPressed[i] {
				bsm.rightmostPressedFret = i
				return
			}
		}
		// If no fret is pressed, set `rightmostPressedFret` to 0 (open string)
		bsm.rightmostPressedFret = 0
	}
}

// Pluck marks the string as vibrating and returns the time
func (bsm *bassStringModel) Pluck(t PluckType) int64 {
	bsm.isVibrating = true
	bsm.pluckType = t
	bsm.lastPluckTime = time.Now().UnixNano()
	return bsm.lastPluckTime
}

// StopVibrating stops the string from vibrating if the time of the last pluck is the same as the parameter `t`
func (bsm *bassStringModel) StopVibrating(t int64) bool {
	if bsm.lastPluckTime != t {
		return false
	}
	bsm.isVibrating = false
	return true
}

func (bsm *bassStringModel) StopVibratingWithoutCheck() {
	bsm.isVibrating = false
}
