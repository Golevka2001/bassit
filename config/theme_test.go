package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

func TestLoadTheme(t *testing.T) {
	// Save original ThemeDir to restore after testing
	originalThemeDir := ThemeDir
	defer func() { ThemeDir = originalThemeDir }()

	t.Run("Successfully load theme", func(t *testing.T) {
		// Create temporary directory as test theme directory
		tempDir := t.TempDir()
		ThemeDir = func() string { return tempDir }

		// Create test theme file
		testTheme := `
fretboard_border_color: "#333333"
fretboard_bg_color: "#222222"
nut_border_color: "#444444"
nut_bg_color: "#555555"
fretwire_color: "#666666"
string_color: "#777777"
inlay_color: "#888888"
pressed_fret_sign_color: "#999999"
plucked_string_sign_color: "#aaaaaa"
base_note_name_fg_color: "#bbbbbb"
base_note_name_bg_color: "#cccccc"
inlay_shape: "dot"

fretboard_v_border_char: "|"
fretboard_h_border_char: "-"
nut_v_border_char: "|"
nut_h_border_char: "="
nut_ul_corner_char: "+"
nut_ll_corner_char: "+"
nut_ur_corner_char: "+"
nut_lr_corner_char: "+"
fretwire_char: "|"
fretwire_top_border_char: "+"
fretwire_bottom_border_char: "+"
string_char: "-"
string_over_fretwire_char: "+"
string_over_boarder_char: "+"
vibrating_string_char: "~"
vibrating_string_over_fretwire_char: "#"
vibrating_string_over_boarder_char: "#"
inlay_char: "o"
pressed_fret_sign_char: "x"
not_plucked_string_char: " "
plucked_string_sign_char: ">"
`
		err := os.WriteFile(filepath.Join(tempDir, "test_theme.yaml"), []byte(testTheme), 0644)
		assert.NoError(t, err)

		// Test loading theme
		theme, err := LoadTheme("test_theme")
		assert.NoError(t, err)
		assert.NotNil(t, theme)

		// Verify converted theme properties
		assert.Equal(t, DotInlayShape, theme.InlayShape)
		assert.Equal(t, '|', theme.FretboardVBorderChar)
		assert.Equal(t, '-', theme.FretboardHBorderChar)
		assert.Equal(t, 'o', theme.InlayChar)
		assert.Equal(t, 'x', theme.PressedFretSignChar)
	})

	t.Run("Load non-existent theme", func(t *testing.T) {
		// Create temporary directory as test theme directory
		tempDir := t.TempDir()
		ThemeDir = func() string { return tempDir }

		// Test loading non-existent theme
		theme, err := LoadTheme("non_existent_theme")
		assert.Error(t, err)
		assert.Nil(t, theme)
	})

	t.Run("Automatically remove file extension", func(t *testing.T) {
		// Create temporary directory as test theme directory
		tempDir := t.TempDir()
		ThemeDir = func() string { return tempDir }

		// Create test theme file
		testTheme := `
fretboard_border_color: "#333333"
inlay_shape: "block"
`
		err := os.WriteFile(filepath.Join(tempDir, "extension_theme.yaml"), []byte(testTheme), 0644)
		assert.NoError(t, err)

		// Test loading theme with extension
		theme, err := LoadTheme("extension_theme.yaml")
		assert.NoError(t, err)
		assert.NotNil(t, theme)
		assert.Equal(t, BlockInlayShape, theme.InlayShape)

		// Test loading theme with .yml extension
		theme, err = LoadTheme("extension_theme.yml")
		assert.NoError(t, err)
		assert.NotNil(t, theme)
	})
}

func TestRawThemeConvert(t *testing.T) {
	t.Run("Test InlayShape conversion", func(t *testing.T) {
		testCases := []struct {
			name          string
			inlayShape    string
			expectedShape InlayShape
		}{
			{"Empty string", "", NoneInlayShape},
			{"none", "none", NoneInlayShape},
			{"dot", "dot", DotInlayShape},
			{"dots", "dots", DotInlayShape},
			{"block", "block", BlockInlayShape},
			{"blocks", "blocks", BlockInlayShape},
			{"Mixed case", "DoT", DotInlayShape},
			{"With spaces", " block ", BlockInlayShape},
			{"Unknown value", "unknown", DotInlayShape}, // Default to DotInlayShape
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				rt := &rawTheme{InlayShape: tc.inlayShape}
				theme := rt.convert()
				assert.Equal(t, tc.expectedShape, theme.InlayShape)
			})
		}
	})

	t.Run("Test character conversion", func(t *testing.T) {
		rt := &rawTheme{
			FretboardVBorderChar: "|",
			StringChar:           "-",
			InlayChar:            "o",
			// Test handling of empty string
			FretwireChar: "",
		}
		theme := rt.convert()
		assert.Equal(t, '|', theme.FretboardVBorderChar)
		assert.Equal(t, '-', theme.StringChar)
		assert.Equal(t, 'o', theme.InlayChar)
		assert.Equal(t, ' ', theme.FretwireChar) // Empty string should convert to space
	})

	t.Run("Test color conversion", func(t *testing.T) {
		rt := &rawTheme{
			FretboardBorderColor: "#333333",
			FretboardBgColor:     "#222222",
			StringColor:          "#777777",
		}
		theme := rt.convert()

		// Verify colors have been converted with specific values
		expectedFretboardBorder := lipgloss.Color("#333333")
		expectedFretboardBg := lipgloss.Color("#222222")
		expectedStringColor := lipgloss.Color("#777777")

		assert.Equal(t, expectedFretboardBorder, theme.FretboardBorderColor)
		assert.Equal(t, expectedFretboardBg, theme.FretboardBgColor)
		assert.Equal(t, expectedStringColor, theme.StringColor)
	})
}
