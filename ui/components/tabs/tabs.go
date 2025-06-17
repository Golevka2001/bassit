package tabs

import (
	"github.com/Golevka2001/bassit/ui/common"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

const (
	minTabLabelWidth = 2
)

type Model struct {
	Width  int
	Height int

	border        lipgloss.Border
	activeStyle   lipgloss.Style
	inactiveStyle lipgloss.Style

	activeIdx int
	labels    []string
}

func New(labels []string, activeIdx int, opts ...Option) Model {
	if activeIdx < 0 {
		activeIdx = 0
	} else if activeIdx >= len(labels) {
		activeIdx = max(len(labels)-1, 0)
	}

	m := Model{
		activeIdx: activeIdx,
		labels:    labels,
		border:    lipgloss.RoundedBorder(),
	}

	// Set default styles
	m.SetActiveStyle(common.NormalStyle)
	m.SetInactiveStyle(common.NormalStyle)

	// Apply options
	for _, opt := range opts {
		opt(&m)
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if len(m.labels) == 0 {
		return m, nil
	}

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case SwitchToNextTabMsg:
		m.activeIdx = min(len(m.labels)-1, m.activeIdx+1)
		cmds = append(cmds, nil)

	case SwitchToPreviousTabMsg:
		m.activeIdx = max(0, m.activeIdx-1)
		cmds = append(cmds, nil)

	case SwitchToTabMsg:
		if msg.TabIdx < 0 {
			msg.TabIdx = 0
		} else if msg.TabIdx >= len(m.labels) {
			msg.TabIdx = len(m.labels) - 1
		}
		m.activeIdx = msg.TabIdx
		cmds = append(cmds, nil)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.activeIdx >= len(m.labels) {
		m.activeIdx = max(0, len(m.labels)-1)
	}

	tabCnt := len(m.labels)
	maxWidth := max(minTabLabelWidth, m.Width/tabCnt-4) // `4` for padding
	var renderedTabs []string
	var tabsWidthTotal int

	for i, tab := range m.labels {
		isFirst := i == 0
		isLast := i == tabCnt-1
		isActive := i == m.activeIdx

		style := m.inactiveStyle
		if isActive {
			style = m.activeStyle
		}

		border, _, _, _, _ := style.GetBorder()
		switch {
		case isFirst && isActive:
			border.BottomLeft = m.border.Left
		case isFirst:
			border.BottomLeft = m.border.MiddleLeft
		}
		switch {
		case isLast && isActive:
			border.BottomRight = m.border.Right
		case isLast:
			border.BottomRight = m.border.MiddleRight
		}
		style = style.Border(border, true)

		tabLabel := tab
		if lipgloss.Width(tab) > maxWidth {
			if maxWidth > 1 {
				var truncated string
				widthAccum := 0
				for _, r := range tab {
					w := lipgloss.Width(string(r))
					if widthAccum+w >= maxWidth {
						break
					}
					truncated += string(r)
					widthAccum += w
				}
				tabLabel = truncated + "…"
			} else {
				tabLabel = "…"
			}
		}

		rendered := style.Render(tabLabel)
		tabsWidthTotal += lipgloss.Width(rendered)

		if isLast && tabsWidthTotal < m.Width {
			tabsWidthTotal -= lipgloss.Width(rendered)
			if isActive {
				border.BottomRight = m.border.BottomLeft
			} else {
				border.BottomRight = m.border.MiddleBottom
			}
			style = style.Border(border, true)
			rendered = style.Render(tabLabel)
			tabsWidthTotal += lipgloss.Width(rendered)
		}

		renderedTabs = append(renderedTabs, rendered)
	}

	remainingWidth := m.Width - tabsWidthTotal
	if remainingWidth > 0 {
		border := lipgloss.HiddenBorder()
		border.BottomLeft = m.border.Bottom
		border.Bottom = m.border.Bottom
		border.BottomRight = m.border.TopRight
		if remainingWidth == 1 {
			border.BottomLeft = m.border.TopRight
			border.Bottom = ""
			border.BottomRight = ""
		}
		fillerStyle := m.inactiveStyle.
			Border(border).
			UnsetPadding().
			UnsetForeground().
			UnsetBackground().
			Width(remainingWidth)
		renderedTabs = append(renderedTabs, fillerStyle.Render(""))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
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
