package view

import (
	_ "math"

	"bassit/audio"
	C "bassit/constant"
	"bassit/model"
	_ "bassit/util"

	"github.com/gdamore/tcell/v2"
)

type ViewManager struct {
	screen       *tcell.Screen
	audioManager *audio.AudioManager
	bassModel    *model.BassModel
	width        int
	height       int
	titleView    *TitleView
	bassView     *BassView
}

func NewViewManager(
	s *tcell.Screen,
	am *audio.AudioManager,
	bm *model.BassModel,
) ViewManager {
	(*s).SetStyle(tcell.StyleDefault)
	(*s).Clear()

	w, h := (*s).Size()

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
		titleView:    tv,
		bassView:     bv,
	}
}

func (v *ViewManager) Draw() {
	(*v.screen).Clear()

	v.titleView.Draw()
	v.bassView.Draw()
	(*v.screen).Show()
}
