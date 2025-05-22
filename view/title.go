package view

import (
	"bassit/audio"
	C "bassit/constant"
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type TitleView struct {
	BaseView
}

func NewTitleView(
	s *tcell.Screen,
	am *audio.AudioManager,
	startX int,
	startY int,
) (*TitleView, int, int) {
	w, h := (*s).Size()

	return &TitleView{
		BaseView: BaseView{
			screen:       s,
			audioManager: am,
			screenW:      w,
			screenH:      h,
			startX:       startX,
			endX:         w,
			startY:       startY,
			endY:         startY + C.TitleViewHeight,
		},
	}, w, startY + C.TitleViewHeight
}

func (tv *TitleView) Draw() {
	s := *tv.screen
	y := (tv.startY + tv.endY) / 2
	style := tcell.StyleDefault.Foreground(C.TitleFgColor).Background(C.TitleBgColor)

	// Background
	for r := tv.startY; r < tv.endY; r++ {
		for c := tv.startX; c < tv.endX; c++ {
			s.SetContent(c, r, ' ', nil, style)
		}
	}

	// Title
	title := fmt.Sprintf("BASSIT v%s", C.Version)
	DrawTextLine(tv.screen, 0, tv.screenW-1, y, title, AlignCenter, style)
}

// RecalcPositions should be called when the window is resized
func (tv *TitleView) RecalcPositions() {
	w, h := (*tv.screen).Size()

	tv.screenW = w
	tv.screenH = h
	tv.endX = w
}
