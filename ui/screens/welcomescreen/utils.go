package welcomescreen

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/Golevka2001/bassit/audio"
	"github.com/Golevka2001/bassit/bass"
	"github.com/Golevka2001/bassit/utils"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/go-music-theory/music-theory/note"
)

func loadResourcesCmd(b *bass.BassModel, a *audio.AudioManager) tea.Cmd {
	return func() tea.Msg {
		lowestNote, highestNote := b.GetLowestAndHighestNotes()

		if err := GenAllPossibleNotes(lowestNote, highestNote); err != nil {
			return err
		}

		a.LoadNoteSamples(lowestNote, highestNote)
		return resourcesLoadedMsg{}
	}
}

func GenAllPossibleNotes(lowestNote, highestNote note.Note) error {
	// Shift up
	lastNote := *note.Named(audio.SrcBassSampleNoteName)
	for {
		curNote := utils.GetNoteStepFrom(lastNote, 1)
		curNoteName := utils.GetNoteNameWithOctave(curNote)

		srcFilePath := filepath.Join(audio.NoteSampleDir, utils.GetNoteNameWithOctave(lastNote)+".wav")
		dstFilePath := filepath.Join(audio.NoteSampleDir, curNoteName+".wav")

		_, err := os.Stat(dstFilePath)
		if err == nil {
			// File already exists
			if curNoteName == utils.GetNoteNameWithOctave(highestNote) {
				break
			}

			lastNote = curNote
			continue
		}

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("powershell", audio.RubberbandCommand, "-p", "1.0", srcFilePath, dstFilePath)
		} else {
			cmd = exec.Command(audio.RubberbandCommand, "-p", "1.0", srcFilePath, dstFilePath)
		}

		err = cmd.Run()
		if err != nil {
			return err
		}

		if curNoteName == utils.GetNoteNameWithOctave(highestNote) {
			break
		}

		lastNote = curNote
	}

	// Shift down
	lastNote = *note.Named(audio.SrcBassSampleNoteName)
	for {
		curNote := utils.GetNoteStepFrom(lastNote, -1)
		curNoteName := utils.GetNoteNameWithOctave(curNote)

		srcFilePath := filepath.Join(audio.NoteSampleDir, utils.GetNoteNameWithOctave(lastNote)+".wav")
		dstFilePath := filepath.Join(audio.NoteSampleDir, curNoteName+".wav")

		_, err := os.Stat(dstFilePath)
		if err == nil {
			// File already exists
			if curNoteName == utils.GetNoteNameWithOctave(lowestNote) {
				break
			}

			lastNote = curNote
			continue
		}

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("powershell", audio.RubberbandCommand, "-p", "-1.0", srcFilePath, dstFilePath)
		} else {
			cmd = exec.Command(audio.RubberbandCommand, "-p", "-1.0", srcFilePath, dstFilePath)
		}

		err = cmd.Run()
		if err != nil {
			return err
		}

		if curNoteName == utils.GetNoteNameWithOctave(lowestNote) {
			break
		}

		lastNote = curNote
	}

	return nil
}
