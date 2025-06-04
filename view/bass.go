package view

import (
	"math"

	"github.com/Golevka2001/bassit/audio"
	"github.com/Golevka2001/bassit/config"
	C "github.com/Golevka2001/bassit/constant"
	"github.com/Golevka2001/bassit/model"
	U "github.com/Golevka2001/bassit/util"

	"github.com/gdamore/tcell/v2"
)

var (
	fretboardStartX  int
	fretboardEndX    int
	fretboardLen     int
	blockInlayStartY int
	blockInlayEndY   int

	fretWireToX  []int
	xToFretWire  map[int]int
	posMarkerToX []int
	stringToY    []int
	yToString    map[int]int
)

var t = &config.Theme

type BassView struct {
	BaseView
	bassModel *model.BassModel
}

func NewBassView(
	s *tcell.Screen,
	am *audio.AudioManager,
	bm *model.BassModel,
	startX int,
	startY int,
) (*BassView, int, int) {
	w, h := (*s).Size()

	bvH := 2*C.StringMarginY + (C.StringCnt-1)*C.StringSpacing
	calcAuxValues(startX, w, startY, startY+bvH)

	return &BassView{
		BaseView: BaseView{
			screen:       s,
			audioManager: am,
			screenW:      w,
			screenH:      h,
			startX:       startX,
			endX:         w,
			startY:       startY,
			endY:         startY + bvH,
		},
		bassModel: bm,
	}, w, startY + bvH
}

// Draw draws the fretboard and strings.
// It should only be called when the view is switched or the size changes.
func (bv *BassView) Draw() {
	bv.drawFretboard()
	bv.drawStrings()
}

// RecalcPositions should be called when the window is resized
func (bv *BassView) RecalcPositions() {
	w, h := (*bv.screen).Size()

	bv.screenW = w
	bv.screenH = h
	bv.endX = w

	calcAuxValues(bv.startX, bv.endX, bv.startY, bv.endY)
}

func (bv *BassView) drawFretboard() {
	s := *bv.screen

	// Draw the fret board, nut and fret wires
	for x := fretboardStartX; x <= fretboardEndX; x++ {
		for y := bv.startY; y <= bv.endY; y++ {
			charToDraw := ' '
			style := tcell.StyleDefault

			if x == fretboardStartX {
				// The left side of the nut
				if y == bv.startY {
					// The upper left corner
					charToDraw = config.Theme.NutULCornerChar
				} else if y == bv.endY {
					// The lower left corner
					charToDraw = config.Theme.NutLLCornerChar
				} else {
					// The vertical border
					charToDraw = config.Theme.NutVBorderChar
				}
				style = style.Foreground(config.Theme.NutBorderColor)
			} else if x > fretboardStartX && x < fretWireToX[0] {
				// The nut
				if y == bv.startY || y == bv.endY {
					// The horizontal border
					charToDraw = config.Theme.NutHBorderChar
					style = style.Foreground(config.Theme.NutBorderColor)
				} else {
					charToDraw = ' '
					style = style.Background(config.Theme.NutBgColor)
				}
			} else if x == fretWireToX[0] {
				// The right side of the nut
				style = style.Foreground(t.NutBorderColor)
				if t.NutBgColor == tcell.ColorDefault {
					if y == bv.startY {
						// The upper right corner
						charToDraw = t.NutURCornerChar
					} else if y == bv.endY {
						// The lower right corner
						charToDraw = t.NutLRCornerChar
					} else {
						charToDraw = t.NutVBorderChar
						style = style.Background(t.NutBgColor)
					}
				} else {
					if y == bv.startY || y == bv.endY {
						// The horizontal border
						charToDraw = t.NutHBorderChar
					} else {
						style = style.Background(t.NutBgColor)
					}
				}
			} else if fretWireIdx, ok := xToFretWire[x]; ok && fretWireIdx != 0 {
				// Fret wires
				if y == bv.startY {
					// Fret wire at the upper border
					charToDraw = t.FretWireUpperChar
				} else if y == bv.endY {
					// Fret wire at the lower border
					charToDraw = t.FretWireLowerChar
				} else {
					charToDraw = t.FretWireChar
				}
				style = style.Foreground(t.FretWireColor)
			} else {
				if y == bv.startY || y == bv.endY {
					// The horizontal border
					charToDraw = t.FretboardHBorderChar
					style = style.Foreground(t.FretboardBorderColor)
				}
			}

			if x > fretWireToX[0] && x <= fretboardEndX && y > bv.startY && y < bv.endY {
				style = style.Background(t.FretboardBgColor)
			}

			s.SetContent(x, y, charToDraw, nil, style)
		}
	}

	// Draw fretboard inlays (or position markers)
	switch t.InlayShape {
	case config.DotInlayShape:
		for fretIdx, x := range posMarkerToX {
			if fretIdx > C.DisplayedFretNum {
				break
			}
			charToDraw := string(t.InlayChar)
			style := tcell.StyleDefault.Foreground(t.InlayColor).Background(t.FretboardBgColor)
			switch fretIdx {
			case 3, 5, 7, 9, 15, 17, 19, 21:
				y := int(math.Round(float64(stringToY[1]+stringToY[2]) / 2)) // TODO: `StringNum` is regarded as 4
				U.DrawTextLine(bv.screen, x-1, x+1, y, charToDraw, U.AlignCenter, style)
			case 12, 24:
				y1 := int(math.Round(float64(stringToY[0]+stringToY[1]) / 2))
				y2 := int(math.Round(float64(stringToY[2]+stringToY[3]) / 2))
				U.DrawTextLine(bv.screen, x-1, x+1, y1, charToDraw, U.AlignCenter, style)
				U.DrawTextLine(bv.screen, x-1, x+1, y2, charToDraw, U.AlignCenter, style)
			}
		}
	case config.BlockInlayShape:
		for _, fretIdx := range []int{1, 3, 5, 7, 9, 12, 15, 17, 19, 21, 24} {
			if fretIdx > C.DisplayedFretNum {
				break
			}
			fretWidth := fretWireToX[fretIdx] - fretWireToX[fretIdx-1]
			margin := int(math.Round(float64(fretWidth) / 3))
			margin = max(margin, 2)
			x1 := fretWireToX[fretIdx-1] + margin
			x2 := fretWireToX[fretIdx] - margin
			style := tcell.StyleDefault.Background(t.InlayColor)
			U.FillArea(bv.screen, x1, x2, blockInlayStartY, blockInlayEndY, true, ' ', style)
		}
	}

	s.Show()
}

func (bv *BassView) drawStrings() {
	s := *bv.screen
	b := bv.bassModel

	for stringIdx, y := range stringToY {
		curString := b.Strings[stringIdx]

		// Draw base note names
		x := bv.startX + C.StringBaseNoteNameMarginLeft
		noteName := U.GetNoteNameWithOctave(curString.BaseNote)
		noteNameStyle := tcell.StyleDefault.Foreground(t.BaseNoteNameFgColor).Background(t.BaseNoteNameBgColor)
		U.DrawTextLine(bv.screen, x, x+2, y, noteName, U.AlignLeft, noteNameStyle)

		// Draw string lines
		bv.restorePluckedString(stringIdx)
		// Draw the string over the left side of the nut
		_, _, origStyle, _ := s.GetContent(fretboardStartX, y)
		style := origStyle.Foreground(t.StringColor)
		s.SetContent(fretboardStartX, y, t.StringOverBoarderChar, nil, style)
	}
	s.Show()
}

func (bv *BassView) drawPressedFret(pressedPos C.PressedPos) {
	s := *bv.screen

	x := calcFretCenterXPos(pressedPos.Fret)
	y := stringToY[pressedPos.String]

	_, _, origStyle, _ := s.GetContent(x, y)
	style := origStyle.Foreground(t.PressedFretSignColor)

	U.DrawTextLine(bv.screen, x-1, x+1, y, string(t.PressedFretSignChar), U.AlignCenter, style)
	s.Show()
}

func (bv *BassView) restorePressedFret(pressedPos C.PressedPos) {
	s := *bv.screen

	x := calcFretCenterXPos(pressedPos.Fret)
	y := stringToY[pressedPos.String]

	_, _, origStyle, _ := s.GetContent(x, y)
	style := origStyle.Foreground(t.StringColor)

	s.SetContent(x, y, t.StringChar, nil, style)
	s.Show()
}

func (bv *BassView) drawPluckedString(stringIdx int) {
	s := *bv.screen

	// Draw the sign
	x := bv.screenW - 1 - C.PluckedStringSignMarginRight
	y := stringToY[stringIdx]
	style := tcell.StyleDefault.Foreground(t.PluckedStringSignColor)
	U.DrawTextLine(bv.screen, x-1, x+1, y, string(t.PluckedStringSignChar), U.AlignCenter, style)

	// Draw the vibrating string
	curString := bv.bassModel.Strings[stringIdx]
	rightMostPressedPos := calcFretCenterXPos(curString.CurValidFret)
	if curString.CurValidFret == 0 {
		rightMostPressedPos = fretboardStartX + C.NutWidth
	}
	if rightMostPressedPos == -1 {
		return
	}
	for x := rightMostPressedPos + 1; x <= fretboardEndX; x++ {
		charToDraw := t.VibratingStringChar
		if x == fretboardStartX+C.NutWidth {
			charToDraw = t.VibratingStringOverBoarderChar
		} else if _, ok := xToFretWire[x]; ok {
			charToDraw = t.VibratingStringOverFretWireChar
		}

		_, _, origStyle, _ := s.GetContent(x, y)
		style := origStyle.Foreground(t.StringColor)
		s.SetContent(x, y, charToDraw, nil, style)
	}
	s.Show()
}

func (bv *BassView) restorePluckedString(stringIdx int) {
	s := *bv.screen

	// Restore the sign
	x := bv.screenW - 1 - C.PluckedStringSignMarginRight
	y := stringToY[stringIdx]
	style := tcell.StyleDefault
	U.DrawTextLine(bv.screen, x-1, x+1, y, string(t.NotPluckedStringChar), U.AlignCenter, style)

	// Restore the vibrating string
	curString := bv.bassModel.Strings[stringIdx]
	rightMostPressedPos := calcFretCenterXPos(curString.CurValidFret)
	if curString.CurValidFret == 0 {
		rightMostPressedPos = fretboardStartX
	}
	if rightMostPressedPos == -1 {
		return
	}
	for x := rightMostPressedPos + 1; x <= fretboardEndX; x++ {
		charToDraw := t.StringChar
		if x == fretboardStartX+C.NutWidth {
			if t.NutBgColor == tcell.ColorDefault {
				charToDraw = t.StringOverBoarderChar
			} else {
				charToDraw = t.StringChar
			}
		} else if _, ok := xToFretWire[x]; ok {
			charToDraw = t.StringOverFretWireChar
		}

		_, _, origStyle, _ := s.GetContent(x, y)
		style := origStyle.Foreground(t.StringColor)
		s.SetContent(x, y, charToDraw, nil, style)
	}
	s.Show()
}

// calcAuxValues calculates the auxiliary values for drawing the fretboard and strings
// Parameters:
//   - x1, x2: The starting and ending x-coordinates of the fretboard
//   - y1, y2: The starting and ending y-coordinates of the fretboard
func calcAuxValues(x1, x2, y1, y2 int) {
	fretboardStartX = x1 + C.FretboardMarginX
	fretboardEndX = x2 - C.FretboardMarginX
	fretboardLen = fretboardEndX - fretboardStartX
	blockInlayStartY = y1 + C.BlockInlayMarginY
	blockInlayEndY = y2 - C.BlockInlayMarginY

	fretWireToX = make([]int, C.DisplayedFretNum+1)
	xToFretWire = make(map[int]int)
	posMarkerToX = make([]int, C.DisplayedFretNum+1)
	stringToY = make([]int, C.StringCnt)
	yToString = make(map[int]int)

	scaleFactor := float64(fretboardLen-C.NutWidth) / calcDn(C.DisplayedFretNum)
	for fretIdx := 0; fretIdx <= C.DisplayedFretNum; fretIdx++ {
		// Calculate the position of fret wires
		x := calcFretWireXPos(fretIdx, scaleFactor, fretboardStartX+C.NutWidth)
		fretWireToX[fretIdx] = x
		xToFretWire[x] = fretIdx

		// Calculate the position of position markers
		switch fretIdx {
		case 3, 5, 7, 9, 12, 15, 17, 19, 21, 24:
			posMarkerToX[fretIdx] = calcFretCenterXPos(fretIdx)
		}
	}

	// Calculate the position of strings
	for stringIdx := range C.StringCnt {
		y := y1 + C.StringMarginY
		if stringIdx > 0 {
			y = stringToY[stringIdx-1] + C.StringSpacing
		}
		stringToY[stringIdx] = y
		yToString[y] = stringIdx
	}
}

// calcFretCenterXPos calculates the x-coordinate of the center of the nth fret
// Parameters:
//   - n: The fret number (0 means the nut)
func calcFretCenterXPos(n int) int {
	if n < 1 || n > C.DisplayedFretNum {
		return -1
	}

	curFretWireX := fretWireToX[n]
	prevFretWireX := fretWireToX[n-1]

	return int(math.Round(float64(curFretWireX+prevFretWireX) / 2))
}

// calcFretWireXPos calculates the x-coordinate of the nth fret wire
// Parameters:
//   - n: The fret number (0 means the nut)
//   - scaleFactor: fretboardLength / d(fretCnt)
//   - offset: The starting position of the fretboard (x-coordinate)
func calcFretWireXPos(n int, scaleFactor float64, offset int) int {
	dist := int(math.Round(calcDn(n) * scaleFactor))
	return int(offset + dist)
}

// calcDn calculates the distance from the nut to the nth fret wire
func calcDn(n int) float64 {
	// Formula: d(n) = SL * (1 - (2^(-n/12)))
	// where `n` is the fret number,
	// `d(n)` is the distance from the nut to the nth fret wire,
	// `SL` is the scale length (not used here)
	// (Reference: https://www.thekimerers.net/brian/YAFCalc/YAFCalc.html)
	return 1 - math.Pow(float64(2), float64(-n)/12)
}
