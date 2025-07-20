package mainscreen

import (
	"reflect"
	"strings"

	"github.com/Golevka2001/bassit/ui/common"
	"github.com/Golevka2001/bassit/ui/components/tabselector"
	"github.com/Golevka2001/bassit/ui/screens/mainscreen/tabcontents"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// staticElementsHeight represents the height of the title, tabs and borders
const staticElementsHeight = 6

// tabState is an enum for the different tabs in the program
type tabState int

const (
	tabStateFreePlay tabState = iota
	tabStateChord
	tabStateSettings
	tabStateExit
)

func (t tabState) String() string {
	return map[tabState]string{
		tabStateFreePlay: "Free Play",
		tabStateChord:    "Chord",
		tabStateSettings: "Settings",
		tabStateExit:     "Exit",
	}[t]
}

type Model struct {
	commonScreenModel *common.CommonScreenModel
	commonTabModel    *common.CommonTabModel

	selector   tabselector.Model
	state      tabState
	focusOnTab bool

	// Sub-models
	freePlayTab tabcontents.FreePlayTabModel
	chordTab    tabcontents.ChordTabModel
	settingsTab tabcontents.SettingsTabModel
	exitTab     tabcontents.ExitTabModel
}

func NewModel(csm *common.CommonScreenModel) Model {
	ctm := common.CommonTabModel{
		Context: csm.Context,
	}

	tabLabels := []string{
		tabStateFreePlay.String(),
		tabStateChord.String(),
		tabStateSettings.String(),
		tabStateExit.String(),
	}

	m := Model{
		commonScreenModel: csm,
		commonTabModel:    &ctm,
		selector: tabselector.New(
			tabLabels,
			0,
			tabselector.WithBorder(lipgloss.RoundedBorder()),
		),
		state:       tabStateFreePlay,
		focusOnTab:  false,
		freePlayTab: tabcontents.NewFreePlayTabModel(&ctm),
		chordTab:    tabcontents.NewChordTabModel(&ctm),
		settingsTab: tabcontents.NewSettingsTabModel(&ctm),
		exitTab:     tabcontents.NewExitTabModel(&ctm),
	}

	return m
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds,
		tea.EnterAltScreen,
		tea.ClearScreen,
	)
	cmds = append(cmds, m.selector.Init())

	switch m.state {
	case tabStateFreePlay:
		cmds = append(cmds, m.freePlayTab.Init())
	case tabStateChord:
		cmds = append(cmds, m.chordTab.Init())
	case tabStateSettings:
		cmds = append(cmds, m.settingsTab.Init())
	}

	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			m.focusOnTab = !m.focusOnTab
			cmds = append(cmds, nil)
			if m.focusOnTab {
				cmds = append(cmds, func() tea.Msg {
					return tabselector.FocusMsg{TabIdx: int(m.state)}
				})
			} else {
				cmds = append(cmds, func() tea.Msg {
					return tabselector.UnfocusMsg{}
				})
			}

		case "enter":
			if m.focusOnTab {
				focusedIdx := m.selector.GetFocusedIdx()
				if int(m.state) != focusedIdx {
					// Switch to the focused tab
					m.focusOnTab = false
					m.state = tabState(focusedIdx)
					cmds = append(cmds, func() tea.Msg {
						return tabselector.SwitchToTabMsg{TabIdx: focusedIdx}
					})
				}
			}
		}
	}

	if m.focusOnTab ||
		reflect.TypeOf(msg) == reflect.TypeOf(tabselector.FocusMsg{}) ||
		reflect.TypeOf(msg) == reflect.TypeOf(tabselector.UnfocusMsg{}) ||
		reflect.TypeOf(msg) == reflect.TypeOf(tabselector.SwitchToTabMsg{}) {
		m.selector, cmd = m.selector.Update(msg)
		cmds = append(cmds, cmd)
	}

	if !m.focusOnTab {
		switch m.state {
		case tabStateFreePlay:
			m.freePlayTab, cmd = m.freePlayTab.Update(msg)
			cmds = append(cmds, cmd)
		case tabStateChord:
			m.chordTab, cmd = m.chordTab.Update(msg)
			cmds = append(cmds, cmd)
		case tabStateSettings:
			m.settingsTab, cmd = m.settingsTab.Update(msg)
			cmds = append(cmds, cmd)
		case tabStateExit:
			m.exitTab, cmd = m.exitTab.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var b strings.Builder

	// Render the title
	rTitle := common.TitleStyle.
		Width(m.commonScreenModel.Width).
		Render(common.TitleText)

	// Render the tabs
	rTabs := common.NormalStyle.
		MarginTop(1).
		Render(m.selector.View())

	// Render the active tab content
	var tabContent string
	switch m.state {
	case tabStateFreePlay:
		tabContent = m.freePlayTab.View()
	case tabStateChord:
		tabContent = m.chordTab.View()
	case tabStateSettings:
		tabContent = m.settingsTab.View()
	case tabStateExit:
		tabContent = m.exitTab.View()
	}
	rTabContent := common.NormalStyle.
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		UnsetBorderTop().
		Render(tabContent)

	b.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		rTitle,
		rTabs,
		rTabContent,
	))

	return b.String()
}

func (m *Model) SyncSize() {
	m.selector.SetWidth(m.commonScreenModel.Width)

	m.commonTabModel.Width = m.commonScreenModel.Width - 4 // `4` for padding and borders
	m.commonTabModel.Height = m.commonScreenModel.Height - staticElementsHeight
	m.freePlayTab.SyncSize()
	m.chordTab.SyncSize()
	m.settingsTab.SyncSize()
	m.exitTab.SyncSize()
}
