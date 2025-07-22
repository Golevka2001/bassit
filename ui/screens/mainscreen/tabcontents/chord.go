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

// staticElementsHeight represents the height of the notes and chord detection result
const staticElementsHeight = 5

type ChordTabModel struct {
	commonModel *common.CommonTabModel

	bass      *bass.BassModel
	audio     *audio.AudioManager
	fretboard fretboard.Model

	ignoredStrings [config.StringCnt]bool
}

func NewChordTabModel(ctm *common.CommonTabModel) ChordTabModel {
	fb := fretboard.New(
		ctm.Width,
		ctm.Height,
		ctm.Context.Theme,
		ctm.Context.Bass.GetBaseNotes(),
		true,
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
				go m.audio.StopBassNote(pos)
				cmds = append(cmds, func() tea.Msg {
					return fretboard.ReleaseFretMsg(pos)
				})
			} else {
				// If the fret is not pressed, press it
				m.bass.PressFret(pos)
				// Play the note
				go m.audio.PlayBassNote(pos, bass.PluckTypeNormal1)
				cmds = append(cmds, func() tea.Msg {
					return fretboard.PressFretMsg(pos)
				})
			}
		}
		// Key event for ignoring a string
		if pluckInfo, ok := config.KeyToPluckInfo[msg.String()]; ok {
			stringIdx := pluckInfo.StringIdx
			if m.ignoredStrings[stringIdx] {
				// If the string is already ignored, restore it
				m.ignoredStrings[stringIdx] = false
				// Play the note
				pos := bass.FretboardPosition{
					StringIdx: pluckInfo.StringIdx,
					FretIdx:   m.bass.GetValidFretIdxOfString(pluckInfo.StringIdx),
				}
				go m.audio.PlayBassNote(pos, bass.PluckTypeNormal1)
				cmds = append(cmds, func() tea.Msg {
					return fretboard.RestoreIgnoredStringMsg(stringIdx)
				})
			} else {
				// If the string is not ignored, ignore it
				m.ignoredStrings[stringIdx] = true
				// Stop all players on the string
				for fretIdx := range config.DisplayedFretCount {
					pos := bass.FretboardPosition{StringIdx: stringIdx, FretIdx: fretIdx}
					go m.audio.StopBassNote(pos)
				}

				cmds = append(cmds, func() tea.Msg {
					return fretboard.IgnoreStringMsg(stringIdx)
				})
			}
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
		// lowest -> highest
		stringIdx = config.StringCnt - stringIdx - 1
		n := m.bass.GetNoteToPlay(stringIdx)
		// If the string is ignored, skip it
		if m.ignoredStrings[stringIdx] || n == nil {
			notesStr.WriteString("  X  ")
			continue
		}
		notesStr.WriteString("  " + utils.GetNoteNameWithOctave(*n) + "  ")
		notes = append(notes, n)
	}
	rNotes := common.NormalStyle.
		MarginTop(1).
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
	rChords := common.BUTextStyle.
		MarginBottom(1).
		Width(m.commonModel.Width).
		Align(lipgloss.Center).
		Foreground(lipgloss.Blue).
		Render(chordsStr.String())

	// Render the fretboard
	rFretboard := common.NormalStyle.
		Width(m.commonModel.Width).
		Height(m.commonModel.Height - staticElementsHeight).
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
		m.commonModel.Height-staticElementsHeight,
	)
}
