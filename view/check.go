package view

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	C "bassit/constant"

	tcheck "github.com/Golevka2001/go-tcheck"
	"github.com/gdamore/tcell/v2"
)

type CheckView struct {
	screen       *tcell.Screen
	checkManager *tcheck.CheckManager
	checkUI      *tcheck.UIRenderer
}

func NewCheckView(s *tcell.Screen) *CheckView {
	// Initialize `tcheck`
	var ui *tcheck.UIRenderer
	cm := tcheck.NewCheckManager(func() {
		if ui != nil {
			ui.Draw()
		}
	}, 3)
	ui = tcheck.NewUIRenderer(*s, cm)

	// Check if `rubberband` is available
	cm.AddCheck("Checking `rubberband`", checkRubberband)
	// Check if the audio files are available
	cm.AddCheck("Checking audio resources", checkAudioResources)

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

	// Collect results
	failed := []string{}
	for _, item := range cv.checkManager.GetItems() {
		if item.Status == tcheck.StatusFailed {
			failed = append(failed, fmt.Sprintf("%s: %v", item.Name, item.Error))
		}
	}

	return failed
}

// Fini stops the event loop and cleans up the UI
func (cv *CheckView) Fini() {
	if cv.checkUI != nil {
		cv.checkUI.Stop()
	}
	// Sleep for a bit...zZZ
	time.Sleep(500 * time.Millisecond)
	(*cv.screen).Clear()
}

func checkRubberband(reporter tcheck.SubProgressReporter) error {
	// Check if the binary exists (only for Windows and macOS)
	if C.OS == "windows" || C.OS == "darwin" {
		isExist := false
		reporter.ReportSubProgress(0, "Checking if binary exists")
		if _, err := os.Stat(C.RubberbandCommand); err == nil {
			isExist = true
		}
		if !isExist {
			msg := "Rubberband binary not found or not executable"
			reporter.ReportSubProgress(100, msg)
			return fmt.Errorf("%s", msg)
		}
	}

	// Check if the command is available
	reporter.ReportSubProgress(50, "Checking if command is available")
	var cmd *exec.Cmd
	switch C.OS {
	case "windows":
		cmd = exec.Command("powershell", C.RubberbandCommand, "-V")
	case "darwin":
		cmd = exec.Command(C.RubberbandCommand, "-V")
	default:
		// Check if `rubberband-r3` is available
		tmp := exec.Command(C.RubberbandCommand, "-V") // A Cmd cannot be reused after calling its Run, Output or CombinedOutput methods
		if err := tmp.Run(); err != nil {
			// If not, check if `rubberband` is available
			C.RubberbandCommand = "rubberband"
		}
		cmd = exec.Command(C.RubberbandCommand, "-V")
	}
	if err := cmd.Run(); err != nil {
		msg := "Command not available"
		reporter.ReportSubProgress(100, msg)
		return fmt.Errorf("%s", msg)
	}

	reporter.ReportSubProgress(100, "Passed")
	return nil
}

func checkAudioResources(reporter tcheck.SubProgressReporter) error {
	reporter.ReportSubProgress(0, "Checking if audio resources are readable")
	filePath := filepath.Join(C.NoteSampleDir, C.SrcBassSampleNoteName+".wav")
	if file, err := os.Open(filePath); err == nil {
		defer file.Close()
	} else {
		msg := "Audio resources not found or not readable"
		reporter.ReportSubProgress(100, msg)
		return fmt.Errorf("%s", msg)
	}

	reporter.ReportSubProgress(100, "Passed")
	return nil
}
