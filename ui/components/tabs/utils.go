package tabs

import "github.com/charmbracelet/lipgloss/v2"

func tabBorderWithCustomBottom(left, middle, right string, borderType lipgloss.Border) lipgloss.Border {
	border := borderType
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

type Option func(*Model)

func WithBorder(b lipgloss.Border) Option {
	return func(m *Model) {
		m.border = b
	}
}

func WithActiveStyle(s lipgloss.Style) Option {
	return func(m *Model) {
		m.SetActiveStyle(s)
	}
}

func WithInactiveStyle(s lipgloss.Style) Option {
	return func(m *Model) {
		m.SetInactiveStyle(s)
	}
}
