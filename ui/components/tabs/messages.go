package tabs

type SwitchToNextTabMsg struct{}

type SwitchToPreviousTabMsg struct{}

type SwitchToTabMsg struct {
	TabIdx int
}
