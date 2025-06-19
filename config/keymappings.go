package config

var (
	KeyToFretboardPos = map[string]struct {
		StringIdx int
		FretIdx   int
	}{
		// The 1st string (G string)
		"1": {StringIdx: 0, FretIdx: 1},
		"2": {StringIdx: 0, FretIdx: 2},
		"3": {StringIdx: 0, FretIdx: 3},
		"4": {StringIdx: 0, FretIdx: 4},
		"5": {StringIdx: 0, FretIdx: 5},
		"6": {StringIdx: 0, FretIdx: 6},
		"7": {StringIdx: 0, FretIdx: 7},
		"8": {StringIdx: 0, FretIdx: 8},
		"9": {StringIdx: 0, FretIdx: 9},
		"0": {StringIdx: 0, FretIdx: 10},
		"-": {StringIdx: 0, FretIdx: 11},
		// The 2nd string (D string)
		"q": {StringIdx: 1, FretIdx: 1},
		"w": {StringIdx: 1, FretIdx: 2},
		"e": {StringIdx: 1, FretIdx: 3},
		"r": {StringIdx: 1, FretIdx: 4},
		"t": {StringIdx: 1, FretIdx: 5},
		"y": {StringIdx: 1, FretIdx: 6},
		"u": {StringIdx: 1, FretIdx: 7},
		"i": {StringIdx: 1, FretIdx: 8},
		"o": {StringIdx: 1, FretIdx: 9},
		"p": {StringIdx: 1, FretIdx: 10},
		// The 3rd string (A string)
		"a": {StringIdx: 2, FretIdx: 1},
		"s": {StringIdx: 2, FretIdx: 2},
		"d": {StringIdx: 2, FretIdx: 3},
		"f": {StringIdx: 2, FretIdx: 4},
		"g": {StringIdx: 2, FretIdx: 5},
		"h": {StringIdx: 2, FretIdx: 6},
		"j": {StringIdx: 2, FretIdx: 7},
		"k": {StringIdx: 2, FretIdx: 8},
		"l": {StringIdx: 2, FretIdx: 9},
		// The 4th string (E string)
		"z": {StringIdx: 3, FretIdx: 1},
		"x": {StringIdx: 3, FretIdx: 2},
		"c": {StringIdx: 3, FretIdx: 3},
		"v": {StringIdx: 3, FretIdx: 4},
		"b": {StringIdx: 3, FretIdx: 5},
		"n": {StringIdx: 3, FretIdx: 6},
		"m": {StringIdx: 3, FretIdx: 7},
		",": {StringIdx: 3, FretIdx: 8},
	}

	KeyToPluckState = map[string]struct {
		StringIdx int
		// Position can be `0` or `1`
		Position int
	}{
		// The 1st string (G string)
		"=":         {StringIdx: 0, Position: 0},
		"backspace": {StringIdx: 0, Position: 1},
		// The 2nd string (D string)
		"[": {StringIdx: 1, Position: 0},
		"]": {StringIdx: 1, Position: 1},
		// The 3rd string (A string)
		";": {StringIdx: 2, Position: 0},
		"'": {StringIdx: 2, Position: 1},
		// The 4th string (E string)
		".": {StringIdx: 3, Position: 0},
		"/": {StringIdx: 3, Position: 1},
	}
)
