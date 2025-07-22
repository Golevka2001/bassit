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

	// Priority: mute > slap > normal
	// It means that if both are enabled, the sample of mute type will be played
	// And normal mode is the default mode
	muteEnabled bool // `\|` key
	slapEnabled bool // `space` key

	pluckType bass.PluckType
}

func NewFreePlayTabModel(ctm *common.CommonTabModel) FreePlayTabModel {
	fb := fretboard.New(
		ctm.Width,
		ctm.Height,
		ctm.Context.Theme,
		ctm.Context.Bass.GetBaseNotes(),
		false,
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
		// Key event for enabling slap mode or mute mode
		switch msg.String() {
		case config.MuteModeKey:
			m.muteEnabled = true
		case config.SlapModeKey:
			m.slapEnabled = true
		}

		if m.muteEnabled {
			m.pluckType = bass.PluckTypeMute1
		} else if m.slapEnabled {
			m.pluckType = bass.PluckTypeSlap1
		} else {
			m.pluckType = bass.PluckTypeNormal1
		}

		// Key event for plucking a string
		if pluckInfo, ok := config.KeyToPluckInfo[msg.String()]; ok {
			pos := bass.FretboardPosition{
				StringIdx: pluckInfo.StringIdx,
				FretIdx:   m.bass.GetValidFretIdxOfString(pluckInfo.StringIdx),
			}
			var pluckType bass.PluckType
			if m.muteEnabled {
				pluckType = bass.PluckTypeMute1
			} else if m.slapEnabled {
				pluckType = bass.PluckTypeSlap1
			} else {
				pluckType = bass.PluckTypeNormal1
			}
			pluckType += bass.PluckType(pluckInfo.Type)
			pluckTime := m.bass.PluckString(pluckInfo.StringIdx, pluckType)
			// Play the note
			go m.audio.PlayBassNote(pos, pluckType) // avoid blocking the main thread
			cmds = append(cmds, func() tea.Msg {
				return fretboard.PluckStringMsg{
					FretboardPosition: pos,
					Type:              pluckType,
				}
			})
			// Restore string after duration
			cmds = append(cmds, tea.Tick(utils.GetVibDuration(), func(t time.Time) tea.Msg {
				// Only restore if no new pluck has been made
				if m.bass.StopVibratingString(pluckInfo.StringIdx, pluckTime) {
					go m.audio.StopBassNote(pos)
					return fretboard.RestorePluckedStringMsg(pos)
				}
				return nil
			}))
		}

		// Key event for pressing a fret
		if pos, ok := config.KeyToFretboardPos[msg.String()]; ok {
			seqCmds := []tea.Cmd{}
			// If the fret is higher than the valid fret, stop the vibrating string
			lastValidPos := bass.FretboardPosition{
				StringIdx: pos.StringIdx,
				FretIdx:   m.bass.GetValidFretIdxOfString(pos.StringIdx),
			}
			if m.bass.IsStringVibrating(pos.StringIdx) && pos.FretIdx > lastValidPos.FretIdx {
				go m.audio.StopBassNote(lastValidPos)
				m.bass.StopVibratingStringWithoutCheck(pos.StringIdx)
				seqCmds = append(seqCmds, func() tea.Msg {
					return fretboard.RestorePluckedStringMsg(lastValidPos)
				})
			}
			// Then press the fret
			m.bass.PressFret(pos)
			seqCmds = append(seqCmds, func() tea.Msg {
				return fretboard.PressFretMsg(pos)
			})
			cmds = append(cmds, tea.Sequence(seqCmds...))
		}

	case tea.KeyReleaseMsg:
		// Key event for disabling slap mode or mute mode
		switch msg.String() {
		case config.MuteModeKey:
			m.muteEnabled = false
		case config.SlapModeKey:
			m.slapEnabled = false
		}

		if m.muteEnabled {
			m.pluckType = bass.PluckTypeMute1
		} else if m.slapEnabled {
			m.pluckType = bass.PluckTypeSlap1
		} else {
			m.pluckType = bass.PluckTypeNormal1
		}

		// Key event for releasing a fret
		if pos, ok := config.KeyToFretboardPos[msg.String()]; ok {
			seqCmds := []tea.Cmd{}
			// If the fret is higher than the valid fret, stop the vibrating string

			if m.bass.IsStringVibrating(pos.StringIdx) && pos.FretIdx >= m.bass.GetValidFretIdxOfString(pos.StringIdx) {
				go m.audio.StopBassNote(pos)
				m.bass.StopVibratingStringWithoutCheck(pos.StringIdx)
				seqCmds = append(seqCmds, func() tea.Msg {
					return fretboard.RestorePluckedStringMsg(pos)
				})
			}
			// Then release the fret
			m.bass.ReleaseFret(pos)
			seqCmds = append(seqCmds, func() tea.Msg {
				return fretboard.ReleaseFretMsg(pos)
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
