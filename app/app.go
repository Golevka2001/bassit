package app

import (
	"fmt"
	"os"

	C "bassit/constant"
	"bassit/model"
	"bassit/util"
	"bassit/view"

	"github.com/gdamore/tcell/v2"
	hook "github.com/robotn/gohook"
	"github.com/shirou/gopsutil/host"
)

func Run() {
	// Get OS information
	platform, _, _, _ := host.PlatformInformation()
	C.OS = platform

	util.MapRawcodeToPressedPos()
	util.MapRawcodeToPluckedString()

	// Initialize
	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create screen:", err)
		os.Exit(1)
	}

	if err = s.Init(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to initialize screen:", err)
		os.Exit(1)
	}

	bm, err := model.NewBassModel(C.StandardTuning)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to create bass model:", err)
		return
	}

	v := view.NewView(s, bm)

	// Start event loop
	runEventLoop(s, v, bm)
}

func runEventLoop(
	s tcell.Screen,
	v view.View,
	bm *model.BassModel,
) {
	// done := false

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.

		// done = true

		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	evChan := hook.Start()
	defer hook.End()

	v.Draw()
	for ev := range evChan {
		if ev.Kind == hook.KeyDown || ev.Kind == hook.KeyUp {
			if C.RawcodeToKey[ev.Rawcode] == "escape" {
				// Escape key pressed, exit the program
				quit()
				return
			}

			v.HandleKeyEvent(ev)

			v.Draw()
		}
	}
}
