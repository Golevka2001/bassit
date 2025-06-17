package tabcontents

import (
	"time"

	"github.com/Golevka2001/bassit/audio"
	"github.com/Golevka2001/bassit/bass"
	"github.com/Golevka2001/bassit/config"
	"github.com/Golevka2001/bassit/ui/common"
	"github.com/Golevka2001/bassit/ui/components/fretboard"
	"github.com/Golevka2001/bassit/utils"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type FreePlayTabModel struct {
	commonModel *common.CommonTabModel

	bass      *bass.BassModel
	audio     *audio.AudioManager
	fretboard fretboard.Model
}

func NewFreePlayTabModel(ctm *common.CommonTabModel) FreePlayTabModel {
	fb := fretboard.New(
		ctm.Width,
		ctm.Height,
		ctm.Context.Theme,
		ctm.Context.Bass.GetBaseNotes(),
	)

	return FreePlayTabModel{
		commonModel: ctm,
		bass:        ctm.Context.Bass,
		audio:       ctm.Context.Audio,
		fretboard:   fb,
	}
}

func (m FreePlayTabModel) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestKeyboardEnhancements(tea.WithKeyReleases),
		m.fretboard.Init(),
	)
}

func (m FreePlayTabModel) Update(msg tea.Msg) (FreePlayTabModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// Key event for plucking a string
		if stringIdx, ok := config.KeyToStringIdx[msg.String()]; ok {
			curString := m.bass.Strings[stringIdx]
			curFretIdx := curString.CurValidFret
			curNote := curString.GetNoteToPlay()

			m.bass.PluckString(stringIdx)
			go m.audio.PlayBassNote(curNote)
			cmds = append(cmds, func() tea.Msg {
				return fretboard.PluckStringMsg{
					StringIdx: stringIdx,
					FretIdx:   curFretIdx,
				}
			})
			// Restore string after duration
			cmds = append(cmds, tea.Tick(utils.GetVibDuration(), func(t time.Time) tea.Msg {
				m.audio.StopBassNote(curNote)
				m.bass.RestoreString(stringIdx)
				return fretboard.RestorePluckedStringMsg{
					StringIdx: stringIdx,
					FretIdx:   curFretIdx,
				}
			}))
		}

		// Key event for pressing a fret
		if pos, ok := config.KeyToFretboardPos[msg.String()]; ok {
			seqCmds := []tea.Cmd{}
			// If the fret is higher than the valid fret, stop the vibrating string
			curString := m.bass.Strings[pos.StringIdx]
			curFretIdx := curString.CurValidFret
			if curString.PluckedState && pos.FretIdx >= curString.CurValidFret {
				m.audio.StopBassNote(curString.GetNoteToPlay())
				m.bass.RestoreString(pos.StringIdx)
				seqCmds = append(seqCmds, func() tea.Msg {
					return fretboard.RestorePluckedStringMsg{
						StringIdx: pos.StringIdx,
						FretIdx:   curFretIdx,
					}
				})
			}
			// Then press the fret
			m.bass.Press(pos.StringIdx, pos.FretIdx)
			seqCmds = append(seqCmds, func() tea.Msg {
				return fretboard.PressFretMsg{
					StringIdx: pos.StringIdx,
					FretIdx:   pos.FretIdx,
				}
			})
			cmds = append(cmds, tea.Sequence(seqCmds...))
		}

	case tea.KeyReleaseMsg:
		// Key event for releasing a fret
		if pos, ok := config.KeyToFretboardPos[msg.String()]; ok {
			seqCmds := []tea.Cmd{}
			// If the fret is higher than the valid fret, stop the vibrating string
			curString := m.bass.Strings[pos.StringIdx]
			curFretIdx := curString.CurValidFret
			if curString.PluckedState && pos.FretIdx >= curString.CurValidFret {
				m.audio.StopBassNote(curString.GetNoteToPlay())
				m.bass.RestoreString(pos.StringIdx)
				seqCmds = append(seqCmds, func() tea.Msg {
					return fretboard.RestorePluckedStringMsg{
						StringIdx: pos.StringIdx,
						FretIdx:   curFretIdx,
					}
				})
			}
			// Then release the fret
			m.bass.Release(pos.StringIdx, pos.FretIdx)
			seqCmds = append(seqCmds, func() tea.Msg {
				return fretboard.ReleaseFretMsg{
					StringIdx: pos.StringIdx,
					FretIdx:   pos.FretIdx,
				}
			})
			cmds = append(cmds, tea.Sequence(seqCmds...))
		}
	}

	var cmd tea.Cmd
	m.fretboard, cmd = m.fretboard.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m FreePlayTabModel) View() string {
	return common.NormalStyle.
		Width(m.commonModel.Width).
		Height(m.commonModel.Height).
		AlignVertical(lipgloss.Center).
		Render(m.fretboard.View())
}

func (m *FreePlayTabModel) SyncSize() {
	m.fretboard.SetSize(
		m.commonModel.Width,
		m.commonModel.Height,
	)
}
