package audio

import (
	"bytes"
	"os"
	"path/filepath"
	"time"

	C "github.com/Golevka2001/bassit/constant"
	"github.com/Golevka2001/bassit/util"

	"github.com/ebitengine/oto/v3"
	"github.com/go-music-theory/music-theory/note"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type AudioManager struct {
	otoCtx           *oto.Context
	noteNameToPlayer map[string]*oto.Player
}

func NewAudioManager() (*AudioManager, error) {
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

func (am *AudioManager) StopBassNote(n note.Note) {
	// FIXME: When the same note is played at different positions (e.g. (0,2) and (1,7) are both A2), this function stops both of them.
	// One possible solution is to assign a unique player to each position, but this would consume more memory and introduce noise.
	noteName := util.GetNoteNameWithOctave(n)

	player, ok := am.noteNameToPlayer[noteName]
	if !ok {
		return
	}
	player.Pause()
	player.Seek(0, 0)
}

func (am *AudioManager) LoadNoteSamples(lowestNote, highestNote note.Note) {
	curNote := lowestNote
	curNoteName := util.GetNoteNameWithOctave(curNote)
	for {
		nextNoteName := util.GetNoteNameWithOctave(util.GetNoteStepFrom(curNote, 1))
		if nextNoteName == util.GetNoteNameWithOctave(highestNote) {
			break
		}

		filePath := filepath.Join(C.NoteSampleDir, curNoteName+".wav")

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
