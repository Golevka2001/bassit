package view

import (
	"bassit/audio"
	C "bassit/constant"
	"bassit/model"
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

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

func NewViewManager(
	s *tcell.Screen,
	am *audio.AudioManager,
	bm *model.BassModel,
) ViewManager {
	// Do some checks
	cv := NewCheckView(s)
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
