package audio

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"time"

	C "bassit/constant"
	"bassit/util"

	"github.com/200sc/klangsynthese/audio"
	"github.com/200sc/klangsynthese/audio/filter"
	"github.com/ebitengine/oto/v3"
	"github.com/go-music-theory/music-theory/note"
	"github.com/hajimehoshi/go-mp3"
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
		Format:       oto.FormatSignedInt16LE, // `go-mp3`'s format is signed 16bit integers
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

	// am.genAllPossibleNotes(highestNote, lowestNote)
	return am, nil
}

func (am *AudioManager) PlayBassNote(n note.Note) {
	// noteName := util.GetNoteNameWithOctave(n)

	// player, ok := am.noteNameToPlayer[noteName]
	// if !ok {
	// 	return
	// }
	// // Reset the player to the beginning
	// player.Seek(0, 0)

	// player.Play()

	// for player.IsPlaying() {
	// 	time.Sleep(time.Millisecond)
	// }
	file, err := os.Open(C.SrcBassSampleFilePath)
	if err != nil {
		// return err
	}
	defer file.Close()

	// Decode MP3 to raw PCM
	decMp3, err := mp3.NewDecoder(file)
	if err != nil {
		// return err
	}

	decPCM := make([]byte, decMp3.Length())
	_, err = io.ReadFull(decMp3, decPCM)
	if err != nil {
		// return err
	}

	// Wrap decoded data for pitch shifting
	baseEnc := audio.Encoding{
		Data: decPCM,
		Format: audio.Format{
			SampleRate: uint32(C.SampleRate),
			Channels:   uint16(C.ChannelCount),
			Bits:       uint16(C.BitDepth),
		},
		CanLoop: audio.CanLoop{Loop: false},
	}

	// Set up pitch shifter
	shifter, err := filter.NewFFTShifter(C.ShifterFFTFrameSize, C.ShifterOverSampleFactor)
	if err != nil {
		// return err
	}

	baseNote := *note.Named(C.SrcBassSampleNoteName)
	step := util.GetStepBetween(baseNote, n)
	encoder := shifter.PitchShift(math.Pow(2, float64(step)/12))
	encoder(&baseEnc)

	player := am.otoCtx.NewPlayer(bytes.NewReader(baseEnc.Data))
	if player == nil {
		// return fmt.Errorf("failed to create player for note %s", noteName)
	}
	defer player.Close()

	player.Play()

	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}
}

func (am *AudioManager) genAllPossibleNotes(lowestNote, highestNote note.Note) error {
	file, err := os.Open(C.SrcBassSampleFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode MP3 to raw PCM
	decMp3, err := mp3.NewDecoder(file)
	if err != nil {
		return err
	}

	decPCM := make([]byte, decMp3.Length())
	_, err = io.ReadFull(decMp3, decPCM)
	if err != nil {
		return err
	}

	// Wrap decoded data for pitch shifting
	baseEnc := audio.Encoding{
		Data: decPCM,
		Format: audio.Format{
			SampleRate: uint32(C.SampleRate),
			Channels:   uint16(C.ChannelCount),
			Bits:       uint16(C.BitDepth),
		},
		CanLoop: audio.CanLoop{Loop: false},
	}

	// Set up pitch shifter
	shifter, err := filter.NewFFTShifter(C.ShifterFFTFrameSize, C.ShifterOverSampleFactor)
	if err != nil {
		return err
	}

	// Shift up
	for {
		baseNote := *note.Named(C.SrcBassSampleNoteName)
		step := util.GetStepBetween(baseNote, highestNote)
		for i := 0; i <= min(step, 12); i++ {
			curNote := util.GetNoteStepFrom(baseNote, i)
			curNoteName := util.GetNoteNameWithOctave(curNote)

			encoder := shifter.PitchShift(C.ShiftUpFactors[i])
			encoder(&baseEnc)

			// Create a new player
			player := am.otoCtx.NewPlayer(bytes.NewReader(baseEnc.Data))
			if player == nil {
				return fmt.Errorf("failed to create player for note %s", curNoteName)
			}

			// Store the player in the map
			am.noteNameToPlayer[curNoteName] = player

			if i == 12 {
				baseNote = curNote
			} else {
				baseEnc.Data = decPCM // Reset to the original PCM data for the next iteration
			}
		}
		if step < 12 {
			break
		}
	}

	// Shift down
	for {
		baseNote := *note.Named(C.SrcBassSampleNoteName)
		step := util.GetStepBetween(baseNote, lowestNote)
		for i := 0; i <= min(step, 12); i++ {
			curNote := util.GetNoteStepFrom(baseNote, -i)
			curNoteName := util.GetNoteNameWithOctave(curNote)

			encoder := shifter.PitchShift(C.ShiftDownFactors[i])
			encoder(&baseEnc)

			// Create a new player
			player := am.otoCtx.NewPlayer(bytes.NewReader(baseEnc.Data))
			if player == nil {
				return fmt.Errorf("failed to create player for note %s", curNoteName)
			}

			// Store the player in the map
			am.noteNameToPlayer[curNoteName] = player

			if i == 12 {
				baseNote = curNote
			} else {
				baseEnc.Data = decPCM // Reset to the original PCM data for the next iteration
			}
		}
		if step < 12 {
			break
		}
	}

	return nil
}

func loadMp3File(filePath string) (*mp3.Decoder, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	fileBytesReader := bytes.NewReader(fileBytes)

	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		return nil, err
	}

	return decodedMp3, nil
}
