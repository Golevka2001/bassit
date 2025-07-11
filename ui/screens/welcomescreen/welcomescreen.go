package welcomescreen

import (
	"strings"
	"time"

	"github.com/Golevka2001/bassit/audio"
	"github.com/Golevka2001/bassit/bass"
	"github.com/Golevka2001/bassit/ui/common"
	"github.com/Golevka2001/bassit/ui/components/asciiart"

	"github.com/charmbracelet/bubbles/v2/spinner"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

const (
	timeout = 2 * time.Second

	// staticElementsHeight represents the height of the title and loading hint
	// TODO: make it dynamic
	staticElementsHeight = 2
)

type timeoutMsg struct{}

type resourcesLoadedMsg struct{}

type Model struct {
	commonModel *common.CommonScreenModel

	spinner  spinner.Model
	asciiArt asciiart.Model

	done        bool
	timeoutFlag bool

	bass  *bass.BassModel
	audio *audio.AudioManager
}

func NewModel(csm *common.CommonScreenModel) Model {
	m := Model{
		commonModel: csm,
		spinner: spinner.Model{
			Spinner: spinner.Pulse,
		},
		asciiArt: asciiart.New(bassHeadstockAsciiArt),
		bass:     csm.Context.Bass,
		audio:    csm.Context.Audio,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tea.ClearScreen,
		m.spinner.Tick,
		loadSoundpackCmd(m.bass, m.audio),
		func() tea.Msg {
			time.Sleep(timeout)
			return timeoutMsg{}
		},
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case resourcesLoadedMsg:
		m.done = true
		if m.timeoutFlag {
			cmds = append(cmds, func() tea.Msg {
				return common.SwitchScreenMsg{}
			})
		}

	case timeoutMsg:
		m.timeoutFlag = true
		if m.done {
			cmds = append(cmds, func() tea.Msg {
				return common.SwitchScreenMsg{}
			})
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var b strings.Builder

	// Render the title
	rTitle := common.TitleStyle.
		Width(m.commonModel.Width).
		Render(common.TitleText)

	// Render the ASCII art
	rAsciiArt := m.asciiArt.View()

	// Render the loading hint
	loadingHint := m.spinner.View() + " " + "Loading..."
	rLoadingHint := common.NormalStyle.Render(loadingHint)

	b.WriteString(lipgloss.JoinVertical(lipgloss.Left,
		rTitle,
		rAsciiArt,
		rLoadingHint,
	))

	return b.String()
}

func (m *Model) SyncSize() {
	m.asciiArt.SetSize(
		m.commonModel.Width,
		m.commonModel.Height-staticElementsHeight,
	)
}

const bassHeadstockAsciiArt = `
            ╭◜‾‾◝╮        ╭◜‾‾◝╮        ╭◜‾‾◝╮        ╭◜‾‾◝╮           
          ◞◜      ◝◟    ◞◜      ◝◟    ◞◜      ◝◟    ◞◜      ◝◟         
         (          )  (          )  (          )  (          )        
          ◝‒◟ ╭╮ ◞‒◜    ◝‒◟ ╭╮ ◞‒◜    ◝‒◟ ╭╮ ◞‒◜    ◝‒◟ ╭╮ ◞‒◜         
             ╰┤├╯          ╰┤├╯          ╰┤├╯          ╰┤├╯            
       ◜‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾--_     
      /                                                           ‾-_  
     /   ◜‾‾╭◝         ◜‾|‾◝         ◜╮‾‾◝         ◜‾|‾◝             \ 
    /    ▏ / ▕         ▏ | ▕         ▏ \ ▕         ▏ | ▕              \
   /     ◟╯__◞         ◟_|_◞         ◟__╰◞         ◟_|_◞              ▕
__◞__--‾‾‾     ___---‾‾‾     ___---‾‾‾     ___---‾‾‾                  ▕
\\    ___---‾‾‾     /‾‾‾‾\‾‾‾     ___---‾‾‾                           /
 \\‾‾‾     ___---‾‾‾▏(◜◞) ▏_---‾‾‾      ◟_                           / 
  \\_---‾‾‾     ___-\____/            _ ‾\◞  __--‾‾-_               /  
   \\  ___---‾‾‾                  (‾   \  _-‾        ‾◟_         _◞‾   
    \\‾                       (‾   ‾‾) _-‾              ‾------‾‾      
     \\__---◝_    ╭‾‾\   ◞-╮   ‾‾) ‾‾-‾                                
              ◝_   \_◞-_  (‾\  ‾‾_-‾                                   
                ◝_  \   \ ╰-╯‾_-‾                                      
                  ◝_ \_◞╯  _-‾                                         
                    \__--‾‾                                            
`
