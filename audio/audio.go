package audio

import (
	"bytes"
	"os"
	"path/filepath"
	"time"

	"github.com/Golevka2001/bassit/utils"

	"github.com/ebitengine/oto/v3"
	"github.com/go-music-theory/music-theory/note"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type AudioManager struct {
	otoCtx           *oto.Context
	noteNameToPlayer map[string]*oto.Player
}

func New() (*AudioManager, error) {
	// Prepare an Oto context
	op := &oto.NewContextOptions{
		SampleRate:   SampleRate,
		ChannelCount: ChannelCount,
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

	return am, nil
}

func (am *AudioManager) PlayBassNote(n note.Note) {
	noteName := utils.GetNoteNameWithOctave(n)

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

func (am *AudioManager) StopBassNote(n note.Note) {
	// FIXME: When the same note is played at different positions (e.g. (0,2) and (1,7) are both A2), this function stops both of them.
	// One possible solution is to assign a unique player to each position, but this would consume more memory and introduce noise.
	noteName := utils.GetNoteNameWithOctave(n)

	player, ok := am.noteNameToPlayer[noteName]
	if !ok {
		return
	}
	player.Pause()
	player.Seek(0, 0)
}

func (am *AudioManager) LoadNoteSamples(lowestNote, highestNote note.Note) {
	curNote := lowestNote
	curNoteName := utils.GetNoteNameWithOctave(curNote)
	for {
		filePath := filepath.Join(NoteSampleDir, curNoteName+".wav")

		// Read the file into memory
		fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			curNote = utils.GetNoteStepFrom(curNote, 1)
			curNoteName = utils.GetNoteNameWithOctave(curNote)
			continue
		}
		// Convert the pure bytes into a reader object
		fileBytesReader := bytes.NewReader(fileBytes)
		// Decode file
		decodedWav, err := wav.DecodeWithoutResampling(fileBytesReader)
		if err != nil {
			curNote = utils.GetNoteStepFrom(curNote, 1)
			curNoteName = utils.GetNoteNameWithOctave(curNote)
			continue
		}

		// Create a new player that will handle our sound
		player := am.otoCtx.NewPlayer(decodedWav)

		// Store to the map
		am.noteNameToPlayer[curNoteName] = player

		if curNoteName == utils.GetNoteNameWithOctave(highestNote) {
			break
		}
		// Move to the next note
		curNote = utils.GetNoteStepFrom(curNote, 1)
		curNoteName = utils.GetNoteNameWithOctave(curNote)
	}
}
