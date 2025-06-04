package view

import (
	"fmt"
	"os"
	"time"

	"github.com/Golevka2001/bassit/audio"
	C "github.com/Golevka2001/bassit/constant"
	"github.com/Golevka2001/bassit/model"
	U "github.com/Golevka2001/bassit/util"

	"github.com/gdamore/tcell/v2"
)

type BaseView struct {
	screen       *tcell.Screen
	audioManager *audio.AudioManager
	screenW      int
	screenH      int
	startX       int
	endX         int
	startY       int
	endY         int
}

type ViewManager struct {
	screen       *tcell.Screen
	audioManager *audio.AudioManager
	bassModel    *model.BassModel
	width        int
	height       int
	checkView    *CheckView
	titleView    *TitleView
	bassView     *BassView
}

// SkipCheck will be set by the `root` command
var SkipCheck = false

func NewViewManager(
	s *tcell.Screen,
	am *audio.AudioManager,
	bm *model.BassModel,
) ViewManager {
	var cv *CheckView
	if !SkipCheck {
		// Do some checks
		cv = NewCheckView(s)
		failed := cv.RunChecks()
		cv.Fini()

		// Quit if any check failed
		if len(failed) > 0 {
			(*s).Fini()

			fmt.Println("❌ Exiting due to failed checks:")
			for _, fail := range failed {
				fmt.Printf(" - %s\n", fail)
			}
			os.Exit(1)
		}
	}

	// Draw welcome screen if not skipped
	U.DrawWelcome(s)
	startTime := time.Now()

	// All checks passed. Then generate and load audio resources
	lowestNote, highestNote := bm.GetLowestAndHighestNotes()
	if err := U.GenAllPossibleNotes(lowestNote, highestNote); err != nil {
		(*s).Fini()
		fmt.Printf("❌ Failed to generate audio resources: %v\n", err)
		os.Exit(1)
	}
	am.LoadNoteSamples(lowestNote, highestNote)

	elapsed := time.Since(startTime)
	if elapsed < time.Duration(C.MinWelcomeDuration)*time.Millisecond {
		time.Sleep(time.Duration(C.MinWelcomeDuration)*time.Millisecond - elapsed)
	}

	w, h := (*s).Size()

	// Initialize other views
	tv, tvEndX, tvEndY := NewTitleView(s, nil, 0, 0)
	bv, bvEndX, bvEndY := NewBassView(s, am, bm, 0, tvEndY+C.BassViewMarginTop)

	_ = tvEndX
	_ = bvEndX
	_ = bvEndY

	return ViewManager{
		screen:       s,
		bassModel:    bm,
		audioManager: am,
		width:        w,
		height:       h,
		checkView:    cv,
		titleView:    tv,
		bassView:     bv,
	}
}

func (vm *ViewManager) Draw() {
	(*vm.screen).Clear()

	vm.titleView.Draw()
	vm.bassView.Draw()
	(*vm.screen).Show()
}
