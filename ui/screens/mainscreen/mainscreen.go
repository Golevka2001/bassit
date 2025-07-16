package mainscreen

import (
	"reflect"
	"strings"

	"github.com/Golevka2001/bassit/ui/common"
	"github.com/Golevka2001/bassit/ui/components/tabs"
	"github.com/Golevka2001/bassit/ui/screens/mainscreen/tabcontents"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// staticElementsHeight represents the height of the title, tabs and borders
// TODO: make it dynamic
const staticElementsHeight = 6

// tabState is an enum for the different tabs in the program
type tabState int

const (
	tabStateFreePlay tabState = iota
	tabStateSettings
	tabStateExit
)

func (t tabState) String() string {
	return map[tabState]string{
		tabStateFreePlay: "Free Play",
		tabStateSettings: "Settings",
		tabStateExit:     "Exit",
	}[t]
}

var (
	tabLabels = []string{
		tabStateFreePlay.String(),
		tabStateSettings.String(),
		tabStateExit.String(),
	}

	tabContentStyle = common.NormalStyle.Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			UnsetBorderTop()
)

type Model struct {
	commonScreenModel *common.CommonScreenModel
	commonTabModel    *common.CommonTabModel

	tabs  tabs.Model
	state tabState

	focusOnTab bool

	// Sub-models
	freePlayTab tabcontents.FreePlayTabModel
	settingsTab tabcontents.SettingsTabModel
	exitTab     tabcontents.ExitTabModel
}

func NewModel(csm *common.CommonScreenModel) Model {
	ctm := common.CommonTabModel{
		Context: csm.Context,
	}

	m := Model{
		commonScreenModel: csm,
		commonTabModel:    &ctm,
		tabs: tabs.New(
			tabLabels,
			0,
			tabs.WithBorder(lipgloss.RoundedBorder()),
		),
		state:       tabStateFreePlay,
		focusOnTab:  false,
		freePlayTab: tabcontents.NewFreePlayTabModel(&ctm),
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
	cmds = append(cmds, m.tabs.Init())

	switch m.state {
	case tabStateFreePlay:
		cmds = append(cmds, m.freePlayTab.Init())
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
					return tabs.FocusMsg{TabIdx: int(m.state)}
				})
			} else {
				cmds = append(cmds, func() tea.Msg {
					return tabs.UnfocusMsg{}
				})
			}

		case "enter":
			if m.focusOnTab {
				focusedIdx := m.tabs.GetFocusedIdx()
				if int(m.state) != focusedIdx {
					// Switch to the focused tab
					m.focusOnTab = false
					m.state = tabState(focusedIdx)
					cmds = append(cmds, func() tea.Msg {
						return tabs.SwitchToTabMsg{TabIdx: focusedIdx}
					})
				}
			}
		}
	}

	if m.focusOnTab ||
		reflect.TypeOf(msg) == reflect.TypeOf(tabs.FocusMsg{}) ||
		reflect.TypeOf(msg) == reflect.TypeOf(tabs.UnfocusMsg{}) ||
		reflect.TypeOf(msg) == reflect.TypeOf(tabs.SwitchToTabMsg{}) {
		m.tabs, cmd = m.tabs.Update(msg)
		cmds = append(cmds, cmd)
	}

	if !m.focusOnTab {
		switch m.state {
		case tabStateFreePlay:
			m.freePlayTab, cmd = m.freePlayTab.Update(msg)
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
		Render(m.tabs.View())

	// Render the active tab content
	var tabContent string
	switch m.state {
	case tabStateFreePlay:
		tabContent = m.freePlayTab.View()
	case tabStateSettings:
		tabContent = m.settingsTab.View()
	case tabStateExit:
		tabContent = m.exitTab.View()
	}
	rTabContent := tabContentStyle.Render(tabContent)

	b.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		rTitle,
		rTabs,
		rTabContent,
	))

	return b.String()
}

func (m *Model) SyncSize() {
	m.tabs.SetWidth(m.commonScreenModel.Width)

	m.commonTabModel.Width = m.commonScreenModel.Width - 4 // `4` for padding and borders
	m.commonTabModel.Height = m.commonScreenModel.Height - staticElementsHeight
	m.freePlayTab.SyncSize()
	m.settingsTab.SyncSize()
	m.exitTab.SyncSize()
}
