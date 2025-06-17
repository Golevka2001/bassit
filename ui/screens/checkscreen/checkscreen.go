package checkscreen

import (
	"fmt"
	"strings"
	"time"

	"github.com/Golevka2001/bassit/ui/common"

	"github.com/charmbracelet/bubbles/v2/progress"
	"github.com/charmbracelet/bubbles/v2/spinner"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

const (
	maxConcurrentRuns = 2

	title            = "🎸 Pre-Gig Checklist"
	normalHelpMsg    = "Press [Esc] to skip checks or [Ctrl+C] to quit"
	passedHelpMsg    = "All checks passed! Starting the program..."
	hasInfoHelpMsg   = "Press any key to continue or [Ctrl+C] to quit"
	hasFailedHelpMsg = "Some checks failed, press any key to quit"
)

var (
	checklistStyle = common.NormalStyle.
		MarginTop(1).
		MarginBottom(2)
)

// checkResultMsg is sent when a check completes
type checkResultMsg struct {
	index  int
	status checkStatus
}

type Model struct {
	commonModel *common.CommonScreenModel

	spinner  spinner.Model
	progress progress.Model

	checkList        []*checkItem
	runningChecksCnt int
	maxConcurrency   int
	hasFailed        bool
	hasInfo          bool
	done             bool
}

func NewModel(csm *common.CommonScreenModel) Model {
	// Register the checks to be performed
	var list []*checkItem
	registerCheck(&list, "Some time-consuming check", func(update stepUpdater) checkStatus {
		update("This is a placeholder for a time-consuming check")
		// Simulate a time-consuming check
		for i := 0; i < 5; i++ {
			update(fmt.Sprintf("Step %d of 5", i+1))
			// Simulate some work
			time.Sleep(1 * time.Second) // Uncomment to simulate delay
		}
		return statusPassed
	})
	registerCheck(&list, "Check Configuration", checkConfiguration)
	registerCheck(&list, "Check Rubberband", checkRubberband)
	registerCheck(&list, "Check Audio Resources", checkAudioResources)

	m := Model{
		commonModel:    csm,
		spinner:        spinner.New(),
		progress:       progress.New(progress.WithDefaultGradient()),
		maxConcurrency: maxConcurrentRuns,
		checkList:      list,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.ExitAltScreen,
		m.spinner.Tick,
		m.progress.SetPercent(0),
		m.runChecksConcurrently(),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		if !m.done {
			// If not done, press [Esc] to skip checks
			if key == "esc" {
				cmds = append(cmds, func() tea.Msg {
					return common.SwitchScreenMsg{}
				})
			}
		} else {
			// When all checks are done
			// If any check has failed, press any key to quit
			// If any check has info, press any key to continue
			if m.hasFailed {
				return m, tea.Quit
			} else if m.hasInfo {
				cmds = append(cmds, func() tea.Msg {
					return common.SwitchScreenMsg{}
				})
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case progress.FrameMsg:
		var cmd tea.Cmd
		m.progress, cmd = m.progress.Update(msg)
		cmds = append(cmds, cmd)

	case checkResultMsg:
		item := m.checkList[msg.index]
		item.status = msg.status
		m.runningChecksCnt--

		// Update progress bar
		completed := 0
		for _, it := range m.checkList {
			if it.status != statusPending && it.status != statusInProgress {
				completed++
			}
		}
		cmd := m.progress.SetPercent(float64(completed) / float64(len(m.checkList)))
		cmds = append(cmds, cmd)

		// Check if all checks are done
		allDone := true
		hasFailed := false
		hasInfo := false
		for _, it := range m.checkList {
			if it.status == statusPending || it.status == statusInProgress {
				allDone = false
				break
			}
			if it.status == statusFailed {
				hasFailed = true
			}
			if it.status == statusInfo {
				hasInfo = true
			}
		}
		if allDone {
			m.hasFailed = hasFailed
			m.hasInfo = hasInfo
			m.done = allDone
			// If all checks are passed, automatically switch to the next view
			if !hasFailed && !hasInfo {
				time.Sleep(500 * time.Millisecond)
				cmds = append(cmds, func() tea.Msg {
					return common.SwitchScreenMsg{}
				})
			}
		} else {
			// If not all checks are done, run more checks if possible
			if m.runningChecksCnt < m.maxConcurrency {
				cmds = append(cmds, m.runChecksConcurrently())
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var b strings.Builder

	// Render the title
	rTitle := common.TitleStyle.
		Width(m.commonModel.Width).
		Render(title)

	// Render the checklist
	rChecklist := checklistStyle.Render(m.renderChecklist())

	// Render the help message
	helpMsg := normalHelpMsg
	if m.done {
		// Check if any check has failed or has info
		if m.hasFailed {
			helpMsg = hasFailedHelpMsg
		} else if m.hasInfo {
			helpMsg = hasInfoHelpMsg
		} else {
			helpMsg = passedHelpMsg
		}
	}
	rHelp := common.HelpStyle.
		MarginBottom(1).
		Render(helpMsg)

	// Render the progress bar
	rProgress := m.progress.View() + "\n"

	b.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		rTitle,
		rChecklist,
		rHelp,
		rProgress,
	))

	return b.String()
}

func (m *Model) SyncSize() {
	// TODO: this might be a bug of `bubbles/progress`, check it later
	m.progress.SetWidth(m.commonModel.Width)
}
