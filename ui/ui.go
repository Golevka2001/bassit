package ui

import (
	"fmt"

	"github.com/Golevka2001/bassit/ui/common"
	"github.com/Golevka2001/bassit/ui/screens/mainscreen"
	"github.com/Golevka2001/bassit/ui/screens/welcomescreen"

	tea "github.com/charmbracelet/bubbletea/v2"
)

// screenState is an enum for the different screens in the program
type screenState int

const (
	stateWelcomeScreen screenState = iota
	stateMainScreen
)

func (s screenState) String() string {
	return map[screenState]string{
		stateWelcomeScreen: "Welcome Screen",
		stateMainScreen:    "Main Screen",
	}[s]
}

type Model struct {
	width  int
	height int

	commonModel *common.CommonScreenModel
	state       screenState
	fatalErr    error

	// Sub-models
	welcome welcomescreen.Model
	main    mainscreen.Model
}

func NewProgram(ctx *common.UIContext) *tea.Program {
	model := NewModel(ctx)

	program := tea.NewProgram(model,
		tea.WithAltScreen(),
		tea.WithKeyboardEnhancements(tea.WithKeyReleases),
	)

	return program
}

func NewModel(ctx *common.UIContext) tea.Model {
	csm := common.CommonScreenModel{
		Context: ctx,
	}

	m := Model{
		commonModel: &csm,
		state:       stateWelcomeScreen,
		welcome:     welcomescreen.NewModel(&csm),
		main:        mainscreen.NewModel(&csm),
	}

	return m
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	switch m.state {
	case stateWelcomeScreen:
		cmds = append(cmds, m.welcome.Init())
	case stateMainScreen:
		cmds = append(cmds, m.main.Init())
	default:
		// Default to the welcome screen if the state is unknown
		cmds = append(cmds, m.welcome.Init())
	}

	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If there's been an error, any key exits
	if m.fatalErr != nil {
		if _, ok := msg.(tea.KeyMsg); ok {
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case common.ErrMsg:
		m.fatalErr = msg
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+z":
			return m, tea.Suspend

		// Ctrl+C always quits no matter where in the application you are
		case "ctrl+c":
			return m, tea.Quit
		}

	// Window size is received when starting up and on every resize
	case tea.WindowSizeMsg:
		m.width = max(0, msg.Width)
		m.height = max(0, msg.Height)
		m.commonModel.Width = m.width - common.HPad*2
		m.commonModel.Height = m.height - common.VPad*2

		// Sync the size of the sub-models
		m.welcome.SyncSize()
		m.main.SyncSize()

	case common.SwitchScreenMsg:
		switch m.state {
		case stateWelcomeScreen:
			m.state = stateMainScreen
			cmds = append(cmds, m.main.Init())
		}
	}

	// Process children
	switch m.state {
	case stateWelcomeScreen:
		newWelcomeScreenModel, cmd := m.welcome.Update(msg)
		m.welcome = newWelcomeScreenModel
		cmds = append(cmds, cmd)

	case stateMainScreen:
		newMainScreenModel, cmd := m.main.Update(msg)
		m.main = newMainScreenModel
		cmds = append(cmds, cmd)

	default:
		m.fatalErr = fmt.Errorf("internal error: unknown screen state: %d", m.state)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.fatalErr != nil {
		return common.ScreenStyle.
			Render(common.ErrorView(m.fatalErr, m.width-common.HPad*2))
	}

	var screen string
	switch m.state {
	case stateWelcomeScreen:
		screen = m.welcome.View()
	case stateMainScreen:
		screen = m.main.View()
	}

	return common.ScreenStyle.Render(screen)
}
