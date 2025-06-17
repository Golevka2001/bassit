package checkscreen

import (
	tea "github.com/charmbracelet/bubbletea/v2"
)

type checkStatus int

const (
	statusPending checkStatus = iota
	statusInProgress
	statusPassed
	statusInfo
	statusFailed
)

// stepUpdater is used to update the current step of a check
type stepUpdater func(string)

// checkFunc will be called to perform a check
type checkFunc func(update stepUpdater) checkStatus

// checkItem represents a single check to be performed
type checkItem struct {
	title     string
	status    checkStatus
	curStep   string
	checkFunc checkFunc
}

// runChecksConcurrently runs checks concurrently based on the maximum concurrency limit
// It should be called after the initial setup and whenever a check completes
func (m *Model) runChecksConcurrently() tea.Cmd {
	var cmds []tea.Cmd

	for i, item := range m.checkList {
		if item.status == statusPending && m.runningChecksCnt < m.maxConcurrency {
			item.status = statusInProgress
			m.runningChecksCnt++
			cmds = append(cmds, runCheck(i, item))
		}
	}

	return tea.Batch(cmds...)
}

// registerCheck is a helper function to add an item to the checklist
func registerCheck(checkList *[]*checkItem, title string, fn checkFunc) {
	*checkList = append(*checkList, &checkItem{
		title:     title,
		status:    statusPending,
		curStep:   "",
		checkFunc: fn,
	})
}

// runCheck runs a check in a separate goroutine and returns a command to update the model
func runCheck(index int, item *checkItem) tea.Cmd {
	return func() tea.Msg {
		stepChan := make(chan string)
		done := make(chan struct{})

		go func() {
			for {
				select {
				case step := <-stepChan:
					item.curStep = step
				case <-done:
					return
				}
			}
		}()

		status := item.checkFunc(func(step string) {
			stepChan <- step
		})

		close(done)
		return checkResultMsg{
			index:  index,
			status: status,
		}
	}
}
