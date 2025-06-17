package audio

import (
	"github.com/Golevka2001/bassit/config"
)

const (
	SampleRate   = 48000 // in Hz
	ChannelCount = 2     // stereo, or 1 for mono
)

var (
	RubberbandCommand     = "rubberband-r3"
	NoteSampleDir         = config.BaseDir() + "/audio/bass/pluck/default/"
	SrcBassSampleNoteName = "C2"
)
