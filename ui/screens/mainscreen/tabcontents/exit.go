package tabcontents

import (
	"strings"

	"github.com/Golevka2001/bassit/ui/common"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type ExitTabModel struct {
	commonModel *common.CommonTabModel
}

func NewExitTabModel(ctm *common.CommonTabModel) ExitTabModel {
	return ExitTabModel{
		commonModel: ctm,
	}
}

func (m ExitTabModel) Init() tea.Cmd {
	return nil
}

func (m ExitTabModel) Update(msg tea.Msg) (ExitTabModel, tea.Cmd) {
	return m, nil
}

func (m ExitTabModel) View() string {
	var b strings.Builder

	b.WriteString(common.NormalStyle.
		Width(m.commonModel.Width).
		Height(m.commonModel.Height).
		Align(lipgloss.Center).
		Render("Comming Soon..."))

	return b.String()
}

func (m *ExitTabModel) SyncSize() {
}
