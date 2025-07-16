package common

import (
	"fmt"

	"github.com/Golevka2001/bassit/config"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"
)

const (
	// VPad defines the vertical padding
	VPad = 1
	// HPad defines the horizontal padding
	HPad = 2
)

var (
	TitleText    = fmt.Sprintf("𝄢 BASSIT v%s", config.Version)
	TitleFgColor = lipgloss.Color("0")
	TitleBgColor = lipgloss.Color("15")
)

var (
	NormalStyle = lipgloss.NewStyle()

	ScreenStyle = NormalStyle.Padding(VPad, HPad)

	BTextStyle  = NormalStyle.Bold(true)
	ITextStyle  = NormalStyle.Italic(true)
	UTextStyle  = NormalStyle.Underline(true)
	BITextStyle = BTextStyle.Italic(true)
	UITextStyle = UTextStyle.Bold(true)

	TitleStyle = BTextStyle.
			Inline(true).
			AlignHorizontal(lipgloss.Center).
			Foreground(TitleFgColor).
			Background(TitleBgColor)

	// HelpStyle is copied from `bubbles/help` component
	HelpStyle = NormalStyle.Foreground(compat.AdaptiveColor{
		Light: lipgloss.Color("#B2B2B2"),
		Dark:  lipgloss.Color("#4A4A4A"),
	})
)
