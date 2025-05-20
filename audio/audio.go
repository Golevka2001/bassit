package audio

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	// "bassit/audio/wav"
	C "bassit/constant"
	"bassit/util"

	"github.com/ebitengine/oto/v3"
	"github.com/go-music-theory/music-theory/note"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	// "github.com/hajimehoshi/go-mp3"
)

type AudioManager struct {
	otoCtx           *oto.Context
	noteNameToPlayer map[string]*oto.Player
}

func NewAudioManager(lowestNote, highestNote note.Note) (*AudioManager, error) {
	// Prepare an Oto context
	op := &oto.NewContextOptions{
		SampleRate:   C.SampleRate,
		ChannelCount: C.ChannelCount,
		Format:       oto.FormatSignedInt16LE,
	}

	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		return nil, err
	}
	<-readyChan

	// Initialize the AudioManager
	am := &AudioManager{
		otoCtx:           otoCtx,
		noteNameToPlayer: make(map[string]*oto.Player),
	}

	genAllPossibleNotes(lowestNote, highestNote)
	am.loadNoteSamples(lowestNote, highestNote)

	return am, nil
}

func (am *AudioManager) PlayBassNote(n note.Note) {
	noteName := util.GetNoteNameWithOctave(n)

	player, ok := am.noteNameToPlayer[noteName]
	if !ok {
		return
	}
	// Reset the player to the beginning
	player.Pause()
	player.Seek(0, 0)

	player.Play()

	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}
}

func genAllPossibleNotes(lowestNote, highestNote note.Note) error {
	// Shift up
	lastNote := *note.Named(C.SrcBassSampleNoteName)
	for {
		curNote := util.GetNoteStepFrom(lastNote, 1)
		curNoteName := util.GetNoteNameWithOctave(curNote)

		// srcFilePath := fmt.Sprintf("%s%s.mp3", C.NoteSampleDir, util.GetNoteNameWithOctave(lastNote))
		// dstFilePath := fmt.Sprintf("%s%s.mp3", C.NoteSampleDir, curNoteName)
		srcFilePath := fmt.Sprintf("%s%s.wav", C.NoteSampleDir, util.GetNoteNameWithOctave(lastNote))
		dstFilePath := fmt.Sprintf("%s%s.wav", C.NoteSampleDir, curNoteName)

		_, err := os.Stat(dstFilePath)
		if err == nil {
			// File already exists
			if curNoteName == util.GetNoteNameWithOctave(highestNote) {
				break
			}

			lastNote = curNote
			continue
		}

		var cmd *exec.Cmd
		switch C.OS {
		case "windows":
			cmd = exec.Command("powershell", C.RubberBandPathForWindows, "-p", "1.0", "--fine", srcFilePath, dstFilePath)
		case "darwin":
			cmd = exec.Command("osascript", "-e", fmt.Sprintf("tell application \"Terminal\" to do script \"%s -p 1.0 --fine %s %s\"", C.RubberBandPathForDarwin, srcFilePath, dstFilePath))
		}
		if cmd == nil {
			return fmt.Errorf("unsupported OS: %s", C.OS)
		}

		err = cmd.Run()
		if err != nil {
			return err
		}

		if curNoteName == util.GetNoteNameWithOctave(highestNote) {
			break
		}

		lastNote = curNote
	}

	// Shift down
	lastNote = *note.Named(C.SrcBassSampleNoteName)
	for {
		curNote := util.GetNoteStepFrom(lastNote, -1)
		curNoteName := util.GetNoteNameWithOctave(curNote)

		// srcFilePath := fmt.Sprintf("%s%s.mp3", C.NoteSampleDir, util.GetNoteNameWithOctave(lastNote))
		// dstFilePath := fmt.Sprintf("%s%s.mp3", C.NoteSampleDir, curNoteName)
		srcFilePath := fmt.Sprintf("%s%s.wav", C.NoteSampleDir, util.GetNoteNameWithOctave(lastNote))
		dstFilePath := fmt.Sprintf("%s%s.wav", C.NoteSampleDir, curNoteName)

		_, err := os.Stat(dstFilePath)
		if err == nil {
			// File already exists
			if curNoteName == util.GetNoteNameWithOctave(lowestNote) {
				break
			}

			lastNote = curNote
			continue
		}

		var cmd *exec.Cmd
		switch C.OS {
		case "windows":
			cmd = exec.Command("powershell", C.RubberBandPathForWindows, "-p", "-1.0", "--fine", srcFilePath, dstFilePath)
		case "darwin":
			cmd = exec.Command("osascript", "-e", fmt.Sprintf("tell application \"Terminal\" to do script \"%s -p -1.0 --fine %s %s\"", C.RubberBandPathForDarwin, srcFilePath, dstFilePath))
		}
		if cmd == nil {
			return fmt.Errorf("unsupported OS: %s", C.OS)
		}

		err = cmd.Run()
		if err != nil {
			return err
		}

		if curNoteName == util.GetNoteNameWithOctave(lowestNote) {
			break
		}

		lastNote = curNote
	}

	return nil
}

func (am *AudioManager) loadNoteSamples(lowestNote, highestNote note.Note) {
	// curNote := lowestNote
	// curNoteName := util.GetNoteNameWithOctave(curNote)
	// for {
	// 	nextNoteName := util.GetNoteNameWithOctave(util.GetNoteStepFrom(curNote, 1))
	// 	if nextNoteName == util.GetNoteNameWithOctave(highestNote) {
	// 		break
	// 	}

	// 	filePath := fmt.Sprintf("%s%s.mp3", C.NoteSampleDir, curNoteName)

	// 	// Read the file into memory
	// 	fileBytes, err := os.ReadFile(filePath)
	// 	if err != nil {
	// 		curNote = util.GetNoteStepFrom(curNote, 1)
	// 		curNoteName = util.GetNoteNameWithOctave(curNote)
	// 		continue
	// 	}
	// 	// Convert the pure bytes into a reader object
	// 	fileBytesReader := bytes.NewReader(fileBytes)
	// 	// Decode file
	// 	decodedMP3, err := mp3.NewDecoder(fileBytesReader)
	// 	if err != nil {
	// 		curNote = util.GetNoteStepFrom(curNote, 1)
	// 		curNoteName = util.GetNoteNameWithOctave(curNote)
	// 		continue
	// 	}

	// 	// Create a new 'player' that will handle our sound. Paused by default.
	// 	player := am.otoCtx.NewPlayer(decodedMP3)
	// 	am.noteNameToPlayer[curNoteName] = player

	// 	// Store to the map
	// 	curNote = util.GetNoteStepFrom(curNote, 1)
	// 	curNoteName = util.GetNoteNameWithOctave(curNote)
	// }

	curNote := lowestNote
	curNoteName := util.GetNoteNameWithOctave(curNote)
	for {
		nextNoteName := util.GetNoteNameWithOctave(util.GetNoteStepFrom(curNote, 1))
		if nextNoteName == util.GetNoteNameWithOctave(highestNote) {
			break
		}

		filePath := fmt.Sprintf("%s%s.wav", C.NoteSampleDir, curNoteName)

		// Read the file into memory
		fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			curNote = util.GetNoteStepFrom(curNote, 1)
			curNoteName = util.GetNoteNameWithOctave(curNote)
			continue
		}
		// Convert the pure bytes into a reader object
		fileBytesReader := bytes.NewReader(fileBytes)
		// Decode file
		// decodedWav, err := wav.NewDecoder(fileBytesReader)
		decodedWav, err := wav.DecodeWithoutResampling(fileBytesReader)
		if err != nil {
			curNote = util.GetNoteStepFrom(curNote, 1)
			curNoteName = util.GetNoteNameWithOctave(curNote)
			continue
		}

		// Create a new 'player' that will handle our sound. Paused by default.
		player := am.otoCtx.NewPlayer(decodedWav)
		am.noteNameToPlayer[curNoteName] = player

		// Store to the map
		curNote = util.GetNoteStepFrom(curNote, 1)
		curNoteName = util.GetNoteNameWithOctave(curNote)
	}
}
