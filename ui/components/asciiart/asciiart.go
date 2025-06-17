package asciiart

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type Model struct {
	Width  int
	Height int

	content string
}

func New(asciiArt string) Model {
	return Model{
		content: asciiArt,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	rows := strings.Split(m.content, "\n")
	artHeight := len(rows)

	// Calculate max width of ASCII art
	var artWidth int
	for _, row := range rows {
		artWidth = max(artWidth, lipgloss.Width(row))
	}

	// Vertically center
	if m.Height <= 0 {
		return ""
	}
	if artHeight > m.Height {
		vStart := (artHeight - m.Height) / 2
		vEnd := vStart + m.Height
		rows = rows[vStart:vEnd]
	} else {
		paddingTop := (m.Height - artHeight) / 2
		paddingBottom := m.Height - artHeight - paddingTop
		rows = append(make([]string, paddingTop), rows...)
		rows = append(rows, make([]string, paddingBottom)...)
	}

	// Horizontally center
	for i, line := range rows {
		lineWidth := lipgloss.Width(line)
		switch {
		case lineWidth <= m.Width:
			padding := (m.Width - lineWidth) / 2
			rows[i] = strings.Repeat(" ", padding) + line
		default:
			rows[i] = cropLineToWidth(line, m.Width)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m *Model) SetSize(width, height int) {
	width = max(0, width)
	height = max(0, height)

	m.Width = width
	m.Height = height
}

func (m *Model) GetContent() string {
	return m.content
}

func (m *Model) SetContent(asciiArt string) {
	m.content = asciiArt
}
