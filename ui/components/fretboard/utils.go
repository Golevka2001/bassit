package fretboard

import (
	"math"

	"github.com/Golevka2001/bassit/config"
)

func (m *Model) updateLayoutMappings() {
	if m.Width <= 0 || m.Height <= 0 || m.layoutMappingsUpdated {
		return
	}

	m.fretboardStartX = baseNoteNameWidth
	m.fretboardEndX = max(0, m.Width-1-pluckSignWidth*2) // 2 signs in free playing mode
	if m.forChordTab {
		m.fretboardEndX = max(0, m.Width-1-pluckSignWidth) // 1 sign in chord detection mode
	}
	m.nutStartX = m.fretboardStartX
	m.nutEndX = m.nutStartX + nutWidth - 1
	m.fretboardLen = m.fretboardEndX - m.fretboardStartX
	m.fretboardHeight = 2 + config.StringCnt + (config.StringCnt-1)*stringSpacing

	m.fretwireX = make([]int, config.DisplayedFretCount+1)
	m.xToFretwire = make(map[int]int)
	m.fretCenterX = make([]int, config.DisplayedFretCount+1)
	m.stringY = make([]int, config.StringCnt)
	m.yToString = make(map[int]int)

	// Calculate the position of fretwires
	m.calcFretwirePositions()

	// Calculate the center of frets
	for fretwireIdx := range config.DisplayedFretCount + 1 {
		m.fretCenterX[fretwireIdx] = m.calcFretCenterXPos(fretwireIdx)
	}

	// Calculate the position of strings
	for stringIdx := range config.StringCnt {
		y := 1 + stringIdx*(stringSpacing+1)
		m.stringY[stringIdx] = y
		m.yToString[y] = stringIdx
	}

	m.layoutMappingsUpdated = true
}

func (m *Model) calcFretwirePositions() {
	availableWidth := m.fretboardLen - nutWidth - 1
	scaleFactor := float64(availableWidth) / calcDn(config.DisplayedFretCount)

	m.fretwireX[0] = m.nutEndX
	m.xToFretwire[m.nutEndX] = 0

	prevWidth := math.Inf(1)
	cumulativeError := 0.0

	for i := 1; i <= config.DisplayedFretCount; i++ {
		// Calculate the theoretical spacing
		theoreticalSpacing := calcDn(i) - calcDn(i-1)
		adjustedSpacing := theoreticalSpacing*scaleFactor + cumulativeError
		intSpacing := max(1, int(math.Round(adjustedSpacing)))

		// If the current width is wider than the previous width, shrink it to the previous width
		if float64(intSpacing) > prevWidth {
			intSpacing = int(prevWidth)
		}
		prevWidth = float64(intSpacing)

		// Update the cumulative error
		cumulativeError = adjustedSpacing - float64(intSpacing)
		m.fretwireX[i] = m.fretwireX[i-1] + intSpacing
		m.xToFretwire[m.fretwireX[i]] = i
	}
}

// calcFretCenterXPos calculates the x-coordinate of the center of the nth fret
// Parameters:
//   - n: The fret number (0 means the nut)
func (m *Model) calcFretCenterXPos(n int) int {
	if n < 1 || n > config.DisplayedFretCount {
		return -1
	}

	curFretwireX := m.fretwireX[n]
	prevFretwireX := m.fretwireX[n-1]

	return int(math.Round(float64(curFretwireX+prevFretwireX) / 2))
}

// calcFretwireXPos calculates the x-coordinate of the nth fretwire
// Parameters:
//   - n: The fret number (0 means the nut)
//   - scaleFactor: fretboardLen / d(fretCnt)
//   - offset: The starting position of the fretboard (x-coordinate)
func calcFretwireXPos(n int, scaleFactor float64, offset int) int {
	d := int(math.Round(calcDn(n) * scaleFactor))
	return int(offset + d)
}

// calcDn calculates the distance from the nut to the nth fretwire
func calcDn(n int) float64 {
	// Formula: d(n) = SL * (1 - (2^(-n/12)))
	// where `n` is the fret number,
	// `d(n)` is the distance from the nut to the nth fretwire,
	// `SL` is the scale length (not used here)
	// (Reference: https://www.thekimerers.net/brian/YAFCalc/YAFCalc.html)
	return 1 - math.Pow(float64(2), float64(-n)/12)
}
