package view

import (
	C "bassit/constant"
	"fmt"
	"os"
	"os/exec"

	tcheck "github.com/Golevka2001/go-tcheck"
	"github.com/gdamore/tcell/v2"
)

type CheckView struct {
	screen       *tcell.Screen
	checkManager *tcheck.CheckManager
	checkUI      *tcheck.UIRenderer
}

func NewCheckView(s *tcell.Screen) *CheckView {
	var ui *tcheck.UIRenderer
	cm := tcheck.NewCheckManager(func() {
		if ui != nil {
			ui.Draw()
		}
	}, 3)
	ui = tcheck.NewUIRenderer(*s, cm)

	addChecks(cm)

	return &CheckView{
		screen:       s,
		checkManager: cm,
		checkUI:      ui,
	}
}

func (cv *CheckView) RunChecks() []string {
	// Start running checks and event loop
	go cv.checkManager.RunAllChecks()
	cv.checkUI.Run()

	failed := []string{}
	for _, item := range cv.checkManager.GetItems() {
		if item.Status == tcheck.StatusFailed {
			failed = append(failed, fmt.Sprintf("%s: %v", item.Name, item.Error))
		}
	}

	return failed
}

func addChecks(cm *tcheck.CheckManager) {
	// Check if `rubberband` is available
	cm.AddCheck(
		"Checking if `rubberband` is available",
		checkRubberband,
	)
}

func checkRubberband(reporter tcheck.SubProgressReporter) error {
	var cmd *exec.Cmd

	// Check if the binary exists
	isExist := false
	reporter.ReportSubProgress(0, "Checking existence")
	if C.OS == "windows" || C.OS == "darwin" {
		if _, err := os.Stat(C.RubberBandPathForWindows); err == nil {
			isExist = true
		}
	}
	if !isExist {
		msg := "Rubberband binary not found or not executable"
		reporter.ReportSubProgress(100, msg)
		return fmt.Errorf("%s", msg)
	}

	// Check if the command is available
	reporter.ReportSubProgress(50, "Checking availability")
	switch C.OS {
	case "windows":
		cmd = exec.Command("powershell", C.RubberBandPathForWindows, "-V")
	case "darwin":
		cmd = exec.Command(C.RubberBandPathForDarwin, "-V")
	default:
		// Check if `rubberband-r3` is available
		cmd = exec.Command(C.RubberBandCommand, "-V")
		if err := cmd.Run(); err != nil {
			// If not, check if `rubberband` is available
			C.RubberBandCommand = "rubberband"
			cmd = exec.Command(C.RubberBandCommand, "-V")
		}
	}
	if err := cmd.Run(); err != nil {
		msg := "Command not available"
		reporter.ReportSubProgress(100, msg)
		return fmt.Errorf("%s", msg)
	}

	reporter.ReportSubProgress(100, "Available")
	return nil
}
