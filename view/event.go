package view

import (
	"time"

	C "bassit/constant"
	"bassit/util"

	hook "github.com/robotn/gohook"
)

func (v *View) HandleKeyEvent(event hook.Event) {
	pluckedStringIdx, ok := C.RawcodeToPluckedString[event.Rawcode]
	if ok {
		if event.Kind == hook.KeyDown {
			v.bassModel.Pluck(pluckedStringIdx)

			// Draw vibrating string
			v.drawPluckedString(pluckedStringIdx)

			// Play corresponding sound
			curString := v.bassModel.Strings[pluckedStringIdx]
			go v.audioManager.PlayBassNote(curString.GetNoteToPlay())

			// Release the pluck and refresh the view after a short delay
			go func() {
				time.Sleep(util.GetVibDuration())
				v.bassModel.ReleasePluck(pluckedStringIdx)
				v.restorePluckedString(pluckedStringIdx)
			}()
		}
		return
	}

	pressedPos, ok := C.RawcodeToPressedPos[event.Rawcode]
	if ok {
		// Stop vibration
		pressedString := v.bassModel.Strings[pressedPos.String]
		if pressedString.PluckedState && pressedPos.Fret >= pressedString.CurValidFret {
			v.restorePluckedString(pressedPos.String)
		}

		if event.Kind == hook.KeyDown {
			v.bassModel.Press(pressedPos)

			// Draw pressed fret
			v.drawPressedFret(pressedPos)
		} else if event.Kind == hook.KeyUp {
			v.bassModel.Release(pressedPos)

			// Restore pressed fret
			v.restorePressedFret(pressedPos)
		}
		return
	}
}
