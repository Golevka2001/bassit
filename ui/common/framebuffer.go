package common

import (
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
)

type Cell struct {
	R     rune
	Style lipgloss.Style
	Set   bool

	// RenderedStr stores the rendered content (ANSI escape codes)
	// When set, R and Style are ignored
	RenderedStr string // optional
}

type FrameBuffer struct {
	Width, Height int
	Cells         [][]Cell
}

func NewFrameBuffer(width, height int) FrameBuffer {
	cells := make([][]Cell, height)
	for y := range cells {
		cells[y] = make([]Cell, width)
	}
	return FrameBuffer{
		Width:  width,
		Height: height,
		Cells:  cells,
	}
}

// Empty returns true if all cells are empty
func (fb *FrameBuffer) Empty() bool {
	for _, row := range fb.Cells {
		for _, cell := range row {
			if cell.Set {
				return false
			}
		}
	}
	return true
}

// DrawRune draws a single rune at the given coordinates
// Note: this method does not support rendered content (ANSI escape codes)
// Parameters:
// - x, y: coordinates of the cell to draw the rune in
// - r: rune to draw
// - style: style to apply to the rune
// - inheritStyle: if true, the new style will be `style.Inherit(originalStyle)`
func (fb *FrameBuffer) DrawRune(x, y int, r rune, style lipgloss.Style, inheritStyle bool) {
	if x >= 0 && x < fb.Width && y >= 0 && y < fb.Height {
		fb.Cells[y][x].R = r
		fb.SetStyle(x, y, style, inheritStyle)
	}
}

// DrawString draws a string at the given coordinates, line by line
// If a new line is encountered, the string is drawn in the next line, starting from the same x coordinate
// Note: this method does not support rendered content (ANSI escape codes)
// Parameters:
// - x, y: coordinates of the first cell to draw the string in
// - s: string to draw (newlines are supported)
// - style: style to apply to the whole string
// - inheritStyle: if true, the new style will be `style.Inherit(originalStyle)`
func (fb *FrameBuffer) DrawString(x, y int, s string, style lipgloss.Style, inheritStyle bool) {
	rows := strings.Split(s, "\n")
	for dy, row := range rows {
		newY := y + dy
		if newY < 0 || newY >= fb.Height {
			break
		}
		for dx, r := range row {
			newX := x + dx
			if newX < 0 || newX >= fb.Width {
				break
			}
			fb.DrawRune(newX, newY, r, style, inheritStyle)
		}
	}
}

// DrawRenderedString draws rendered content at the given coordinates
// Parameters:
// - x, y: coordinates of the first cell to draw the raw content in
// - rendered: rendered content to draw (ANSI escape codes)
func (fb *FrameBuffer) DrawRenderedString(x, y int, rendered string) {
	rows := strings.Split(rendered, "\n")
	for dy, row := range rows {
		newY := y + dy
		if newY < 0 || newY >= fb.Height {
			break
		}
		runes := []rune(row)
		for dx, r := range runes {
			newX := x + dx
			if newX < 0 || newX >= fb.Width {
				break
			}
			fb.Cells[newY][newX] = Cell{RenderedStr: string(r), Set: true}
		}
	}
}

// SetStyle sets the style of a cell
// Parameters:
// - x, y: coordinates of the cell to set the style of
// - style: style to set
// - inheritStyle: if true, the new style will be `style.Inherit(originalStyle)`
func (fb *FrameBuffer) SetStyle(x, y int, style lipgloss.Style, inheritStyle bool) {
	if x >= 0 && x < fb.Width && y >= 0 && y < fb.Height {
		if inheritStyle {
			fb.Cells[y][x].Style = style.Inherit(fb.Cells[y][x].Style)
		} else {
			fb.Cells[y][x].Style = style
		}
		fb.Cells[y][x].Set = true
	}
}

func (fb *FrameBuffer) Clear() {
	for y := range fb.Cells {
		for x := range fb.Cells[y] {
			fb.Cells[y][x] = Cell{}
		}
	}
}

func (fb *FrameBuffer) Render() string {
	var b strings.Builder
	for i, row := range fb.Cells {
		for _, cell := range row {
			if cell.Set {
				if cell.RenderedStr != "" {
					b.WriteString(cell.RenderedStr)
				} else {
					b.WriteString(cell.Style.Render(string(cell.R)))
				}
			} else {
				b.WriteRune(' ')
			}
		}
		if i < len(fb.Cells)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}
