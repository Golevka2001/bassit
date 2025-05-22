package view

import (
	"math"

	"bassit/audio"

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

type TextAlign int

const (
	AlignLeft TextAlign = iota
	AlignCenter
	AlignRight
)

// DrawTextLine draws a text line with a given alignment at the specified position.
// Parameters:
// - s: The screen to draw on
// - x1, x2: The horizontal boundary of the text line (inclusive)
// - y: The y-coordinate of the text line
// - text: The text to draw
// - align: AlignLeft|AlignCenter|AlignRight
// - style: The style of the text
func DrawTextLine(
	s *tcell.Screen,
	x1, x2 int,
	y int,
	text string,
	align TextAlign,
	style tcell.Style,
) {
	tLen := len(text)
	if x1 > x2 || y < 0 || tLen < 1 {
		return
	}

	startX := x1
	switch align {
	case AlignCenter:
		startX = int(math.Round(float64(x1+x2-tLen) / 2))
	case AlignRight:
		startX = x2 - tLen + 1
	}

	tOffset := 0
	if startX < x1 {
		tOffset = x1 - startX
		startX = x1
	}

	for i, c := range text[tOffset:] {
		if i+startX > x2 {
			break
		}
		(*s).SetContent(startX+i, y, c, nil, style)
	}
	(*s).Show()
}
