package bass

type FretboardPosition struct {
	// StringIdx is the index from the thinnest string to the thickest string (ranges from 0 to config.StringCnt-1)
	// 0 represents the thinnest string, config.StringCnt-1 represents the thickest string
	StringIdx int

	// FretIdx is the index from the nut to the last fret (ranges from 0 to config.MaxFretCnt)
	// 0 represents the nut, other values represent the frets
	FretIdx int
}

// PluckType defines different types of techniques for plucking the string
type PluckType int

const (
	PluckTypeNormal1 PluckType = iota // recommend defining this as "mp dynamic"
	PluckTypeNormal2                  // recommend defining this as "mf dynamic"
	PluckTypeSlap1                    // recommend defining this as "slap"
	PluckTypeSlap2                    // recommend defining this as "pop"
	PluckTypeMute1                    // recommend defining this as "left-hand mute"/"dead note"/"ghost note"
	PluckTypeMute2                    // recommend defining this as "palm mute"
)

func (t PluckType) String() string {
	return []string{"normal1", "normal2", "slap1", "slap2", "mute1", "mute2"}[t]
}
