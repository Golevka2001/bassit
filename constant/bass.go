package constant

const (
	// StringCnt defines the number of strings on a bass. Currently, 5-string bass is not supported
	StringCnt = 4
	// MaxFretCnt is NOT the number of frets to be displayed
	MaxFretCnt = 24
	// DisplayedFretNum defines the number of frets to be displayed
	DisplayedFretNum = 12
	// StandardPitch for tuning bass strings, default is A4 = 440Hz
	StandardPitch = 440.0 // A4 in Hz
)

var (
	// StandardTuning for a 4-string bass
	StandardTuning = []string{"G2", "D2", "A1", "E1"}
)
