package common

import (
	"github.com/Golevka2001/bassit/audio"
	"github.com/Golevka2001/bassit/bass"
	"github.com/Golevka2001/bassit/config"
)

type UIContext struct {
	Audio     *audio.AudioManager
	Bass      *bass.BassModel
	Config    *config.Config
	Theme     *config.Theme
	Soundpack *config.SoundpackInfo
}
