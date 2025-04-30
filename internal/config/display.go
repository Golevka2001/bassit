// Package config bassit/internal/config/display.go
package config

/* Spacing properties */
var (
    FretboardMarginLeft          = 5
    FretboardMarginRight         = 6
    FretboardMarginTop           = 5
    StringBaseNoteNameMarginLeft = 2
    NutWidth                     = 1
    // StringSpacing defines the spacing between strings
    StringSpacing = 2
    // StringMargin defines the spacing between the string and the fretboard edge
    StringMargin = 1
)

/* Characters for display */
var (
    FretboardVBorderChar      = '│'
    FretboardHBorderChar      = '─'
    FretboardULCornerChar     = '┌'
    FretboardLLCornerChar     = '└'
    FretWireChar              = '│'
    FretWireUpperChar         = '┬'
    FretWireLowerChar         = '┴'
    PositionMarkerChar        = '●'
    StringChar                = '═'
    StringOverFretChar        = '╪'
    PluckedStringChar         = '━'
    PluckedStringOverFretChar = '┿'
    PressedFretChar           = '◎'
)
