// bassit/internal/config/display.go
package config

const (
	/* Spacing properties */
	FretboardMarginLeft  = 5
	FretboardMarginRight = 5
	FretboardMarginTop   = 5
	StringNameMarginLeft = 2
	NutWidth             = 1
	StringSpacing        = 2 // the spacing between strings
	StringMargin         = 1 // the spacing between the string and the fretboard edge

	/* Characters for display */
	FretboardVBorderChar  = '│'
	FretboardHBorderChar  = '─'
	FretboardULCornerChar = '┌'
	FretboardLLCornerChar = '└'
	FretWireChar          = '│'
	FretWireUpperChar     = '┬'
	FretWireLowerChar     = '┴'
	PositionMarkerChar    = '●'
	StringChar            = '═'
	StringOverFretChar    = '╪'
	PressedFretChar       = '◎'
)
