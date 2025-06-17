package common

// ErrMsg is a message that contains an error
type ErrMsg struct{ err error }

func (e ErrMsg) Error() string { return e.err.Error() }

// SwitchScreenMsg informs the ui manager to switch to next screen
type SwitchScreenMsg struct{}
