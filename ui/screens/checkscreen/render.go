package checkscreen

import (
	"strings"

	"github.com/Golevka2001/bassit/ui/common"

	"github.com/charmbracelet/lipgloss/v2"
)

const (
	// minGap is the minimum gap between the item and current step/info text
	minGap = 2

	pendingMark = "-"
	passedMark  = "✓"
	infoMark    = "i"
	failedMark  = "✗"
)

var (
	pendingColor    = lipgloss.Color("8")
	inProgressColor = lipgloss.Color("3")
	passedColor     = lipgloss.Color("2")
	infoColor       = lipgloss.Color("6")
	failedColor     = lipgloss.Color("1")

	pendingStyle    = common.NormalStyle.Foreground(pendingColor)
	inProgressStyle = common.NormalStyle.Foreground(inProgressColor)
	passedStyle     = common.NormalStyle.Foreground(passedColor)
	infoStyle       = common.ITextStyle.Foreground(infoColor)
	failedStyle     = common.BITextStyle.Foreground(failedColor)
)

func (m *Model) renderChecklist() string {
	var renderedItems []string

	for _, item := range m.checkList {
		var style lipgloss.Style
		var marker string
		var curStep string

		switch item.status {
		case statusPending:
			style = pendingStyle
			marker = pendingMark
		case statusInProgress:
			style = inProgressStyle
			marker = m.spinner.View()
			curStep = item.curStep
		case statusPassed:
			style = passedStyle
			marker = passedMark
		case statusInfo:
			style = infoStyle
			marker = infoMark
			curStep = item.curStep
		case statusFailed:
			style = failedStyle
			marker = failedMark
			curStep = item.curStep
		}
		curStep = strings.TrimSpace(curStep)
		if len(curStep) > 0 {
			curStep = "(" + curStep + ")"
		}

		line := marker + " " + item.title
		if curStep != "" {
			cellsRemaining := max(minGap, m.commonModel.Width-lipgloss.Width(line+curStep))
			gap := strings.Repeat(" ", cellsRemaining)
			line += gap + curStep
		}
		renderedItems = append(renderedItems, style.Width(m.commonModel.Width).Render(line))
	}

	return lipgloss.JoinVertical(lipgloss.Left, renderedItems...)
}
