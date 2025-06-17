package common

import (
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
)

const (
	title   = "ERROR"
	helpMsg = "Press any key to exit"
)

var (
	errorTitleStyle = TitleStyle.
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("1"))

	errMsgStyle = NormalStyle.
			MarginTop(1).
			MarginBottom(2)
)

func ErrorView(err error, width int) string {
	var b strings.Builder

	// Render the title
	rTitle := errorTitleStyle.
		Width(width).
		Render(title)

	// Render the error message
	rErr := errMsgStyle.
		Width(width).
		Render(err.Error())

	// Render the help message
	rHelp := HelpStyle.
		Width(width).
		Render(helpMsg)

	b.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		rTitle,
		rErr,
		rHelp,
	))

	return b.String()
}
