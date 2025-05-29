package app

import (
	"fmt"
	"os"

	"github.com/Golevka2001/bassit/audio"
	"github.com/Golevka2001/bassit/config"
	C "github.com/Golevka2001/bassit/constant"
	"github.com/Golevka2001/bassit/model"
	"github.com/Golevka2001/bassit/util"
	"github.com/Golevka2001/bassit/view"

	"github.com/gdamore/tcell/v2"
	hook "github.com/robotn/gohook"
)

func Run() {
	util.MapRawcodeToPressedPos()
	util.MapRawcodeToPluckedString()

	// Initialize screen
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

	// Initialize bass
	bm, err := model.NewBassModel(config.Config.Tuning)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to create bass model:", err)
		os.Exit(1)
	}

	// Initialize audio manager
	am, err := audio.NewAudioManager()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to create audio manager:", err)
		os.Exit(1)
	}

	// Initialize view manager
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
