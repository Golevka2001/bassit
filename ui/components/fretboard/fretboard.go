package fretboard

import (
	"github.com/Golevka2001/bassit/config"
	"github.com/Golevka2001/bassit/ui/common"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/go-music-theory/music-theory/note"
)

type Model struct {
	Width  int
	Height int

	frameBuf common.FrameBuffer

	baseNotes []note.Note
	style     Style
	theme     *config.Theme

	layoutMappingsUpdated bool
	fretboardLen          int
	fretboardHeight       int
	nutStartX             int
	nutEndX               int
	fretboardStartX       int
	fretboardEndX         int
	fretwireX             []int
	xToFretwire           map[int]int
	fretCenterX           []int
	stringY               []int
	yToString             map[int]int
}

func New(width, height int, theme *config.Theme, baseNotes []note.Note) Model {
	m := Model{
		Width:     width,
		Height:    height,
		frameBuf:  common.NewFrameBuffer(width, height),
		baseNotes: baseNotes,
		style:     NewStyle(theme),
		theme:     theme,
	}

	m.updateLayoutMappings()

	return m
}

func (m Model) Init() tea.Cmd {
	// Draw the fretboard
	m.drawFretboard()
	m.drawFretboardInlays()
	m.drawStrings()

	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case PressFretMsg:
		m.drawPressedFret(msg.StringIdx, msg.FretIdx)

	case ReleaseFretMsg:
		m.restorePressedFret(msg.StringIdx, msg.FretIdx)

	case PluckStringMsg:
		m.drawVibratingString(msg.StringIdx, msg.FretIdx, int(msg.Type)%2)

	case RestorePluckedStringMsg:
		m.restoreVibratingString(msg.StringIdx, msg.FretIdx)
	}

	return m, nil
}

func (m Model) View() string {
	if m.Width <= 0 || m.Height <= 0 {
		return ""
	}

	if m.frameBuf.Empty() {
		m.drawFretboard()
		m.drawFretboardInlays()
		m.drawStrings()
	}

	return m.frameBuf.Render()
}

func (m *Model) SetSize(width, height int) {
	width = max(0, width)
	height = max(0, height)

	m.Width = width
	m.Height = height

	m.layoutMappingsUpdated = false
	m.updateLayoutMappings()
	m.frameBuf.Clear()
	m.frameBuf = common.NewFrameBuffer(
		width,
		min(height, m.fretboardHeight),
	)
	m.drawFretboard()
	m.drawFretboardInlays()
	m.drawStrings()
}
