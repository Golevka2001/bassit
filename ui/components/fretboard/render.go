package fretboard

import (
	"strings"

	"github.com/Golevka2001/bassit/config"
	"github.com/Golevka2001/bassit/utils"

	"github.com/charmbracelet/lipgloss/v2"
)

const (
	stringSpacing     = 1
	nutWidth          = 3
	baseNoteNameWidth = 4
	pluckSignWidth    = 2
	blockInlayMarginX = 1
	blockInlayMarginY = 2
)

func (m *Model) drawFretboard() {
	for y := range m.fretboardHeight {
		if y < 0 || y >= m.Height {
			break
		}

		for x := m.fretboardStartX; x <= m.fretboardEndX; x++ {
			var r rune
			var style lipgloss.Style

			// Draw the nut
			if x >= m.nutStartX && x <= m.nutEndX {
				switch x {
				// Left border
				case m.nutStartX:
					style = m.style.nutBorderStyle
					switch y {
					case 0:
						r = m.theme.NutULCornerChar
					case m.fretboardHeight - 1:
						r = m.theme.NutLLCornerChar
					default:
						r = m.theme.NutVBorderChar
					}

				// Right border
				case m.nutEndX:
					switch y {
					case 0:
						r = m.theme.NutURCornerChar
						style = m.style.nutBorderStyle
					case m.fretboardHeight - 1:
						r = m.theme.NutLRCornerChar
						style = m.style.nutBorderStyle
					default:
						r = m.theme.NutVBorderChar
						style = m.style.nutBgStyle
					}

				default:
					switch y {
					case 0:
						r = m.theme.NutHBorderChar
						style = m.style.nutBorderStyle
					case m.fretboardHeight - 1:
						r = m.theme.NutHBorderChar
						style = m.style.nutBorderStyle
					default:
						r = ' '
						style = m.style.nutBgStyle
					}
				}
				m.frameBuf.DrawRune(x, y, r, style, false)
			}

			// Draw the fretboard background and fretwires
			if x >= m.nutEndX+1 && x <= m.fretboardEndX {
				switch y {
				case 0, m.fretboardHeight - 1:
					r = m.theme.FretboardHBorderChar
					style = m.style.fretboardBorderStyle
				default:
					r = ' '
					style = m.style.fretboardBgStyle
				}
				m.frameBuf.DrawRune(x, y, r, style, false)

				// Draw the fretwires
				if _, ok := m.xToFretwire[x]; ok {
					switch y {
					case 0:
						r = m.theme.FretwireTopBorderChar
						style = m.style.fretboardBorderStyle
					case m.fretboardHeight - 1:
						r = m.theme.FretwireBottomBorderChar
						style = m.style.fretboardBorderStyle
					default:
						r = m.theme.FretwireChar
						style = m.style.fretwireStyle
					}
					m.frameBuf.DrawRune(x, y, r, style, true)
				}
			}
		}
	}
}

func (m *Model) drawFretboardInlays() {
	switch m.theme.InlayShape {
	case config.DotInlayShape:
		for fretIdx, fretCenterX := range m.fretCenterX {
			switch fretIdx {
			// Single dot
			case 3, 5, 7, 9, 15, 17, 19, 21:
				y := m.fretboardHeight / 2
				m.frameBuf.DrawRune(fretCenterX, y, m.theme.InlayChar, m.style.dotInlayStyle, true)

			// Double dot
			case 12, 24:
				y1 := m.stringY[0] + (stringSpacing+1)/2
				y2 := m.stringY[config.StringCnt-1] - (stringSpacing+1)/2
				m.frameBuf.DrawRune(fretCenterX, y1, m.theme.InlayChar, m.style.dotInlayStyle, true)
				m.frameBuf.DrawRune(fretCenterX, y2, m.theme.InlayChar, m.style.dotInlayStyle, true)
			}

		}

	case config.BlockInlayShape:
		for fretIdx, curFretwireX := range m.fretwireX {
			switch fretIdx {
			case 1, 3, 5, 7, 9, 12, 15, 17, 19, 21, 24:
				prevFretwireX := m.fretwireX[fretIdx-1]
				x := prevFretwireX + blockInlayMarginX + 1
				y := blockInlayMarginY
				w := max(0, curFretwireX-prevFretwireX-1-blockInlayMarginX*2)
				h := max(0, m.fretboardHeight-blockInlayMarginY*2)
				inlay := strings.Repeat(strings.Repeat(" ", w)+"\n", h)
				m.frameBuf.DrawString(x, y, inlay, m.style.blockInlayStyle, true)
			}
		}
	}
}

func (m *Model) drawStrings() {
	for stringIdx, y := range m.stringY {
		if y < 0 || y >= m.Height {
			break
		}

		// Draw the base notes
		m.frameBuf.DrawString(0, y, utils.GetNoteNameWithOctave(m.baseNotes[stringIdx]), m.style.baseNoteStyle, false)

		// Draw the pluck signs
		m.frameBuf.DrawRune(m.fretboardEndX+2, y, m.theme.NotPluckedStringChar, m.style.stringStyle, false)
		if !m.forChordTab {
			m.frameBuf.DrawRune(m.fretboardEndX+4, y, m.theme.NotPluckedStringChar, m.style.stringStyle, false)
		}

		// Draw the strings
		for x := m.fretboardStartX; x <= m.fretboardEndX; x++ {
			var r rune
			_, crossFretwire := m.xToFretwire[x]

			switch {
			case x == m.nutStartX || x == m.nutEndX:
				r = m.theme.StringOverBoarderChar
			case crossFretwire:
				r = m.theme.StringOverFretwireChar
			default:
				r = m.theme.StringChar
			}
			m.frameBuf.DrawRune(x, y, r, m.style.stringStyle, true)
		}
	}
}

func (m *Model) drawPressedFret(stringIdx, fretIdx int) {
	if stringIdx < 0 || stringIdx >= config.StringCnt ||
		fretIdx < 0 || fretIdx > config.DisplayedFretCount {
		return
	}
	x := m.fretCenterX[fretIdx]
	y := m.stringY[stringIdx]

	m.frameBuf.DrawRune(x, y, m.theme.PressedFretSignChar, m.style.pressedFretSignStyle, true)
}

func (m *Model) restorePressedFret(stringIdx, fretIdx int) {
	if stringIdx < 0 || stringIdx >= config.StringCnt ||
		fretIdx < 0 || fretIdx > config.DisplayedFretCount {
		return
	}
	x := m.fretCenterX[fretIdx]
	y := m.stringY[stringIdx]

	var r rune
	if _, ok := m.xToFretwire[x]; ok {
		r = m.theme.StringOverFretwireChar
	} else {
		r = m.theme.StringChar
	}

	m.frameBuf.DrawRune(x, y, r, m.style.stringStyle, true)
}

func (m *Model) drawVibratingString(stringIdx int, fretIdx int, position int) {
	if stringIdx < 0 || stringIdx >= config.StringCnt ||
		fretIdx < 0 || fretIdx > config.DisplayedFretCount {
		return
	}
	// Draw the sign
	y := m.stringY[stringIdx]
	if position == 0 {
		m.frameBuf.DrawRune(m.fretboardEndX+2, y, m.theme.PluckedStringSignChar, m.style.pluckedStringSignStyle, false)
		m.frameBuf.DrawRune(m.fretboardEndX+4, y, m.theme.NotPluckedStringChar, m.style.stringStyle, false)
	} else {
		m.frameBuf.DrawRune(m.fretboardEndX+4, y, m.theme.PluckedStringSignChar, m.style.pluckedStringSignStyle, false)
		m.frameBuf.DrawRune(m.fretboardEndX+2, y, m.theme.NotPluckedStringChar, m.style.stringStyle, false)
	}

	// Draw the vibrating string
	rightMostPressedPos := m.fretCenterX[fretIdx]
	if fretIdx == 0 {
		rightMostPressedPos = m.nutEndX
	}
	if rightMostPressedPos == -1 {
		return
	}
	for x := rightMostPressedPos + 1; x <= m.fretboardEndX; x++ {
		var r rune
		_, crossFretwire := m.xToFretwire[x]

		switch {
		case x == m.nutStartX || x == m.nutEndX:
			r = m.theme.VibratingStringOverBoarderChar
		case crossFretwire:
			r = m.theme.VibratingStringOverFretwireChar
		default:
			r = m.theme.VibratingStringChar
		}
		m.frameBuf.DrawRune(x, y, r, m.style.stringStyle, true)
	}
}

func (m *Model) restoreVibratingString(stringIdx int, fretIdx int) {
	if stringIdx < 0 || stringIdx >= config.StringCnt ||
		fretIdx < 0 || fretIdx > config.DisplayedFretCount {
		return
	}
	// Restore the sign
	y := m.stringY[stringIdx]
	m.frameBuf.DrawRune(m.fretboardEndX+2, y, m.theme.NotPluckedStringChar, m.style.stringStyle, false)
	m.frameBuf.DrawRune(m.fretboardEndX+4, y, m.theme.NotPluckedStringChar, m.style.stringStyle, false)

	// Restore the vibrating string
	rightMostPressedPos := m.fretCenterX[fretIdx]
	if fretIdx == 0 {
		rightMostPressedPos = m.nutEndX
	}
	if rightMostPressedPos == -1 {
		return
	}
	for x := rightMostPressedPos + 1; x <= m.fretboardEndX; x++ {
		var r rune
		_, crossFretwire := m.xToFretwire[x]

		switch {
		case x == m.nutStartX || x == m.nutEndX:
			r = m.theme.StringOverBoarderChar
		case crossFretwire:
			r = m.theme.StringOverFretwireChar
		default:
			r = m.theme.StringChar
		}
		m.frameBuf.DrawRune(x, y, r, m.style.stringStyle, true)
	}
}

func (m *Model) drawXOnString(stringIdx int) {
	if stringIdx < 0 || stringIdx >= config.StringCnt {
		return
	}
	y := m.stringY[stringIdx]
	m.frameBuf.DrawRune(m.fretboardEndX+2, y, 'X', m.style.stringStyle, false)
}

func (m *Model) removeXFromString(stringIdx int) {
	if stringIdx < 0 || stringIdx >= config.StringCnt {
		return
	}
	y := m.stringY[stringIdx]
	m.frameBuf.DrawRune(m.fretboardEndX+2, y, m.theme.NotPluckedStringChar, m.style.stringStyle, false)
}
