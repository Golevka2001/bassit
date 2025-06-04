package constant

const (
	SampleRate   = 48000 // in Hz
	ChannelCount = 2     // stereo, or 1 for mono
	// BitDepth     = 24    // in bits
)

var (
	RubberbandCommand     = "rubberband-r3"
	NoteSampleDir         = BaseDir() + "/audio/bass/pluck/default/"
	SrcBassSampleNoteName = "C2"
)
