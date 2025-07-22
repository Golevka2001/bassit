package config

const (
	Version = "1.1.0"
)

const (
	// StringCnt defines the number of strings on a bass
	// Currently, 5-string bass is not supported
	StringCnt = 4

	// MaxFretCnt is NOT the number of frets to be displayed
	MaxFretCnt = 24

	// PluckTypeCount defines the number of pluck types
	PluckTypeCount = 6
)

const (
	StringVibDurMean = 850.0  // in ms
	StringVibDurStd  = 50.0   // in ms
	StringVibDurMin  = 700.0  // in ms
	StringVibDurMax  = 1000.0 // in ms
)
