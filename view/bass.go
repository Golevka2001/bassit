package view

import (
	"bassit/audio"
	C "bassit/constant"
	"bassit/model"
	"bassit/util"
	"math"

	"github.com/gdamore/tcell/v2"
)

var (
	fretboardStartX int
	fretboardEndX   int
	fretboardLen    int
	fretWireToX     []int
	xToFretWire     map[int]int
	posMarkerToX    []int
	stringToY       []int
	yToString       map[int]int
)

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

	bvH := 2*C.StringMarginTop + (C.StringCnt-1)*C.StringSpacing
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

func (bv *BassView) drawFretboard() {
	s := *bv.screen

	// Draw the fret board, nut and fret wires
	for x := fretboardStartX; x <= fretboardEndX; x++ {
		for y := bv.startY; y <= bv.endY; y++ {
			charToDraw := ' '
			style := tcell.StyleDefault.Foreground(C.FretboardBorderColor).Background(C.FretboardBgColor)

			if x >= bv.startX && x < fretWireToX[0] {
				// The nut
				if y == bv.startY {
					// The upper left corner
					charToDraw = C.FretboardULCornerChar
				} else if y == bv.endY {
					// The lower left corner
					charToDraw = C.FretboardLLCornerChar
				} else {
					// The vertical border
					charToDraw = C.FretboardVBorderChar
				}

				style = tcell.StyleDefault.Foreground(C.NutBorderColor).Background(C.NutBgColor)
			} else if fretWireIdx, ok := xToFretWire[x]; ok {
				// The fret wire
				if y == bv.startY {
					// Fret wire at the upper border
					charToDraw = C.FretWireUpperChar
				} else if y == bv.endY {
					// Fret wire at the lower border
					charToDraw = C.FretWireLowerChar
				} else {
					// The fret wire
					charToDraw = C.FretWireChar
				}

				if fretWireIdx == 0 {
					// The right side of the nut
					style = tcell.StyleDefault.Foreground(C.NutBorderColor).Background(C.NutBgColor)
				} else {
					style = tcell.StyleDefault.Foreground(C.FretWireColor).Background(C.FretboardBgColor)
				}
			} else {
				if y == bv.startY || y == bv.endY {
					// The horizontal border
					charToDraw = C.FretboardHBorderChar
				}
			}

			s.SetContent(x, y, charToDraw, nil, style)
		}
	}

	// Draw the position markers
	style := tcell.StyleDefault.Foreground(C.PosMarkerColor).Background(C.FretboardBgColor)
	for posMarkerIdx, x := range posMarkerToX {
		if posMarkerIdx > C.DisplayedFretNum {
			break
		}

		switch posMarkerIdx {
		case 3, 5, 7, 9, 15, 17, 19, 21:
			y := int(math.Round(float64(stringToY[1]+stringToY[2]) / 2)) // TODO: `StringNum` is regarded as 4
			s.SetContent(x, y, C.PosMarkerChar, nil, style)
		case 12, 24:
			y1 := int(math.Round(float64(stringToY[0]+stringToY[1]) / 2))
			y2 := int(math.Round(float64(stringToY[2]+stringToY[3]) / 2))
			s.SetContent(x, y1, C.PosMarkerChar, nil, style)
			s.SetContent(x, y2, C.PosMarkerChar, nil, style)
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
		noteName := util.GetNoteNameWithOctave(curString.BaseNote)
		noteNameStyle := tcell.StyleDefault.Foreground(C.BaseNoteNameFgColor).Background(C.BaseNoteNameBgColor)
		DrawTextLine(bv.screen, x, x+2, y, noteName, AlignLeft, noteNameStyle)

		// Draw string lines
		style := tcell.StyleDefault.Foreground(C.StringColor).Background(C.NutBgColor)
		s.SetContent(fretboardStartX, y, C.StringOverFretChar, nil, style)
		bv.restorePluckedString(stringIdx)
	}
	s.Show()
}

func (bv *BassView) drawPressedFret(pressedPos C.PressedPos) {
	s := *bv.screen

	x := calcFretCenterXPos(pressedPos.Fret)
	y := stringToY[pressedPos.String]
	style := tcell.StyleDefault.Foreground(C.PressedFretSignColor).Background(C.FretboardBgColor)

	s.SetContent(x, y, C.PressedFretSignChar, nil, style)
	s.Show()
}

func (bv *BassView) restorePressedFret(pressedPos C.PressedPos) {
	s := *bv.screen

	x := calcFretCenterXPos(pressedPos.Fret)
	y := stringToY[pressedPos.String]
	style := tcell.StyleDefault.Foreground(C.StringColor).Background(C.FretboardBgColor)

	s.SetContent(x, y, C.StringChar, nil, style)
	s.Show()
}

func (bv *BassView) drawPluckedString(stringIdx int) {
	s := *bv.screen

	// Draw the sign
	x := bv.screenW - 1 - C.PluckedStringSignMarginRight
	y := stringToY[stringIdx]
	style := tcell.StyleDefault.Foreground(C.PluckedStringSignColor)
	s.SetContent(x, y, C.PluckedStringSignChar, nil, style)

	// Draw the vibrating string
	curString := bv.bassModel.Strings[stringIdx]
	rightMostPressedPos := calcFretCenterXPos(curString.CurValidFret)
	if curString.CurValidFret == 0 {
		rightMostPressedPos = fretboardStartX
	}
	if rightMostPressedPos == -1 {
		return
	}
	for x := rightMostPressedPos + 1; x <= fretboardEndX; x++ {
		charToDraw := C.VibratingStringChar
		if _, ok := xToFretWire[x]; ok || x == fretboardStartX {
			charToDraw = C.VibratingStringOverFretChar
		}

		_, _, originStyle, _ := s.GetContent(x, y)
		style := originStyle.Foreground(C.StringColor)
		s.SetContent(x, y, charToDraw, nil, style)
	}
	s.Show()
}

func (bv *BassView) restorePluckedString(stringIdx int) {
	s := *bv.screen

	// Restore the sign
	x := bv.screenW - 1 - C.PluckedStringSignMarginRight
	y := stringToY[stringIdx]
	style := tcell.StyleDefault.Foreground(C.StringColor).Background(C.FretboardBgColor)
	s.SetContent(x, y, C.NotPluckedStringChar, nil, style)

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
		charToDraw := C.StringChar
		if _, ok := xToFretWire[x]; ok || x == fretboardStartX {
			charToDraw = C.StringOverFretChar
		}

		_, _, originStyle, _ := s.GetContent(x, y)
		style := originStyle.Foreground(C.StringColor)
		s.SetContent(x, y, charToDraw, nil, style)
	}
	s.Show()
}

// calcAuxValues calculates the auxiliary values for drawing the fretboard and strings
// Parameters:
//   - x1, x2: The starting and ending x-coordinates of the fretboard
//   - y1, y2: The starting and ending y-coordinates of the fretboard
func calcAuxValues(x1, x2, y1, y2 int) {
	fretboardStartX = x1 + C.FretboardMarginLeft
	fretboardEndX = x2 - C.FretboardMarginRight
	fretboardLen = fretboardEndX - fretboardStartX

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
	for stringIdx := 0; stringIdx < C.StringCnt; stringIdx++ {
		y := y1 + C.StringMarginTop
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
