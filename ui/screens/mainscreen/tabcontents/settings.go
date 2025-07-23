package tabcontents

import (
	"strings"

	"github.com/Golevka2001/bassit/ui/common"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type SettingsTabModel struct {
	commonModel *common.CommonTabModel
}

func NewSettingsTabModel(ctm *common.CommonTabModel) SettingsTabModel {
	return SettingsTabModel{
		commonModel: ctm,
	}
}

func (m SettingsTabModel) Init() tea.Cmd {
	return nil
}

func (m SettingsTabModel) Update(msg tea.Msg) (SettingsTabModel, tea.Cmd) {
	return m, nil
}

func (m SettingsTabModel) View() string {
	var b strings.Builder

	b.WriteString(common.NormalStyle.
		Width(m.commonModel.Width).
		Height(m.commonModel.Height).
		Align(lipgloss.Center).
		Render("Coming Soon..."))

	return b.String()
}

func (m *SettingsTabModel) SyncSize() {
}
