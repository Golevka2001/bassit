package tabselector

type FocusMsg struct {
	TabIdx int
}

type UnfocusMsg struct{}

type SwitchToTabMsg struct {
	TabIdx int
}
