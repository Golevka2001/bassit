package util

import (
	"os"
	"os/exec"
	"path/filepath"

	C "bassit/constant"

	"github.com/go-music-theory/music-theory/note"
)

func GenAllPossibleNotes(lowestNote, highestNote note.Note) error {
	// Shift up
	lastNote := *note.Named(C.SrcBassSampleNoteName)
	for {
		curNote := GetNoteStepFrom(lastNote, 1)
		curNoteName := GetNoteNameWithOctave(curNote)

		srcFilePath := filepath.Join(C.NoteSampleDir, GetNoteNameWithOctave(lastNote)+".wav")
		dstFilePath := filepath.Join(C.NoteSampleDir, curNoteName+".wav")

		_, err := os.Stat(dstFilePath)
		if err == nil {
			// File already exists
			if curNoteName == GetNoteNameWithOctave(highestNote) {
				break
			}

			lastNote = curNote
			continue
		}

		var cmd *exec.Cmd
		if C.OS == "windows" {
			cmd = exec.Command("powershell", C.RubberbandCommand, "-p", "1.0", srcFilePath, dstFilePath)
		} else {
			cmd = exec.Command(C.RubberbandCommand, "-p", "1.0", srcFilePath, dstFilePath)
		}

		err = cmd.Run()
		if err != nil {
			return err
		}

		if curNoteName == GetNoteNameWithOctave(highestNote) {
			break
		}

		lastNote = curNote
	}

	// Shift down
	lastNote = *note.Named(C.SrcBassSampleNoteName)
	for {
		curNote := GetNoteStepFrom(lastNote, -1)
		curNoteName := GetNoteNameWithOctave(curNote)

		srcFilePath := filepath.Join(C.NoteSampleDir, GetNoteNameWithOctave(lastNote)+".wav")
		dstFilePath := filepath.Join(C.NoteSampleDir, curNoteName+".wav")

		_, err := os.Stat(dstFilePath)
		if err == nil {
			// File already exists
			if curNoteName == GetNoteNameWithOctave(lowestNote) {
				break
			}

			lastNote = curNote
			continue
		}

		var cmd *exec.Cmd
		if C.OS == "windows" {
			cmd = exec.Command("powershell", C.RubberbandCommand, "-p", "-1.0", srcFilePath, dstFilePath)
		} else {
			cmd = exec.Command(C.RubberbandCommand, "-p", "-1.0", srcFilePath, dstFilePath)
		}

		err = cmd.Run()
		if err != nil {
			return err
		}

		if curNoteName == GetNoteNameWithOctave(lowestNote) {
			break
		}

		lastNote = curNote
	}

	return nil
}
