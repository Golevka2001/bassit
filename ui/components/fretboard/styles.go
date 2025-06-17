package fretboard

import (
	"github.com/Golevka2001/bassit/config"
	"github.com/Golevka2001/bassit/ui/common"
	
	"github.com/charmbracelet/lipgloss/v2"
)

type Style struct {
	fretboardBorderStyle   lipgloss.Style
	fretboardBgStyle       lipgloss.Style
	nutBorderStyle         lipgloss.Style
	nutBgStyle             lipgloss.Style
	fretwireStyle          lipgloss.Style
	baseNoteStyle          lipgloss.Style
	stringStyle            lipgloss.Style
	dotInlayStyle          lipgloss.Style
	blockInlayStyle        lipgloss.Style
	pressedFretSignStyle   lipgloss.Style
	pluckedStringSignStyle lipgloss.Style
}

func NewStyle(theme *config.Theme) Style {
	ns := common.NormalStyle
	s := Style{}
	s.fretboardBorderStyle = ns.Foreground(theme.FretboardBorderColor)
	s.fretboardBgStyle = ns.Background(theme.FretboardBgColor)
	s.nutBorderStyle = ns.Foreground(theme.NutBorderColor)
	s.nutBgStyle = ns.Background(theme.NutBgColor)
	s.fretwireStyle = ns.Foreground(theme.FretwireColor)
	s.baseNoteStyle = ns
	s.stringStyle = ns.Foreground(theme.StringColor)
	s.dotInlayStyle = ns.Foreground(theme.InlayColor)
	s.blockInlayStyle = ns.Background(theme.InlayColor)
	s.pressedFretSignStyle = ns.Foreground(theme.PressedFretSignColor)
	s.pluckedStringSignStyle = ns.Foreground(theme.PluckedStringSignColor)
	return s
}
