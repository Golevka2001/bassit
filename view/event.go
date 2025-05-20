package view

import (
	"time"

	C "bassit/constant"
	"bassit/util"

	hook "github.com/robotn/gohook"
)

func (v *View) HandleKeyEvent(event hook.Event) {
	pluckedString, ok := C.RawcodeToPluckedString[event.Rawcode]
	if ok {
		if event.Kind == hook.KeyDown {
			v.bassModel.Pluck(pluckedString)

			// Play corresponding sound
			curString := v.bassModel.Strings[pluckedString]
			go v.audioManager.PlayBassNote(curString.GetNoteToPlay())

			// Release the pluck and refresh the view after a short delay
			go func() {
				time.Sleep(util.GetVibDuration())
				v.bassModel.ReleasePluck(pluckedString)
				v.Draw()
			}()
		}
		return
	}

	pressedPos, ok := C.RawcodeToPressedPos[event.Rawcode]
	if ok {
		if event.Kind == hook.KeyDown {
			v.bassModel.Press(pressedPos)
		} else {
			v.bassModel.Release(pressedPos)
		}
		return
	}
}
