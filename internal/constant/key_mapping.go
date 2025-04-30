// Package constant bassit/internal/constant/key_mapping.go
package constant

var (
	// Key2Fret maps the keys to the corresponding fret numbers for each string
	Key2Fret = [StringNum]map[rune]int{
		// 1st string (G)
		{'1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, '0': 10, '-': 11},
		// 2nd string (D)
		{'q': 1, 'w': 2, 'e': 3, 'r': 4, 't': 5, 'y': 6, 'u': 7, 'i': 8, 'o': 9, 'p': 10},
		// 3rd string (A)
		{'a': 1, 's': 2, 'd': 3, 'f': 4, 'g': 5, 'h': 6, 'j': 7, 'k': 8, 'l': 9},
		// 4th string (E)
		{'z': 1, 'x': 2, 'c': 3, 'v': 4, 'b': 5, 'n': 6, 'm': 7, ',': 8},
	}

	// ShiftedKey2Fret maps the shifted keys to the corresponding fret numbers for each string
	ShiftedKey2Fret = [StringNum]map[rune]int{
		// 1st string (G)
		{'!': 1, '@': 2, '#': 3, '$': 4, '%': 5, '^': 6, '&': 7, '*': 8, '(': 9, ')': 10, '_': 11},
		// 2nd string (D)
		{'Q': 1, 'W': 2, 'E': 3, 'R': 4, 'T': 5, 'Y': 6, 'U': 7, 'I': 8, 'O': 9, 'P': 10},
		// 3rd string (A)
		{'A': 1, 'S': 2, 'D': 3, 'F': 4, 'G': 5, 'H': 6, 'J': 7, 'K': 8, 'L': 9},
		// 4th string (E)
		{'Z': 1, 'X': 2, 'C': 3, 'V': 4, 'B': 5, 'N': 6, 'M': 7, '<': 8},
	}

	// Key2Pluck maps the keys to the corresponding plucking actions for each string
	Key2Pluck = [StringNum]map[rune]int{
		// 1st string (G)
		{'=': 1, '\b': 2},
		// 2nd string (D)
		{'[': 1, ']': 2},
		// 3rd string (A)
		{';': 1, '\'': 2},
		// 4th string (E)
		{'.': 1, '/': 2},
	}

	// ShiftedKey2Pluck maps the shifted keys to the corresponding plucking actions for each string
	ShiftedKey2Pluck = [StringNum]map[rune]int{
		// 1st string (G)
		{'+': 1, '\b': 2},
		// 2nd string (D)
		{'{': 1, '}': 2},
		// 3rd string (A)
		{':': 1, '"': 2},
		// 4th string (E)
		{'>': 1, '?': 2},
	}
)
