package tabselector

import (
	"github.com/Golevka2001/bassit/ui/common"

	"github.com/charmbracelet/lipgloss/v2"
)

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

func (m *Model) SetWidth(width int) {
	width = max(0, width)
	m.Width = width
}

func (m *Model) SetBorder(border lipgloss.Border) {
	m.border = border
	m.SetActiveStyle(common.NormalStyle)
	m.SetInactiveStyle(common.NormalStyle)
}

func (m *Model) SetActiveStyle(style lipgloss.Style) {
	border := tabBorderWithCustomBottom(
		m.border.BottomRight,
		" ",
		m.border.BottomLeft,
		m.border,
	)

	m.activeStyle = common.NormalStyle.
		Inherit(style).
		Padding(0, 1).
		Border(border, true)
}

func (m *Model) SetInactiveStyle(style lipgloss.Style) {
	border := tabBorderWithCustomBottom(
		m.border.MiddleBottom,
		m.border.Bottom,
		m.border.MiddleBottom,
		m.border,
	)

	m.inactiveStyle = common.NormalStyle.
		Inherit(style).
		Padding(0, 1).
		Border(border, true)
}

func (m *Model) GetFocusedIdx() int {
	return m.focusedIdx
}
