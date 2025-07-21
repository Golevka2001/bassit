package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInlayShape(t *testing.T) {
	tests := []struct {
		name     string
		shape    InlayShape
		expected int
	}{
		{
			name:     "NoneInlayShape",
			shape:    NoneInlayShape,
			expected: 0,
		},
		{
			name:     "DotInlayShape",
			shape:    DotInlayShape,
			expected: 1,
		},
		{
			name:     "BlockInlayShape",
			shape:    BlockInlayShape,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, int(tt.shape))
		})
	}
}

func TestThemeStruct(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T)
	}{
		{
			name: "Theme struct has all required fields",
			validate: func(t *testing.T) {
				theme := Theme{
					InlayShape:                      DotInlayShape,
					FretboardVBorderChar:            '│',
					FretboardHBorderChar:            '─',
					NutVBorderChar:                  '│',
					NutHBorderChar:                  '─',
					NutULCornerChar:                 '┌',
					NutLLCornerChar:                 '└',
					NutURCornerChar:                 '┬',
					NutLRCornerChar:                 '┴',
					FretwireChar:                    '║',
					FretwireTopBorderChar:           '╥',
					FretwireBottomBorderChar:        '╨',
					StringChar:                      '━',
					StringOverFretwireChar:          '╫',
					StringOverBoarderChar:           '┿',
					VibratingStringChar:             '═',
					VibratingStringOverFretwireChar: '╬',
					VibratingStringOverBoarderChar:  '╪',
					InlayChar:                       '●',
					PressedFretSignChar:             '●',
					NotPluckedStringChar:            '░',
					PluckedStringSignChar:           '█',
				}

				assert.Equal(t, DotInlayShape, theme.InlayShape)
				assert.Equal(t, '│', theme.FretboardVBorderChar)
				assert.Equal(t, '─', theme.FretboardHBorderChar)
				assert.Equal(t, '│', theme.NutVBorderChar)
				assert.Equal(t, '─', theme.NutHBorderChar)
				assert.Equal(t, '┌', theme.NutULCornerChar)
				assert.Equal(t, '└', theme.NutLLCornerChar)
				assert.Equal(t, '┬', theme.NutURCornerChar)
				assert.Equal(t, '┴', theme.NutLRCornerChar)
				assert.Equal(t, '║', theme.FretwireChar)
				assert.Equal(t, '╥', theme.FretwireTopBorderChar)
				assert.Equal(t, '╨', theme.FretwireBottomBorderChar)
				assert.Equal(t, '━', theme.StringChar)
				assert.Equal(t, '╫', theme.StringOverFretwireChar)
				assert.Equal(t, '┿', theme.StringOverBoarderChar)
				assert.Equal(t, '═', theme.VibratingStringChar)
				assert.Equal(t, '╬', theme.VibratingStringOverFretwireChar)
				assert.Equal(t, '╪', theme.VibratingStringOverBoarderChar)
				assert.Equal(t, '●', theme.InlayChar)
				assert.Equal(t, '●', theme.PressedFretSignChar)
				assert.Equal(t, '░', theme.NotPluckedStringChar)
				assert.Equal(t, '█', theme.PluckedStringSignChar)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t)
		})
	}
}

func TestRawThemeStruct(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T)
	}{
		{
			name: "rawTheme struct has all required fields",
			validate: func(t *testing.T) {
				rawTheme := rawTheme{
					FretboardBorderColor:            "#ffffff",
					FretboardBgColor:                "#000000",
					NutBorderColor:                  "#cccccc",
					NutBgColor:                      "#f0f0f0",
					FretwireColor:                   "#e1d6c9",
					StringColor:                     "#dedcda",
					InlayColor:                      "#e4e3df",
					PressedFretSignColor:            "#FF6666",
					PluckedStringSignColor:          "#FF6666",
					BaseNoteNameFgColor:             "#000000",
					BaseNoteNameBgColor:             "#ffffff",
					InlayShape:                      "dot",
					FretboardVBorderChar:            "│",
					FretboardHBorderChar:            "─",
					NutVBorderChar:                  "│",
					NutHBorderChar:                  "─",
					NutULCornerChar:                 "┌",
					NutLLCornerChar:                 "└",
					NutURCornerChar:                 "┬",
					NutLRCornerChar:                 "┴",
					FretwireChar:                    "║",
					FretwireTopBorderChar:           "╥",
					FretwireBottomBorderChar:        "╨",
					StringChar:                      "━",
					StringOverFretwireChar:          "╫",
					StringOverBoarderChar:           "┿",
					VibratingStringChar:             "═",
					VibratingStringOverFretwireChar: "╬",
					VibratingStringOverBoarderChar:  "╪",
					InlayChar:                       "●",
					PressedFretSignChar:             "●",
					NotPluckedStringChar:            "░",
					PluckedStringSignChar:           "█",
				}

				assert.Equal(t, "#ffffff", rawTheme.FretboardBorderColor)
				assert.Equal(t, "#000000", rawTheme.FretboardBgColor)
				assert.Equal(t, "#cccccc", rawTheme.NutBorderColor)
				assert.Equal(t, "#f0f0f0", rawTheme.NutBgColor)
				assert.Equal(t, "#e1d6c9", rawTheme.FretwireColor)
				assert.Equal(t, "#dedcda", rawTheme.StringColor)
				assert.Equal(t, "#e4e3df", rawTheme.InlayColor)
				assert.Equal(t, "#FF6666", rawTheme.PressedFretSignColor)
				assert.Equal(t, "#FF6666", rawTheme.PluckedStringSignColor)
				assert.Equal(t, "#000000", rawTheme.BaseNoteNameFgColor)
				assert.Equal(t, "#ffffff", rawTheme.BaseNoteNameBgColor)
				assert.Equal(t, "dot", rawTheme.InlayShape)
				assert.Equal(t, "│", rawTheme.FretboardVBorderChar)
				assert.Equal(t, "─", rawTheme.FretboardHBorderChar)
				assert.Equal(t, "│", rawTheme.NutVBorderChar)
				assert.Equal(t, "─", rawTheme.NutHBorderChar)
				assert.Equal(t, "┌", rawTheme.NutULCornerChar)
				assert.Equal(t, "└", rawTheme.NutLLCornerChar)
				assert.Equal(t, "┬", rawTheme.NutURCornerChar)
				assert.Equal(t, "┴", rawTheme.NutLRCornerChar)
				assert.Equal(t, "║", rawTheme.FretwireChar)
				assert.Equal(t, "╥", rawTheme.FretwireTopBorderChar)
				assert.Equal(t, "╨", rawTheme.FretwireBottomBorderChar)
				assert.Equal(t, "━", rawTheme.StringChar)
				assert.Equal(t, "╫", rawTheme.StringOverFretwireChar)
				assert.Equal(t, "┿", rawTheme.StringOverBoarderChar)
				assert.Equal(t, "═", rawTheme.VibratingStringChar)
				assert.Equal(t, "╬", rawTheme.VibratingStringOverFretwireChar)
				assert.Equal(t, "╪", rawTheme.VibratingStringOverBoarderChar)
				assert.Equal(t, "●", rawTheme.InlayChar)
				assert.Equal(t, "●", rawTheme.PressedFretSignChar)
				assert.Equal(t, "░", rawTheme.NotPluckedStringChar)
				assert.Equal(t, "█", rawTheme.PluckedStringSignChar)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t)
		})
	}
}

func TestRawThemeConvert(t *testing.T) {
	tests := []struct {
		name     string
		rawTheme rawTheme
		validate func(t *testing.T, theme Theme)
	}{
		{
			name: "convert with dot inlay shape",
			rawTheme: rawTheme{
				FretboardBorderColor:            "#ffffff",
				FretboardBgColor:                "#000000",
				NutBorderColor:                  "#cccccc",
				NutBgColor:                      "#f0f0f0",
				FretwireColor:                   "#e1d6c9",
				StringColor:                     "#dedcda",
				InlayColor:                      "#e4e3df",
				PressedFretSignColor:            "#FF6666",
				PluckedStringSignColor:          "#FF6666",
				BaseNoteNameFgColor:             "#000000",
				BaseNoteNameBgColor:             "#ffffff",
				InlayShape:                      "dot",
				FretboardVBorderChar:            "│",
				FretboardHBorderChar:            "─",
				NutVBorderChar:                  "│",
				NutHBorderChar:                  "─",
				NutULCornerChar:                 "┌",
				NutLLCornerChar:                 "└",
				NutURCornerChar:                 "┬",
				NutLRCornerChar:                 "┴",
				FretwireChar:                    "║",
				FretwireTopBorderChar:           "╥",
				FretwireBottomBorderChar:        "╨",
				StringChar:                      "━",
				StringOverFretwireChar:          "╫",
				StringOverBoarderChar:           "┿",
				VibratingStringChar:             "═",
				VibratingStringOverFretwireChar: "╬",
				VibratingStringOverBoarderChar:  "╪",
				InlayChar:                       "●",
				PressedFretSignChar:             "●",
				NotPluckedStringChar:            "░",
				PluckedStringSignChar:           "█",
			},
			validate: func(t *testing.T, theme Theme) {
				assert.Equal(t, DotInlayShape, theme.InlayShape)
				assert.Equal(t, '│', theme.FretboardVBorderChar)
				assert.Equal(t, '─', theme.FretboardHBorderChar)
				assert.Equal(t, '│', theme.NutVBorderChar)
				assert.Equal(t, '─', theme.NutHBorderChar)
				assert.Equal(t, '┌', theme.NutULCornerChar)
				assert.Equal(t, '└', theme.NutLLCornerChar)
				assert.Equal(t, '┬', theme.NutURCornerChar)
				assert.Equal(t, '┴', theme.NutLRCornerChar)
				assert.Equal(t, '║', theme.FretwireChar)
				assert.Equal(t, '╥', theme.FretwireTopBorderChar)
				assert.Equal(t, '╨', theme.FretwireBottomBorderChar)
				assert.Equal(t, '━', theme.StringChar)
				assert.Equal(t, '╫', theme.StringOverFretwireChar)
				assert.Equal(t, '┿', theme.StringOverBoarderChar)
				assert.Equal(t, '═', theme.VibratingStringChar)
				assert.Equal(t, '╬', theme.VibratingStringOverFretwireChar)
				assert.Equal(t, '╪', theme.VibratingStringOverBoarderChar)
				assert.Equal(t, '●', theme.InlayChar)
				assert.Equal(t, '●', theme.PressedFretSignChar)
				assert.Equal(t, '░', theme.NotPluckedStringChar)
				assert.Equal(t, '█', theme.PluckedStringSignChar)
			},
		},
		{
			name: "convert with block inlay shape",
			rawTheme: rawTheme{
				InlayShape:                      "block",
				FretboardVBorderChar:            "│",
				FretboardHBorderChar:            "─",
				NutVBorderChar:                  "│",
				NutHBorderChar:                  "─",
				NutULCornerChar:                 "┌",
				NutLLCornerChar:                 "└",
				NutURCornerChar:                 "┬",
				NutLRCornerChar:                 "┴",
				FretwireChar:                    "║",
				FretwireTopBorderChar:           "╥",
				FretwireBottomBorderChar:        "╨",
				StringChar:                      "━",
				StringOverFretwireChar:          "╫",
				StringOverBoarderChar:           "┿",
				VibratingStringChar:             "═",
				VibratingStringOverFretwireChar: "╬",
				VibratingStringOverBoarderChar:  "╪",
				InlayChar:                       "█",
				PressedFretSignChar:             "●",
				NotPluckedStringChar:            "░",
				PluckedStringSignChar:           "◆",
			},
			validate: func(t *testing.T, theme Theme) {
				assert.Equal(t, BlockInlayShape, theme.InlayShape)
				assert.Equal(t, '│', theme.FretboardVBorderChar)
				assert.Equal(t, '─', theme.FretboardHBorderChar)
				assert.Equal(t, '│', theme.NutVBorderChar)
				assert.Equal(t, '─', theme.NutHBorderChar)
				assert.Equal(t, '┌', theme.NutULCornerChar)
				assert.Equal(t, '└', theme.NutLLCornerChar)
				assert.Equal(t, '┬', theme.NutURCornerChar)
				assert.Equal(t, '┴', theme.NutLRCornerChar)
				assert.Equal(t, '║', theme.FretwireChar)
				assert.Equal(t, '╥', theme.FretwireTopBorderChar)
				assert.Equal(t, '╨', theme.FretwireBottomBorderChar)
				assert.Equal(t, '━', theme.StringChar)
				assert.Equal(t, '╫', theme.StringOverFretwireChar)
				assert.Equal(t, '┿', theme.StringOverBoarderChar)
				assert.Equal(t, '═', theme.VibratingStringChar)
				assert.Equal(t, '╬', theme.VibratingStringOverFretwireChar)
				assert.Equal(t, '╪', theme.VibratingStringOverBoarderChar)
				assert.Equal(t, '█', theme.InlayChar)
				assert.Equal(t, '●', theme.PressedFretSignChar)
				assert.Equal(t, '░', theme.NotPluckedStringChar)
				assert.Equal(t, '◆', theme.PluckedStringSignChar)
			},
		},
		{
			name: "convert with none inlay shape",
			rawTheme: rawTheme{
				InlayShape:            "none",
				FretboardVBorderChar:  "│",
				FretboardHBorderChar:  "─",
				StringChar:            "━",
				InlayChar:             " ",
				PressedFretSignChar:   "●",
				PluckedStringSignChar: "█",
			},
			validate: func(t *testing.T, theme Theme) {
				assert.Equal(t, NoneInlayShape, theme.InlayShape)
				assert.Equal(t, '│', theme.FretboardVBorderChar)
				assert.Equal(t, '─', theme.FretboardHBorderChar)
				assert.Equal(t, '━', theme.StringChar)
				assert.Equal(t, ' ', theme.InlayChar)
				assert.Equal(t, '●', theme.PressedFretSignChar)
				assert.Equal(t, '█', theme.PluckedStringSignChar)
			},
		},
		{
			name: "convert with invalid inlay shape - defaults to dot",
			rawTheme: rawTheme{
				InlayShape:            "invalid",
				FretboardVBorderChar:  "│",
				FretboardHBorderChar:  "─",
				StringChar:            "━",
				InlayChar:             "●",
				PressedFretSignChar:   "●",
				PluckedStringSignChar: "█",
			},
			validate: func(t *testing.T, theme Theme) {
				assert.Equal(t, DotInlayShape, theme.InlayShape)
			},
		},
		{
			name: "convert with empty strings - defaults to space",
			rawTheme: rawTheme{
				InlayShape:                      "dot",
				FretboardVBorderChar:            "",
				FretboardHBorderChar:            "",
				NutVBorderChar:                  "",
				NutHBorderChar:                  "",
				NutULCornerChar:                 "",
				NutLLCornerChar:                 "",
				NutURCornerChar:                 "",
				NutLRCornerChar:                 "",
				FretwireChar:                    "",
				FretwireTopBorderChar:           "",
				FretwireBottomBorderChar:        "",
				StringChar:                      "",
				StringOverFretwireChar:          "",
				StringOverBoarderChar:           "",
				VibratingStringChar:             "",
				VibratingStringOverFretwireChar: "",
				VibratingStringOverBoarderChar:  "",
				InlayChar:                       "",
				PressedFretSignChar:             "",
				NotPluckedStringChar:            "",
				PluckedStringSignChar:           "",
			},
			validate: func(t *testing.T, theme Theme) {
				assert.Equal(t, DotInlayShape, theme.InlayShape)
				assert.Equal(t, ' ', theme.FretboardVBorderChar)
				assert.Equal(t, ' ', theme.FretboardHBorderChar)
				assert.Equal(t, ' ', theme.NutVBorderChar)
				assert.Equal(t, ' ', theme.NutHBorderChar)
				assert.Equal(t, ' ', theme.NutULCornerChar)
				assert.Equal(t, ' ', theme.NutLLCornerChar)
				assert.Equal(t, ' ', theme.NutURCornerChar)
				assert.Equal(t, ' ', theme.NutLRCornerChar)
				assert.Equal(t, ' ', theme.FretwireChar)
				assert.Equal(t, ' ', theme.FretwireTopBorderChar)
				assert.Equal(t, ' ', theme.FretwireBottomBorderChar)
				assert.Equal(t, ' ', theme.StringChar)
				assert.Equal(t, ' ', theme.StringOverFretwireChar)
				assert.Equal(t, ' ', theme.StringOverBoarderChar)
				assert.Equal(t, ' ', theme.VibratingStringChar)
				assert.Equal(t, ' ', theme.VibratingStringOverFretwireChar)
				assert.Equal(t, ' ', theme.VibratingStringOverBoarderChar)
				assert.Equal(t, ' ', theme.InlayChar)
				assert.Equal(t, ' ', theme.PressedFretSignChar)
				assert.Equal(t, ' ', theme.NotPluckedStringChar)
				assert.Equal(t, ' ', theme.PluckedStringSignChar)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme := tt.rawTheme.convert()
			if tt.validate != nil {
				tt.validate(t, theme)
			}
		})
	}
}

func TestLoadTheme(t *testing.T) {
	originalBaseDir := BaseDir
	defer func() { BaseDir = originalBaseDir }()

	tests := []struct {
		name     string
		setup    func() (string, string)
		expected error
		validate func(t *testing.T, theme *Theme, err error, tempDir string)
	}{
		{
			name: "load valid theme",
			setup: func() (string, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_theme_test_*")

				themeDir := filepath.Join(tempDir, "themes")
				os.MkdirAll(themeDir, 0755)

				themeContent := `
fretboard_border_color: "#ffffff"
fretboard_bg_color: "#000000"
nut_border_color: "#cccccc"
nut_bg_color: "#f0f0f0"
fretwire_color: "#e1d6c9"
string_color: "#dedcda"
inlay_color: "#e4e3df"
pressed_fret_sign_color: "#FF6666"
plucked_string_sign_color: "#FF6666"
base_note_name_fg_color: "#000000"
base_note_name_bg_color: "#ffffff"

inlay_shape: "dot"

fretboard_v_border_char: "│"
fretboard_h_border_char: "─"
nut_v_border_char: "│"
nut_h_border_char: "─"
nut_ul_corner_char: "┌"
nut_ll_corner_char: "└"
nut_ur_corner_char: "┬"
nut_lr_corner_char: "┴"
fretwire_char: "║"
fretwire_top_border_char: "╥"
fretwire_bottom_border_char: "╨"
string_char: "━"
string_over_fretwire_char: "╫"
string_over_boarder_char: "┿"
vibrating_string_char: "═"
vibrating_string_over_fretwire_char: "╬"
vibrating_string_over_boarder_char: "╪"
inlay_char: "●"
pressed_fret_sign_char: "●"
not_plucked_string_char: "░"
plucked_string_sign_char: "█"
`
				themeFile := filepath.Join(themeDir, "test_theme.yaml")
				os.WriteFile(themeFile, []byte(themeContent), 0644)

				BaseDir = tempDir
				return "test_theme", tempDir
			},
			expected: nil,
			validate: func(t *testing.T, theme *Theme, err error, tempDir string) {
				assert.NoError(t, err)
				assert.NotNil(t, theme)
				assert.Equal(t, DotInlayShape, theme.InlayShape)
				assert.Equal(t, '│', theme.FretboardVBorderChar)
				assert.Equal(t, '─', theme.FretboardHBorderChar)
				assert.Equal(t, '│', theme.NutVBorderChar)
				assert.Equal(t, '─', theme.NutHBorderChar)
				assert.Equal(t, '┌', theme.NutULCornerChar)
				assert.Equal(t, '└', theme.NutLLCornerChar)
				assert.Equal(t, '┬', theme.NutURCornerChar)
				assert.Equal(t, '┴', theme.NutLRCornerChar)
				assert.Equal(t, '║', theme.FretwireChar)
				assert.Equal(t, '╥', theme.FretwireTopBorderChar)
				assert.Equal(t, '╨', theme.FretwireBottomBorderChar)
				assert.Equal(t, '━', theme.StringChar)
				assert.Equal(t, '╫', theme.StringOverFretwireChar)
				assert.Equal(t, '┿', theme.StringOverBoarderChar)
				assert.Equal(t, '═', theme.VibratingStringChar)
				assert.Equal(t, '╬', theme.VibratingStringOverFretwireChar)
				assert.Equal(t, '╪', theme.VibratingStringOverBoarderChar)
				assert.Equal(t, '●', theme.InlayChar)
				assert.Equal(t, '●', theme.PressedFretSignChar)
				assert.Equal(t, '░', theme.NotPluckedStringChar)
				assert.Equal(t, '█', theme.PluckedStringSignChar)
			},
		},
		{
			name: "load theme with .yaml extension",
			setup: func() (string, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_theme_test_*")

				themeDir := filepath.Join(tempDir, "themes")
				os.MkdirAll(themeDir, 0755)

				themeContent := `
inlay_shape: "block"
fretboard_v_border_char: "│"
fretboard_h_border_char: "─"
nut_v_border_char: "│"
nut_h_border_char: "─"
fretwire_char: "║"
string_char: "━"
inlay_char: "█"
pressed_fret_sign_char: "●"
plucked_string_sign_char: "█"
`
				themeFile := filepath.Join(themeDir, "test_theme.yaml")
				os.WriteFile(themeFile, []byte(themeContent), 0644)

				BaseDir = tempDir
				return "test_theme.yaml", tempDir
			},
			expected: nil,
			validate: func(t *testing.T, theme *Theme, err error, tempDir string) {
				assert.NoError(t, err)
				assert.NotNil(t, theme)
				assert.Equal(t, BlockInlayShape, theme.InlayShape)
				assert.Equal(t, '│', theme.FretboardVBorderChar)
				assert.Equal(t, '─', theme.FretboardHBorderChar)
				assert.Equal(t, '│', theme.NutVBorderChar)
				assert.Equal(t, '─', theme.NutHBorderChar)
				assert.Equal(t, '║', theme.FretwireChar)
				assert.Equal(t, '━', theme.StringChar)
				assert.Equal(t, '█', theme.InlayChar)
				assert.Equal(t, '●', theme.PressedFretSignChar)
				assert.Equal(t, '█', theme.PluckedStringSignChar)
			},
		},
		{
			name: "theme file not found",
			setup: func() (string, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_theme_test_*")

				themeDir := filepath.Join(tempDir, "themes")
				os.MkdirAll(themeDir, 0755)

				BaseDir = tempDir
				return "nonexistent", tempDir
			},
			expected: assert.AnError,
			validate: func(t *testing.T, theme *Theme, err error, tempDir string) {
				assert.Error(t, err)
				assert.Nil(t, theme)
				assert.Contains(t, err.Error(), "error reading theme")
			},
		},
		{
			name: "invalid YAML syntax",
			setup: func() (string, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_theme_test_*")

				themeDir := filepath.Join(tempDir, "themes")
				os.MkdirAll(themeDir, 0755)

				themeContent := `
fretboard_border_color: "#ffffff"
invalid yaml syntax here: [unclosed bracket
inlay_shape: "dot"
fretboard_v_border_char: "│"
`
				themeFile := filepath.Join(themeDir, "invalid_theme.yaml")
				os.WriteFile(themeFile, []byte(themeContent), 0644)

				BaseDir = tempDir
				return "invalid_theme", tempDir
			},
			expected: assert.AnError,
			validate: func(t *testing.T, theme *Theme, err error, tempDir string) {
				assert.Error(t, err)
				assert.Nil(t, theme)
				assert.Contains(t, err.Error(), "error reading theme")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			themeName, tempDir := tt.setup()
			defer os.RemoveAll(tempDir)

			theme, err := LoadTheme(themeName)

			if tt.validate != nil {
				tt.validate(t, theme, err, tempDir)
			}
		})
	}
}
