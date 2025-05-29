package view

import (
	"time"

	C "github.com/Golevka2001/bassit/constant"
	"github.com/Golevka2001/bassit/util"

	hook "github.com/robotn/gohook"
)

func (v *ViewManager) HandleKeyEvent(event hook.Event) {
	am := v.audioManager
	bm := v.bassModel
	bv := v.bassView

	pluckedStringIdx, ok := C.RawcodeToPluckedString[event.Rawcode]
	if ok {
		if event.Kind == hook.KeyDown {
			bm.Pluck(pluckedStringIdx)

			// Draw vibrating string
			bv.drawPluckedString(pluckedStringIdx)

			// Play corresponding sound
			curNote := bm.Strings[pluckedStringIdx].GetNoteToPlay()
			go am.PlayBassNote(curNote)

			// Release the pluck and refresh the view after a short delay
			go func() {
				time.Sleep(util.GetVibDuration())
				bm.ReleasePluck(pluckedStringIdx)
				bv.restorePluckedString(pluckedStringIdx)
				am.StopBassNote(curNote)
			}()
		}
		return
	}

	pressedPos, ok := C.RawcodeToPressedPos[event.Rawcode]
	if ok {
		// Stop vibration
		pressedString := bm.Strings[pressedPos.String]
		if pressedString.PluckedState && pressedPos.Fret >= pressedString.CurValidFret {
			bv.restorePluckedString(pressedPos.String)
			am.StopBassNote(pressedString.GetNoteToPlay())
		}

		if event.Kind == hook.KeyDown {
			bm.Press(pressedPos)

			// Draw pressed fret
			bv.drawPressedFret(pressedPos)
		} else if event.Kind == hook.KeyUp {
			bm.Release(pressedPos)

			// Restore pressed fret
			bv.restorePressedFret(pressedPos)
		}
		return
	}
}

func (vm *ViewManager) HandleWindowResizing() {
	(*vm.screen).Sync()

	vm.titleView.RecalcPositions()
	vm.bassView.RecalcPositions()

	vm.Draw()
}
