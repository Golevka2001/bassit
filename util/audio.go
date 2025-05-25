package util

import (
	"fmt"
	"os"
	"os/exec"

	C "bassit/constant"

	"github.com/go-music-theory/music-theory/note"
)

func GenAllPossibleNotes(lowestNote, highestNote note.Note) error {
	// Shift up
	lastNote := *note.Named(C.SrcBassSampleNoteName)
	for {
		curNote := GetNoteStepFrom(lastNote, 1)
		curNoteName := GetNoteNameWithOctave(curNote)

		srcFilePath := fmt.Sprintf("%s%s.wav", C.NoteSampleDir, GetNoteNameWithOctave(lastNote))
		dstFilePath := fmt.Sprintf("%s%s.wav", C.NoteSampleDir, curNoteName)

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
		switch C.OS {
		case "windows":
			cmd = exec.Command("powershell", C.RubberBandPathForWindows, "-p", "1.0", srcFilePath, dstFilePath)
		case "darwin":
			cmd = exec.Command(C.RubberBandPathForDarwin, "-p", "1.0", srcFilePath, dstFilePath)
		default:
			cmd = exec.Command(C.RubberBandCommand, "-p", "1.0", srcFilePath, dstFilePath)
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

		srcFilePath := fmt.Sprintf("%s%s.wav", C.NoteSampleDir, GetNoteNameWithOctave(lastNote))
		dstFilePath := fmt.Sprintf("%s%s.wav", C.NoteSampleDir, curNoteName)

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
		switch C.OS {
		case "windows":
			cmd = exec.Command("powershell", C.RubberBandPathForWindows, "-p", "-1.0", srcFilePath, dstFilePath)
		case "darwin":
			cmd = exec.Command(C.RubberBandPathForDarwin, "-p", "-1.0", srcFilePath, dstFilePath)
		default:
			cmd = exec.Command(C.RubberBandCommand, "-p", "-1.0", srcFilePath, dstFilePath)
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
