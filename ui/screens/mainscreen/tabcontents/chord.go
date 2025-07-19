package tabcontents

import (
	"strings"

	"github.com/Golevka2001/bassit/audio"
	"github.com/Golevka2001/bassit/bass"
	"github.com/Golevka2001/bassit/config"
	"github.com/Golevka2001/bassit/ui/common"
	"github.com/Golevka2001/bassit/ui/components/fretboard"
	"github.com/Golevka2001/bassit/utils"

	detector "github.com/Golevka2001/go-chord-detector"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/go-music-theory/music-theory/note"
)

var (
	notesStyle  = common.NormalStyle.MarginTop(1)
	chordsStyle = common.BUTextStyle.Foreground(lipgloss.Blue).MarginBottom(1)
)

type ChordTabModel struct {
	commonModel *common.CommonTabModel

	bass      *bass.BassModel
	audio     *audio.AudioManager
	fretboard fretboard.Model
}

func NewChordTabModel(ctm *common.CommonTabModel) ChordTabModel {
	fb := fretboard.New(
		ctm.Width,
		ctm.Height,
		ctm.Context.Theme,
		ctm.Context.Bass.GetBaseNotes(),
	)

	return ChordTabModel{
		commonModel: ctm,
		bass:        ctm.Context.Bass,
		audio:       ctm.Context.Audio,
		fretboard:   fb,
	}
}

func (m ChordTabModel) Init() tea.Cmd {
	return tea.Batch(
		m.fretboard.Init(),
	)
}

func (m ChordTabModel) Update(msg tea.Msg) (ChordTabModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// Key event for pressing a fret
		if pos, ok := config.KeyToFretboardPos[msg.String()]; ok {
			// Switch fret pressed state
			if m.bass.IsFretPressed(pos) {
				// If the fret is already pressed, release it
				m.bass.ReleaseFret(pos)
				cmds = append(cmds, func() tea.Msg {
					return fretboard.ReleaseFretMsg(pos)
				})
			} else {
				// If the fret is not pressed, press it
				m.bass.PressFret(pos)
				cmds = append(cmds, func() tea.Msg {
					return fretboard.PressFretMsg(pos)
				})
			}
		}
		if pluckInfo, ok := config.KeyToPluckInfo[msg.String()]; ok {
			stringIdx := pluckInfo.StringIdx
			cmds = append(cmds, func() tea.Msg {
				return fretboard.PressFretMsg(bass.FretboardPosition{
					StringIdx: stringIdx,
					FretIdx:   0,
				})
			})
		}
	}

	var cmd tea.Cmd
	m.fretboard, cmd = m.fretboard.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ChordTabModel) View() string {
	var b strings.Builder

	// Render the notes
	notes := []*note.Note{}
	var notesStr strings.Builder
	for stringIdx := range config.StringCnt {
		n := m.bass.GetNoteToPlay(config.StringCnt - stringIdx - 1) // lowest -> highest
		if n == nil {
			notesStr.WriteString("  X  ")
		} else {
			notesStr.WriteString("  " + utils.GetNoteNameWithOctave(*n) + "  ")
			notes = append(notes, n)
		}
	}
	rNotes := notesStyle.
		Width(m.commonModel.Width).
		Align(lipgloss.Center).
		Render(notesStr.String())

	// Render the down arrow
	rDownArrow := common.NormalStyle.
		Width(m.commonModel.Width).
		Align(lipgloss.Center).
		Render("↓")

	// Render the chord
	var chordsStr strings.Builder
	chords := detector.Detect(notes)
	if len(chords) > 0 {
		chordsStr.WriteString(strings.Join(chords, "  |  "))
	} else {
		chordsStr.WriteString("Unknown")
	}
	rChords := chordsStyle.
		Width(m.commonModel.Width).
		Align(lipgloss.Center).
		Render(chordsStr.String())

	// Render the fretboard
	rFretboard := common.NormalStyle.
		Width(m.commonModel.Width).
		Height(m.commonModel.Height - 5).
		AlignVertical(lipgloss.Center).
		Render(m.fretboard.View())

	b.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		rNotes,
		rDownArrow,
		rChords,
		rFretboard,
	))

	return b.String()
}

func (m *ChordTabModel) SyncSize() {
	m.fretboard.SetSize(
		m.commonModel.Width,
		m.commonModel.Height-5,
	)
}
