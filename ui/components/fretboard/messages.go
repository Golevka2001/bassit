package fretboard

type PressFretMsg struct {
	StringIdx int
	FretIdx   int
}

type ReleaseFretMsg struct {
	StringIdx int
	FretIdx   int
}

type PluckStringMsg struct {
	StringIdx int
	FretIdx   int
	Position  int
}

type RestorePluckedStringMsg struct {
	StringIdx int
	FretIdx   int
}
