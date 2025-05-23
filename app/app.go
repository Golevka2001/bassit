package app

import (
	"fmt"
	"os"
	"runtime"

	"bassit/audio"
	C "bassit/constant"
	"bassit/model"
	"bassit/util"
	"bassit/view"

	"github.com/gdamore/tcell/v2"
	hook "github.com/robotn/gohook"
)

func Run() {
	// Detect OS
	C.OS = runtime.GOOS

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
	s.Clear()
	defer s.Fini()

	bm, err := model.NewBassModel(C.StandardTuning)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to create bass model:", err)
		os.Exit(1)
	}

	am, err := audio.NewAudioManager(bm.GetLowestAndHighestNotes())
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to create audio manager:", err)
		os.Exit(1)
	}

	vm := view.NewViewManager(&s, am, bm)

	// Start event loop
	runEventLoop(&s, &vm)
}

func runEventLoop(s *tcell.Screen, vm *view.ViewManager) {
	vm.Draw()
	// done := false

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.

		// done = true

		maybePanic := recover()
		(*s).Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	// Tcell event is used to handle window resizing
	go func() {
		for {
			tcellEv := (*s).PollEvent()
			switch tcellEv.(type) {
			case *tcell.EventResize:
				vm.HandleWindowResizing()
			}
		}
	}()

	// Gohook event is used to handle key presses
	gohookEvChan := hook.Start()
	defer hook.End()
	for ev := range gohookEvChan {
		if ev.Kind == hook.KeyDown || ev.Kind == hook.KeyUp {
			curKey := C.RawcodeToKey[ev.Rawcode]
			if C.OS == "darwin" {
				curKey = C.RawcodeToKeyForDarwin[ev.Rawcode]
			}
			if curKey == "escape" {
				// Escape key pressed, exit the program
				quit()
				return
			}

			go vm.HandleKeyEvent(ev)
		}
	}
}
