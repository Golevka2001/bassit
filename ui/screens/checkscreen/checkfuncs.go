package checkscreen

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/Golevka2001/bassit/audio"
)

func checkConfiguration(update stepUpdater) checkStatus {
	// TODO
	update("I will implement this later")
	return statusInfo
}

func checkRubberband(update stepUpdater) checkStatus {
	update("Checking if binary exists")
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		if _, err := os.Stat(audio.RubberbandCommand); err != nil {
			update("Rubberband binary not found or not executable")
			return statusFailed
		}
	}

	update("Checking if command is available")
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell", audio.RubberbandCommand, "-V")
	case "darwin":
		cmd = exec.Command(audio.RubberbandCommand, "-V")
	default:
		// Try rubberband-r3 first
		tmp := exec.Command(audio.RubberbandCommand, "-V")
		if err := tmp.Run(); err != nil {
			audio.RubberbandCommand = "rubberband"
		}
		cmd = exec.Command(audio.RubberbandCommand, "-V")
	}

	if err := cmd.Run(); err != nil {
		update("Command not available")
		return statusFailed
	}

	update("Rubberband check passed")
	return statusPassed
}

func checkAudioResources(update stepUpdater) checkStatus {
	update("Checking if audio resources are readable")
	filePath := filepath.Join(audio.NoteSampleDir, audio.SrcBassSampleNoteName+".wav")
	file, err := os.Open(filePath)
	if err != nil {
		update("Audio resources not found or not readable")
		return statusFailed
	}
	file.Close()

	update("Audio resources check passed")
	return statusPassed
}
