package util

import (
	C "github.com/Golevka2001/bassit/constant"
)

func MapRawcodeToPressedPos() {
	for stringIdx, keyMap := range C.KeyToPressedPos {
		for key, fretIdx := range keyMap {
			rawcode, ok := C.KeyToRawcode[key]
			if C.OS == "darwin" {
				rawcode, ok = C.KeyToRawcodeForDarwin[key]
			}
			if !ok {
				continue
			}

			C.RawcodeToPressedPos[rawcode] = C.PressedPos{
				String: stringIdx,
				Fret:   fretIdx,
			}
		}
	}
}

func MapRawcodeToPluckedString() {
	for key, stringIdx := range C.KeyToPluckedString {
		rawcode, ok := C.KeyToRawcode[key]
		if C.OS == "darwin" {
			rawcode, ok = C.KeyToRawcodeForDarwin[key]
		}
		if !ok {
			continue
		}

		C.RawcodeToPluckedString[rawcode] = stringIdx
	}
}
