package view

import (
	"math"

	"bassit/audio"
	C "bassit/constant"
	"bassit/model"
	"bassit/util"

	"github.com/gdamore/tcell/v2"
)

var (
	borderStartX int
	borderEndX   int
	borderStartY int
	borderEndY   int
	fretboardLen int

	fretWireToX  []int
	xToFretWire  map[int]int
	posMarkerToX []int
	stringToY    []int
	yToString    map[int]int
)

type View struct {
	tcellScreen   tcell.Screen
	bassModel     *model.BassModel
	audioManager  *audio.AudioManager
	width, height int
}

func NewView(
	screen tcell.Screen,
	bassModel *model.BassModel,
	audioManager *audio.AudioManager,
) View {
	screen.SetStyle(tcell.StyleDefault)
	screen.Clear()

	w, h := screen.Size()

	borderStartX = C.FretboardMarginLeft
	borderEndX = w - C.FretboardMarginRight
	borderStartY = C.FretboardMarginTop
	borderEndY = borderStartY + 2*C.StringMarginTop + (C.StringCnt-1)*C.StringSpacing
	fretboardLen = borderEndX - borderStartX

	fretWireToX = make([]int, C.DisplayedFretNum+1)
	xToFretWire = make(map[int]int)
	posMarkerToX = make([]int, C.DisplayedFretNum+1)
	stringToY = make([]int, C.StringCnt)
	yToString = make(map[int]int)

	scaleFactor := float64(fretboardLen-C.NutWidth) / calcDn(C.DisplayedFretNum)
	for fretWireIdx := 0; fretWireIdx <= C.DisplayedFretNum; fretWireIdx++ {
		// Calculate the position of fret wires
		x := calcFretWireXPos(fretWireIdx, scaleFactor, borderStartX+C.NutWidth)
		fretWireToX[fretWireIdx] = x
		xToFretWire[x] = fretWireIdx

		// Calculate the position of position markers
		switch fretWireIdx {
		case 3, 5, 7, 9, 12, 15, 17, 19, 21, 24:
			posMarkerToX[fretWireIdx] = int(math.Round(float64(x+fretWireToX[fretWireIdx-1]) / 2))
		}
	}

	// Calculate the position of strings
	for stringIdx := 0; stringIdx < C.StringCnt; stringIdx++ {
		y := borderStartY + C.StringMarginTop
		if stringIdx > 0 {
			y = stringToY[stringIdx-1] + C.StringSpacing
		}
		stringToY[stringIdx] = y
		yToString[y] = stringIdx
	}

	return View{
		tcellScreen:  screen,
		bassModel:    bassModel,
		audioManager: audioManager,
		width:        w,
		height:       h,
	}
}

func (v *View) Draw() {
	v.tcellScreen.Clear()

	v.drawBassFretboard()
	v.drawBassStrings()
	v.tcellScreen.Show()
}

func (v *View) drawBassFretboard() {
	s := v.tcellScreen

	// Draw the fret board, nut and fret wires
	for x := borderStartX; x <= borderEndX; x++ {
		for y := borderStartY; y <= borderEndY; y++ {
			charToDraw := ' '
			style := tcell.StyleDefault.Foreground(C.FretboardBorderColor).Background(C.FretboardBgColor)

			if x >= borderStartX && x < fretWireToX[0] {
				// The nut
				if y == borderStartY {
					// The upper left corner
					charToDraw = C.FretboardULCornerChar
				} else if y == borderEndY {
					// The lower left corner
					charToDraw = C.FretboardLLCornerChar
				} else {
					// The vertical border
					charToDraw = C.FretboardVBorderChar
				}

				style = tcell.StyleDefault.Foreground(C.NutBorderColor).Background(C.NutBgColor)
			} else if fretWireIdx, ok := xToFretWire[x]; ok {
				// The fret wire
				if y == borderStartY {
					// Fret wire at the upper border
					charToDraw = C.FretWireUpperChar
				} else if y == borderEndY {
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
				if y == borderStartY || y == borderEndY {
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
}

func (v *View) drawBassStrings() {
	s := v.tcellScreen
	b := v.bassModel

	for stringIdx, y := range stringToY {
		curString := b.Strings[stringIdx]

		// Draw base note names
		noteNameStyle := tcell.StyleDefault.Foreground(C.BaseNoteNameFgColor).Background(C.BaseNoteNameBgColor)
		noteName := util.GetNoteNameWithOctave(curString.BaseNote)
		v.DrawLineText(C.StringBaseNoteNameMarginLeft, y, noteName, noteNameStyle)

		// Draw string lines
		rightMostPressedPos := -1
		if curString.CurValidFret > 0 {
			rightMostPressedPos = int(math.Round(float64(fretWireToX[curString.CurValidFret-1]+fretWireToX[curString.CurValidFret]) / 2))
		}
		for x := borderStartX; x <= borderEndX; x++ {
			isVibrating := false
			// If current string is plucked and current position is right to the pressed position
			if curString.PluckedState && x > rightMostPressedPos {
				isVibrating = true
			}

			charToDraw := C.StringChar
			if isVibrating {
				charToDraw = C.VibratingStringChar
			}
			if _, ok := xToFretWire[x]; ok || x == borderStartX {
				charToDraw = C.StringOverFretChar
				if isVibrating {
					charToDraw = C.VibratingStringOverFretChar
				}
			}

			_, _, originStyle, _ := s.GetContent(x, y)
			stringStyle := originStyle.Foreground(C.StringColor)
			s.SetContent(x, y, charToDraw, nil, stringStyle)
		}

		// Draw pressed fret signs
		for fretIdx := curString.CurValidFret; fretIdx > 0; fretIdx-- {
			if fretIdx > C.DisplayedFretNum {
				continue
			}

			if curString.FretPressedStates[fretIdx] {
				charToDraw := C.PressedFretSignChar
				x := int(math.Round(float64(fretWireToX[fretIdx-1]+fretWireToX[fretIdx]) / 2))
				y := stringToY[stringIdx]
				style := tcell.StyleDefault.Foreground(C.PressedFretSignColor).Background(C.FretboardBgColor)
				s.SetContent(x, y, charToDraw, nil, style)
			}
		}

		// Draw plucked string signs
		x := v.width - C.PluckedStringSignMarginRight
		style := tcell.StyleDefault.Foreground(C.StringColor).Background(C.FretboardBgColor)
		if curString.PluckedState {
			style = style.Foreground(C.PluckedStringSignColor)
		}
		s.SetContent(x, y, C.PluckedStringSignChar, nil, style)
	}
}

// calcFretWireXPos calculates the position (on the Screen) of the nth fret wire
// Parameters:
//   - n: The fret number (0 means the nut)
//   - scaleFactor: fretboardLength / d(fretCnt)
//   - offset: The starting position of the fretboard (x-coordinate)
func calcFretWireXPos(n int, scaleFactor float64, offset int) int {
	dist := int(math.Round(calcDn(n) * scaleFactor))
	return int(offset + dist)
}

// calcDn calculates the distance from the nut to the nth fret wire
// Formula: d(n) = SL * (1 - (2^(-n/12)))
// where `n` is the fret number,
// `d(n)` is the distance from the nut to the nth fret wire,
// `SL` is the scale length (not used here)
// (Reference: https://www.thekimerers.net/brian/YAFCalc/YAFCalc.html)
func calcDn(n int) float64 {
	return 1 - math.Pow(float64(2), float64(-n)/12)
}

func (v *View) DrawLineText(x, y int, text string, style tcell.Style) {
	// Draw the text at the specified position
	for offset, char := range text {
		v.tcellScreen.SetContent(x+offset, y, char, nil, style)
	}
}

func (v *View) DrawCenteredText(y int, text string, style tcell.Style) {
	// Calculate the starting position to center the text
	textWidth := len(text)
	startX := (v.width - textWidth) / 2

	// Draw the text at the calculated position
	for offset, char := range text {
		v.tcellScreen.SetContent(startX+offset, y, char, nil, style)
	}
}
