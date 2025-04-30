// Package bass bassit/internal/bass/display.go
package bass

import (
	"math"

	"github.com/Golevka2001/bassit/internal/config"

	"github.com/gdamore/tcell/v2"
)

// RenderTitle draws the title and instructions on the terminal screen
func RenderTitle(bs *BassSimulator) {
	width, height := bs.Screen.Size()

	// Draw the title
	title := "Bassit"
	titleX := (width - len(title)) / 2
	drawText(bs.Screen, titleX, 1, title, tcell.StyleDefault.Foreground(tcell.ColorYellow))

	// Draw the instructions
	instructions := "Press ESC or Ctrl+C to quit"
	instructionsX := (width - len(instructions)) / 2
	drawText(bs.Screen, instructionsX, height-2, instructions, tcell.StyleDefault)
}

// RenderBassFretboard draws a bass fretboard and strings on the terminal screen
func RenderBassFretboard(bs *BassSimulator) {
	width, _ := bs.Screen.Size()

	// The area for the fretboard
	fretboardStartX := config.FretboardMarginLeft
	fretboardEndX := width - config.FretboardMarginRight
	fretboardStartY := config.FretboardMarginTop
	fretboardEndY := config.FretboardMarginTop + 2*config.StringMargin + (len(bs.Strings)-1)*config.StringSpacing
	fretboardLength := fretboardEndX - fretboardStartX // Nut is included

	// Mark the positions to draw the strings and frets
	stringYPosArray := make([]int, len(bs.Strings))
	fretXPosArray := make([]int, config.DisplayedFretNum+1)
	fretXPosSet := make(map[int]bool, config.DisplayedFretNum+1)

	stringYPosArray[0] = fretboardStartY + config.StringMargin
	for i := 1; i < len(bs.Strings); i++ {
		stringYPosArray[i] = stringYPosArray[i-1] + config.StringSpacing
	}

	scaleFactor := float64(fretboardLength-config.NutWidth) / calcDn(config.DisplayedFretNum)
	for i := 0; i <= config.DisplayedFretNum; i++ {
		fretXPosSet[calcFretPosition(i, scaleFactor, fretboardStartX+config.NutWidth)] = true
		fretXPosArray[i] = calcFretPosition(i, scaleFactor, fretboardStartX+config.NutWidth)
	}

	// Draw the border of the fretboard
	for x := fretboardStartX; x <= fretboardEndX; x++ {
		for y := fretboardStartY; y <= fretboardEndY; y++ {
			var char rune
			if x == fretboardStartX {
				if y == fretboardStartY {
					char = config.FretboardULCornerChar
				} else if y == fretboardEndY {
					char = config.FretboardLLCornerChar
				} else {
					char = config.FretboardVBorderChar
				}
			} else {
				if y == fretboardStartY || y == fretboardEndY {
					char = config.FretboardHBorderChar
				}
			}
			if char != 0 {
				bs.Screen.SetContent(x, y, char, nil, tcell.StyleDefault)
			}
		}
	}

	// Draw the nut, fret wires and position markers
	for fretIdx, x := range fretXPosArray {
		// Draw fret wires
		for y := fretboardStartY; y <= fretboardEndY; y++ {
			char := config.FretWireChar
			if y == fretboardStartY {
				char = config.FretWireUpperChar
			} else if y == fretboardEndY {
				char = config.FretWireLowerChar
			}
			bs.Screen.SetContent(x, y, char, nil, tcell.StyleDefault)
		}

		// Draw position markers
		switch fretIdx {
		case 3, 5, 7, 9, 15, 17, 19, 21:
			// Single dot
			markXPos := int(math.Round(float64(fretXPosArray[fretIdx]+fretXPosArray[fretIdx-1]) / 2))
			markYPos := int(math.Round(float64(stringYPosArray[0]+stringYPosArray[len(bs.Strings)-1]) / 2))
			bs.Screen.SetContent(markXPos, markYPos, config.PositionMarkerChar, nil, tcell.StyleDefault)
		case 12, 24:
			// Double dot
			markXPos := int(math.Round(float64(fretXPosArray[fretIdx]+fretXPosArray[fretIdx-1]) / 2))
			markY1Pos := int(math.Round(float64(stringYPosArray[0]+stringYPosArray[1]) / 2))
			markY2Pos := int(math.Round(float64(stringYPosArray[len(bs.Strings)-2]+stringYPosArray[len(bs.Strings)-1]) / 2))
			bs.Screen.SetContent(markXPos, markY1Pos, config.PositionMarkerChar, nil, tcell.StyleDefault)
			bs.Screen.SetContent(markXPos, markY2Pos, config.PositionMarkerChar, nil, tcell.StyleDefault)
		}
	}

	// Draw the strings
	for stringIdx, y := range stringYPosArray {
		// Draw string names
		drawText(bs.Screen, config.StringBaseNoteNameMarginLeft, y, bs.Strings[stringIdx].BaseNoteName, tcell.StyleDefault.Foreground(tcell.ColorGreen))

		// Calculate the pressed fret position
		pressedFret := bs.Strings[stringIdx].PressedFret
		pressedFretXPos := -1
		if pressedFret > 0 {
			pressedFretXPos = int(math.Round(float64(fretXPosArray[pressedFret]+fretXPosArray[pressedFret-1]) / 2))
		}

		// Draw string lines
		for x := fretboardStartX; x <= fretboardEndX; x++ {
			// If any string is pressed, draw a mark at the corresponding position
			if pressedFretXPos == x {
				bs.Screen.SetContent(x, y, config.PressedFretChar, nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
				continue
			}

			char := config.StringChar
			if bs.Strings[stringIdx].IsPlucked {
				char = config.PluckedStringChar
			}
			if _, exists := fretXPosSet[x]; exists || x == fretboardStartX || x == fretboardEndX {
				if bs.Strings[stringIdx].IsPlucked {
					char = config.PluckedStringOverFretChar
				} else {
					char = config.StringOverFretChar
				}
			}
			bs.Screen.SetContent(x, y, char, nil, tcell.StyleDefault)
		}
	}

	bs.Screen.Show()
}

// drawText is a helper function to draw text at a specific position on the screen
func drawText(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	for i, r := range []rune(text) {
		screen.SetContent(x+i, y, r, nil, style)
	}
}

// calcFretPosition calculates the position (on the screen) of the nth fret
// Parameters:
//   - n: The fret number (0 means the nut)
//   - scaleFactor: fretboardLength / d(fretCnt)
//   - offset: The starting position of the fretboard (x-coordinate)
func calcFretPosition(n int, scaleFactor float64, offset int) int {
	dist := int(math.Round(calcDn(n) * scaleFactor))
	return offset + dist
}

// calcDn calculates the distance from the nut to the nth fret
// Formula: d(n) = SL * (1 - (2^(-n/12)))
// where `n` is the fret number,
// `d(n)` is the distance from the nut to the nth fret,
// `SL` is the scale length (not used here)
// (Reference: https://www.thekimerers.net/brian/YAFCalc/YAFCalc.html)
func calcDn(n int) float64 {
	return 1 - math.Pow(float64(2), float64(-n)/12)
}
