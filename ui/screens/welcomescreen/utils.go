package welcomescreen

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/Golevka2001/bassit/audio"
	"github.com/Golevka2001/bassit/bass"
	"github.com/Golevka2001/bassit/config"
	"github.com/Golevka2001/bassit/utils"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/go-music-theory/music-theory/note"
)

func loadSoundpackCmd(b *bass.BassModel, a *audio.AudioManager) tea.Cmd {
	return func() tea.Msg {
		lowestNote, highestNote := b.GetRange()

		if err := genPitchShiftedSamples(lowestNote, highestNote); err != nil {
			return err // FIXME: return error to the screen
		}

		a.LoadSoundpackSamples(b)
		return resourcesLoadedMsg{}
	}
}

func genPitchShiftedSamples(lowestNote, highestNote note.Note) error {
	curSoundpackDir := filepath.Join(config.SoundpackDir(), config.SoundpackName)
	subDirNames := []string{"normal1", "normal2", "slap1", "slap2", "mute1", "mute2"}
	baseNote := *note.Named(config.BasePitch)

	var wg sync.WaitGroup
	errChan := make(chan error, len(subDirNames))

	for i, subDirName := range subDirNames {
		wg.Add(1)
		go func(i int, subDirName string) {
			defer wg.Done()
			if err := processSubDir(i, subDirName, curSoundpackDir, baseNote, lowestNote, highestNote); err != nil {
				errChan <- err
			}
		}(i, subDirName)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

// processSubDir generates audio samples for a single subdirectory
func processSubDir(i int, subDirName, curSoundpackDir string, baseNote, lowestNote, highestNote note.Note) error {
	// Check if directory exists, create if not
	if _, err := os.Stat(filepath.Join(curSoundpackDir, subDirName)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(curSoundpackDir, subDirName), 0755); err != nil {
			return err
		}
	}

	// Copy base audio files
	var copySrcFilePath, copyDstFilePath string
	switch subDirName {
	case "normal1", "normal2":
		copySrcFilePath = filepath.Join(curSoundpackDir, config.Normal[i].File)
		copyDstFilePath = filepath.Join(curSoundpackDir, subDirName, config.BasePitch+".wav")
	case "slap1", "slap2":
		copySrcFilePath = filepath.Join(curSoundpackDir, config.Slap[i-2].File)
		copyDstFilePath = filepath.Join(curSoundpackDir, subDirName, config.BasePitch+".wav")
	case "mute1", "mute2":
		copySrcFilePath = filepath.Join(curSoundpackDir, config.Mute[i-4].File)
		copyDstFilePath = filepath.Join(curSoundpackDir, subDirName, config.BasePitch+".wav")
	}

	if _, err := os.Stat(copyDstFilePath); os.IsNotExist(err) {
		if err := copyFile(copySrcFilePath, copyDstFilePath); err != nil {
			return err
		}
	}

	// Generate pitch shifted samples in both directions
	for _, direction := range []struct {
		step    int
		pitch   string
		endNote note.Note
	}{
		{1, "1.0", highestNote},  // Shift up
		{-1, "-1.0", lowestNote}, // Shift down
	} {
		lastNote := baseNote
		for {
			curNote := utils.GetNoteStepFrom(lastNote, direction.step)
			rbSrcFilePath := filepath.Join(curSoundpackDir, subDirName, utils.GetNoteNameWithOctave(lastNote)+".wav")
			rbDstFilePath := filepath.Join(curSoundpackDir, subDirName, utils.GetNoteNameWithOctave(curNote)+".wav")
			if _, err := os.Stat(rbDstFilePath); err == nil {
				// File already exists
				if curNote == direction.endNote {
					break
				}
				lastNote = curNote
				continue
			}

			if err := shiftPitch(rbSrcFilePath, rbDstFilePath, direction.step); err != nil {
				return err
			}

			if curNote == direction.endNote {
				break
			}
			lastNote = curNote
		}
	}
	return nil
}

func copyFile(srcFilePath, dstFilePath string) error {
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstFilePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

// shiftPitch calls rubberband to shift the pitch of the audio file
// Parameters:
//   - semitoneToShift: the number of semitones to shift the pitch, positive value shifts up, negative value shifts down
func shiftPitch(src, dst string, semitoneToShift int) error {
	semitoneToShiftStr := fmt.Sprintf("%d", semitoneToShift)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", audio.RubberbandCommand, "-p", semitoneToShiftStr, src, dst)
	} else {
		cmd = exec.Command(audio.RubberbandCommand, "-p", semitoneToShiftStr, src, dst)
	}

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
