package common

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewFrameBuffer(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{
			name:   "small buffer",
			width:  5,
			height: 3,
		},
		{
			name:   "single cell",
			width:  1,
			height: 1,
		},
		{
			name:   "wide buffer",
			width:  100,
			height: 10,
		},
		{
			name:   "tall buffer",
			width:  10,
			height: 100,
		},
		{
			name:   "zero width",
			width:  0,
			height: 5,
		},
		{
			name:   "zero height",
			width:  5,
			height: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := NewFrameBuffer(tt.width, tt.height)

			assert.Equal(t, tt.width, fb.Width)
			assert.Equal(t, tt.height, fb.Height)
			assert.Equal(t, tt.height, len(fb.Cells))

			if tt.height > 0 {
				assert.Equal(t, tt.width, len(fb.Cells[0]))
			}

			for y := 0; y < tt.height; y++ {
				for x := 0; x < tt.width; x++ {
					cell := fb.Cells[y][x]
					assert.False(t, cell.Set)
					assert.Equal(t, rune(0), cell.R)
					assert.Equal(t, "", cell.RenderedStr)
				}
			}
		})
	}
}

func TestFrameBuffer_Empty(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() FrameBuffer
		expected bool
	}{
		{
			name: "new buffer is empty",
			setup: func() FrameBuffer {
				return NewFrameBuffer(5, 3)
			},
			expected: true,
		},
		{
			name: "buffer with one set cell is not empty",
			setup: func() FrameBuffer {
				fb := NewFrameBuffer(5, 3)
				fb.DrawRune(2, 1, 'A', lipgloss.NewStyle(), false)
				return fb
			},
			expected: false,
		},
		{
			name: "cleared buffer is empty",
			setup: func() FrameBuffer {
				fb := NewFrameBuffer(5, 3)
				fb.DrawRune(2, 1, 'A', lipgloss.NewStyle(), false)
				fb.Clear()
				return fb
			},
			expected: true,
		},
		{
			name: "buffer with rendered string is not empty",
			setup: func() FrameBuffer {
				fb := NewFrameBuffer(5, 3)
				fb.DrawRenderedString(0, 0, "test")
				return fb
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := tt.setup()
			assert.Equal(t, tt.expected, fb.Empty())
		})
	}
}

func TestFrameBuffer_DrawRune(t *testing.T) {
	tests := []struct {
		name         string
		width        int
		height       int
		x            int
		y            int
		r            rune
		style        lipgloss.Style
		inheritStyle bool
		expectSet    bool
	}{
		{
			name:      "draw within bounds",
			width:     5,
			height:    3,
			x:         2,
			y:         1,
			r:         'A',
			style:     lipgloss.NewStyle().Foreground(lipgloss.Color("red")),
			expectSet: true,
		},
		{
			name:      "draw at origin",
			width:     5,
			height:    3,
			x:         0,
			y:         0,
			r:         'B',
			style:     lipgloss.NewStyle(),
			expectSet: true,
		},
		{
			name:      "draw at bottom right",
			width:     5,
			height:    3,
			x:         4,
			y:         2,
			r:         'C',
			style:     lipgloss.NewStyle(),
			expectSet: true,
		},
		{
			name:      "draw out of bounds - negative x",
			width:     5,
			height:    3,
			x:         -1,
			y:         1,
			r:         'D',
			style:     lipgloss.NewStyle(),
			expectSet: false,
		},
		{
			name:      "draw out of bounds - negative y",
			width:     5,
			height:    3,
			x:         2,
			y:         -1,
			r:         'E',
			style:     lipgloss.NewStyle(),
			expectSet: false,
		},
		{
			name:      "draw out of bounds - x too large",
			width:     5,
			height:    3,
			x:         5,
			y:         1,
			r:         'F',
			style:     lipgloss.NewStyle(),
			expectSet: false,
		},
		{
			name:      "draw out of bounds - y too large",
			width:     5,
			height:    3,
			x:         2,
			y:         3,
			r:         'G',
			style:     lipgloss.NewStyle(),
			expectSet: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := NewFrameBuffer(tt.width, tt.height)
			fb.DrawRune(tt.x, tt.y, tt.r, tt.style, tt.inheritStyle)

			if tt.expectSet {
				assert.True(t, fb.Cells[tt.y][tt.x].Set)
				assert.Equal(t, tt.r, fb.Cells[tt.y][tt.x].R)
				assert.Equal(t, tt.style, fb.Cells[tt.y][tt.x].Style)
			} else {
				for y := 0; y < tt.height; y++ {
					for x := 0; x < tt.width; x++ {
						assert.False(t, fb.Cells[y][x].Set)
					}
				}
			}
		})
	}
}

func TestFrameBuffer_DrawString(t *testing.T) {
	tests := []struct {
		name         string
		width        int
		height       int
		x            int
		y            int
		s            string
		style        lipgloss.Style
		inheritStyle bool
		validate     func(t *testing.T, fb FrameBuffer)
	}{
		{
			name:   "draw simple string",
			width:  10,
			height: 5,
			x:      2,
			y:      1,
			s:      "hello",
			style:  lipgloss.NewStyle().Foreground(lipgloss.Color("blue")),
			validate: func(t *testing.T, fb FrameBuffer) {
				expected := []rune("hello")
				for i, r := range expected {
					cell := fb.Cells[1][2+i]
					assert.True(t, cell.Set)
					assert.Equal(t, r, cell.R)
				}
			},
		},
		{
			name:   "draw multiline string",
			width:  10,
			height: 5,
			x:      1,
			y:      1,
			s:      "line1\nline2\nline3",
			style:  lipgloss.NewStyle(),
			validate: func(t *testing.T, fb FrameBuffer) {
				lines := []string{"line1", "line2", "line3"}
				for lineIdx, line := range lines {
					for charIdx, r := range line {
						cell := fb.Cells[1+lineIdx][1+charIdx]
						assert.True(t, cell.Set)
						assert.Equal(t, r, cell.R)
					}
				}
			},
		},
		{
			name:   "draw string that goes out of bounds horizontally",
			width:  5,
			height: 3,
			x:      3,
			y:      1,
			s:      "toolong",
			style:  lipgloss.NewStyle(),
			validate: func(t *testing.T, fb FrameBuffer) {
				assert.True(t, fb.Cells[1][3].Set)
				assert.Equal(t, 't', fb.Cells[1][3].R)
				assert.True(t, fb.Cells[1][4].Set)
				assert.Equal(t, 'o', fb.Cells[1][4].R)
			},
		},
		{
			name:   "draw string that goes out of bounds vertically",
			width:  10,
			height: 2,
			x:      0,
			y:      1,
			s:      "line1\nline2\nline3",
			style:  lipgloss.NewStyle(),
			validate: func(t *testing.T, fb FrameBuffer) {
				line1 := "line1"
				for i, r := range line1 {
					cell := fb.Cells[1][i]
					assert.True(t, cell.Set)
					assert.Equal(t, r, cell.R)
				}
			},
		},
		{
			name:   "draw empty string",
			width:  5,
			height: 3,
			x:      2,
			y:      1,
			s:      "",
			style:  lipgloss.NewStyle(),
			validate: func(t *testing.T, fb FrameBuffer) {
				assert.True(t, fb.Empty())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := NewFrameBuffer(tt.width, tt.height)
			fb.DrawString(tt.x, tt.y, tt.s, tt.style, tt.inheritStyle)
			tt.validate(t, fb)
		})
	}
}

func TestFrameBuffer_DrawRenderedString(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		height   int
		x        int
		y        int
		rendered string
		validate func(t *testing.T, fb FrameBuffer)
	}{
		{
			name:     "draw simple rendered string",
			width:    10,
			height:   5,
			x:        2,
			y:        1,
			rendered: "hello",
			validate: func(t *testing.T, fb FrameBuffer) {
				expected := []rune("hello")
				for i, r := range expected {
					cell := fb.Cells[1][2+i]
					assert.True(t, cell.Set)
					assert.Equal(t, string(r), cell.RenderedStr)
				}
			},
		},
		{
			name:     "draw multiline rendered string",
			width:    10,
			height:   5,
			x:        1,
			y:        1,
			rendered: "line1\nline2\nline3",
			validate: func(t *testing.T, fb FrameBuffer) {
				lines := []string{"line1", "line2", "line3"}
				for lineIdx, line := range lines {
					for charIdx, r := range line {
						cell := fb.Cells[1+lineIdx][1+charIdx]
						assert.True(t, cell.Set)
						assert.Equal(t, string(r), cell.RenderedStr)
					}
				}
			},
		},
		{
			name:     "draw rendered string with ANSI codes",
			width:    10,
			height:   3,
			x:        0,
			y:        0,
			rendered: "\x1b[31mred\x1b[0m",
			validate: func(t *testing.T, fb FrameBuffer) {
				rendered := "\x1b[31mred\x1b[0m"
				runes := []rune(rendered)
				for i, r := range runes {
					if i < fb.Width {
						cell := fb.Cells[0][i]
						assert.True(t, cell.Set)
						assert.Equal(t, string(r), cell.RenderedStr)
					}
				}
			},
		},
		{
			name:     "draw rendered string out of bounds",
			width:    3,
			height:   2,
			x:        2,
			y:        1,
			rendered: "toolong",
			validate: func(t *testing.T, fb FrameBuffer) {
				assert.True(t, fb.Cells[1][2].Set)
				assert.Equal(t, "t", fb.Cells[1][2].RenderedStr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := NewFrameBuffer(tt.width, tt.height)
			fb.DrawRenderedString(tt.x, tt.y, tt.rendered)
			tt.validate(t, fb)
		})
	}
}

func TestFrameBuffer_SetStyle(t *testing.T) {
	tests := []struct {
		name         string
		width        int
		height       int
		x            int
		y            int
		style        lipgloss.Style
		inheritStyle bool
		expectSet    bool
	}{
		{
			name:      "set style within bounds",
			width:     5,
			height:    3,
			x:         2,
			y:         1,
			style:     lipgloss.NewStyle().Foreground(lipgloss.Color("green")),
			expectSet: true,
		},
		{
			name:      "set style out of bounds - negative x",
			width:     5,
			height:    3,
			x:         -1,
			y:         1,
			style:     lipgloss.NewStyle(),
			expectSet: false,
		},
		{
			name:      "set style out of bounds - x too large",
			width:     5,
			height:    3,
			x:         5,
			y:         1,
			style:     lipgloss.NewStyle(),
			expectSet: false,
		},
		{
			name:      "set style out of bounds - negative y",
			width:     5,
			height:    3,
			x:         2,
			y:         -1,
			style:     lipgloss.NewStyle(),
			expectSet: false,
		},
		{
			name:      "set style out of bounds - y too large",
			width:     5,
			height:    3,
			x:         2,
			y:         3,
			style:     lipgloss.NewStyle(),
			expectSet: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := NewFrameBuffer(tt.width, tt.height)
			fb.SetStyle(tt.x, tt.y, tt.style, tt.inheritStyle)

			if tt.expectSet {
				assert.True(t, fb.Cells[tt.y][tt.x].Set)
				assert.Equal(t, tt.style, fb.Cells[tt.y][tt.x].Style)
			} else {
				for y := 0; y < tt.height; y++ {
					for x := 0; x < tt.width; x++ {
						assert.False(t, fb.Cells[y][x].Set)
					}
				}
			}
		})
	}
}

func TestFrameBuffer_SetStyle_Inherit(t *testing.T) {
	fb := NewFrameBuffer(5, 3)

	originalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("red"))
	fb.SetStyle(2, 1, originalStyle, false)

	newStyle := lipgloss.NewStyle().Background(lipgloss.Color("blue"))
	fb.SetStyle(2, 1, newStyle, true)

	expectedStyle := newStyle.Inherit(originalStyle)
	assert.Equal(t, expectedStyle, fb.Cells[1][2].Style)
	assert.True(t, fb.Cells[1][2].Set)
}

func TestFrameBuffer_Clear(t *testing.T) {
	fb := NewFrameBuffer(5, 3)

	fb.DrawRune(1, 1, 'A', lipgloss.NewStyle().Foreground(lipgloss.Color("red")), false)
	fb.DrawString(2, 0, "test", lipgloss.NewStyle(), false)
	fb.DrawRenderedString(0, 2, "rendered")

	assert.False(t, fb.Empty())

	fb.Clear()

	assert.True(t, fb.Empty())
	for y := 0; y < fb.Height; y++ {
		for x := 0; x < fb.Width; x++ {
			cell := fb.Cells[y][x]
			assert.False(t, cell.Set)
			assert.Equal(t, rune(0), cell.R)
			assert.Equal(t, "", cell.RenderedStr)
			assert.Equal(t, lipgloss.Style{}, cell.Style)
		}
	}
}

func TestFrameBuffer_Render(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() FrameBuffer
		expected string
	}{
		{
			name: "empty buffer",
			setup: func() FrameBuffer {
				return NewFrameBuffer(3, 2)
			},
			expected: "   \n   ",
		},
		{
			name: "buffer with single character",
			setup: func() FrameBuffer {
				fb := NewFrameBuffer(3, 2)
				fb.DrawRune(1, 0, 'A', lipgloss.NewStyle(), false)
				return fb
			},
			expected: " A \n   ",
		},
		{
			name: "buffer with string",
			setup: func() FrameBuffer {
				fb := NewFrameBuffer(5, 2)
				fb.DrawString(0, 0, "hello", lipgloss.NewStyle(), false)
				return fb
			},
			expected: "hello\n     ",
		},
		{
			name: "buffer with multiline string",
			setup: func() FrameBuffer {
				fb := NewFrameBuffer(5, 3)
				fb.DrawString(0, 0, "hi\nbye", lipgloss.NewStyle(), false)
				return fb
			},
			expected: "hi   \nbye  \n     ",
		},
		{
			name: "buffer with rendered string",
			setup: func() FrameBuffer {
				fb := NewFrameBuffer(4, 2)
				fb.DrawRenderedString(0, 0, "test")
				return fb
			},
			expected: "test\n    ",
		},
		{
			name: "buffer with mixed content",
			setup: func() FrameBuffer {
				fb := NewFrameBuffer(6, 3)
				fb.DrawRune(0, 0, 'A', lipgloss.NewStyle(), false)
				fb.DrawString(2, 0, "BC", lipgloss.NewStyle(), false)
				fb.DrawRenderedString(0, 1, "def")
				fb.DrawRune(5, 2, 'Z', lipgloss.NewStyle(), false)
				return fb
			},
			expected: "A BC  \ndef   \n     Z",
		},
		{
			name: "single row buffer",
			setup: func() FrameBuffer {
				fb := NewFrameBuffer(4, 1)
				fb.DrawString(0, 0, "test", lipgloss.NewStyle(), false)
				return fb
			},
			expected: "test",
		},
		{
			name: "single column buffer",
			setup: func() FrameBuffer {
				fb := NewFrameBuffer(1, 3)
				fb.DrawRune(0, 0, 'A', lipgloss.NewStyle(), false)
				fb.DrawRune(0, 2, 'C', lipgloss.NewStyle(), false)
				return fb
			},
			expected: "A\n \nC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fb := tt.setup()
			result := fb.Render()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFrameBuffer_Render_WithStyles(t *testing.T) {
	fb := NewFrameBuffer(3, 2)

	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("red"))
	fb.DrawRune(0, 0, 'A', redStyle, false)
	fb.DrawRune(1, 0, 'B', lipgloss.NewStyle(), false)

	result := fb.Render()

	lines := strings.Split(result, "\n")
	assert.Equal(t, 2, len(lines))

	firstLine := lines[0]
	assert.Contains(t, firstLine, "A")
	assert.Contains(t, firstLine, "B")

	assert.Equal(t, "   ", lines[1])
}

func TestFrameBuffer_EdgeCases(t *testing.T) {
	t.Run("zero size buffer", func(t *testing.T) {
		fb := NewFrameBuffer(0, 0)
		assert.True(t, fb.Empty())
		assert.Equal(t, "", fb.Render())

		fb.DrawRune(0, 0, 'A', lipgloss.NewStyle(), false)
		fb.DrawString(0, 0, "test", lipgloss.NewStyle(), false)
		fb.DrawRenderedString(0, 0, "test")
		fb.SetStyle(0, 0, lipgloss.NewStyle(), false)
		fb.Clear()
	})

	t.Run("unicode characters", func(t *testing.T) {
		fb := NewFrameBuffer(10, 3)

		fb.DrawRune(0, 0, '🎵', lipgloss.NewStyle(), false)
		fb.DrawRune(1, 0, '中', lipgloss.NewStyle(), false)
		fb.DrawRune(2, 0, 'ñ', lipgloss.NewStyle(), false)
		fb.DrawString(0, 1, "🎸🎹🎤", lipgloss.NewStyle(), false)
		fb.DrawRenderedString(0, 2, "αβγδε")

		assert.False(t, fb.Empty())
		result := fb.Render()
		assert.Contains(t, result, "🎵")
		assert.Contains(t, result, "中")
		assert.Contains(t, result, "ñ")
	})

	t.Run("overwrite cells", func(t *testing.T) {
		fb := NewFrameBuffer(5, 3)

		fb.DrawRune(2, 1, 'A', lipgloss.NewStyle().Foreground(lipgloss.Color("red")), false)
		assert.Equal(t, 'A', fb.Cells[1][2].R)

		fb.DrawRune(2, 1, 'B', lipgloss.NewStyle().Foreground(lipgloss.Color("blue")), false)
		assert.Equal(t, 'B', fb.Cells[1][2].R)

		fb.DrawRenderedString(2, 1, "C")
		assert.Equal(t, "C", fb.Cells[1][2].RenderedStr)
		assert.Equal(t, rune(0), fb.Cells[1][2].R)
	})

	t.Run("mixed rendered and styled content", func(t *testing.T) {
		fb := NewFrameBuffer(6, 2)

		fb.DrawRune(0, 0, 'A', lipgloss.NewStyle().Foreground(lipgloss.Color("red")), false)
		fb.DrawRenderedString(1, 0, "BC")
		fb.DrawString(3, 0, "DE", lipgloss.NewStyle().Background(lipgloss.Color("blue")), false)

		result := fb.Render()
		lines := strings.Split(result, "\n")

		firstLine := lines[0]
		assert.Contains(t, firstLine, "A")
		assert.Contains(t, firstLine, "B")
		assert.Contains(t, firstLine, "C")
		assert.Contains(t, firstLine, "D")
		assert.Contains(t, firstLine, "E")
	})
}

func TestFrameBuffer_Integration(t *testing.T) {
	t.Run("complex drawing scenario", func(t *testing.T) {
		fb := NewFrameBuffer(20, 10)

		borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("white"))
		for x := 0; x < fb.Width; x++ {
			fb.DrawRune(x, 0, '─', borderStyle, false)
			fb.DrawRune(x, fb.Height-1, '─', borderStyle, false)
		}
		for y := 0; y < fb.Height; y++ {
			fb.DrawRune(0, y, '│', borderStyle, false)
			fb.DrawRune(fb.Width-1, y, '│', borderStyle, false)
		}

		fb.DrawRune(0, 0, '┌', borderStyle, false)
		fb.DrawRune(fb.Width-1, 0, '┐', borderStyle, false)
		fb.DrawRune(0, fb.Height-1, '└', borderStyle, false)
		fb.DrawRune(fb.Width-1, fb.Height-1, '┘', borderStyle, false)

		titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("yellow"))
		fb.DrawString(2, 2, "Frame Buffer Test", titleStyle, false)

		contentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("green"))
		fb.DrawString(2, 4, "Line 1\nLine 2\nLine 3", contentStyle, false)

		fb.DrawRenderedString(2, 7, "Rendered content")

		result := fb.Render()

		lines := strings.Split(result, "\n")
		assert.Equal(t, fb.Height, len(lines))

		assert.Contains(t, lines[0], "┌")
		assert.Contains(t, lines[0], "┐")
		assert.Contains(t, lines[fb.Height-1], "└")
		assert.Contains(t, lines[fb.Height-1], "┘")

		assert.Contains(t, result, "Frame Buffer Test")
		assert.Contains(t, result, "Line 1")
		assert.Contains(t, result, "Line 2")
		assert.Contains(t, result, "Line 3")
		assert.Contains(t, result, "Rendered content")
	})
}
