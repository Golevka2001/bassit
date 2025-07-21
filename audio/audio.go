package audio

import (
	"bytes"
	"os"
	"path/filepath"
	"time"

	"github.com/Golevka2001/bassit/bass"
	"github.com/Golevka2001/bassit/config"
	"github.com/Golevka2001/bassit/utils"

	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/ebitengine/oto/v3"
)

type BassNotePlayer struct {
	players [config.PluckTypeCount]*oto.Player
}

type AudioManager struct {
	otoCtx  *oto.Context
	players map[bass.FretboardPosition]BassNotePlayer
}

func NewAudioManager() (*AudioManager, error) {
	// Prepare an Oto context
	op := &oto.NewContextOptions{
		SampleRate:   config.SampleRate,
		ChannelCount: 2,
		Format:       oto.FormatSignedInt16LE,
	}

	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		return nil, err
	}
	<-readyChan

	// Initialize the AudioManager
	am := &AudioManager{
		otoCtx:  otoCtx,
		players: make(map[bass.FretboardPosition]BassNotePlayer),
	}

	return am, nil
}

func (am *AudioManager) PlayBassNote(pos bass.FretboardPosition, t bass.PluckType) {
	player := am.getPlayer(pos, t)
	if player == nil {
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

func (am *AudioManager) StopBassNote(pos bass.FretboardPosition) {
	for t := bass.PluckTypeNormal1; t < config.PluckTypeCount; t++ {
		player := am.getPlayer(pos, t)
		if player == nil {
			return
		}
		player.Pause()
		player.Seek(0, 0)
	}
}

func (am *AudioManager) LoadSoundpackSamples(b *bass.BassModel) {
	for t := bass.PluckTypeNormal1; t < config.PluckTypeCount; t++ {
		dirPath := filepath.Join(config.SoundpackDir(), config.SoundpackName, t.String())
		for stringIdx := range config.StringCnt {
			for fretIdx := 0; fretIdx < config.DisplayedFretCount; fretIdx++ {
				pos := bass.FretboardPosition{StringIdx: stringIdx, FretIdx: fretIdx}
				n := (*b).GetNoteAt(pos)
				if n == nil {
					continue
				}
				filePath := filepath.Join(dirPath, utils.GetNoteNameWithOctave(*n)+".wav")
				// Read the file into memory
				fileBytes, err := os.ReadFile(filePath)
				if err != nil {
					continue
				}
				// Convert the pure bytes into a reader object
				fileBytesReader := bytes.NewReader(fileBytes)
				// Decode file
				decodedWav, err := wav.DecodeWithoutResampling(fileBytesReader)
				if err != nil {
					continue
				}
				// Store to the map
				player, ok := am.players[pos]
				if !ok {
					player = BassNotePlayer{}
				}
				player.players[t] = am.otoCtx.NewPlayer(decodedWav)
				am.players[pos] = player
			}
		}
	}
}

func (am *AudioManager) getPlayer(pos bass.FretboardPosition, t bass.PluckType) *oto.Player {
	players, ok := am.players[pos]
	if !ok {
		return nil
	}
	return players.players[t]
}
