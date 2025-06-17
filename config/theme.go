package config

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type InlayShape int

const (
	NoneInlayShape InlayShape = iota
	DotInlayShape
	BlockInlayShape
)

type Theme struct {
	FretboardBorderColor   color.Color
	FretboardBgColor       color.Color
	NutBorderColor         color.Color
	NutBgColor             color.Color
	FretwireColor          color.Color
	StringColor            color.Color
	InlayColor             color.Color
	PressedFretSignColor   color.Color
	PluckedStringSignColor color.Color
	BaseNoteNameFgColor    color.Color
	BaseNoteNameBgColor    color.Color

	InlayShape InlayShape

	FretboardVBorderChar            rune
	FretboardHBorderChar            rune
	NutVBorderChar                  rune
	NutHBorderChar                  rune
	NutULCornerChar                 rune
	NutLLCornerChar                 rune
	NutURCornerChar                 rune
	NutLRCornerChar                 rune
	FretwireChar                    rune
	FretwireTopBorderChar           rune
	FretwireBottomBorderChar        rune
	StringChar                      rune
	StringOverFretwireChar          rune
	StringOverBoarderChar           rune
	VibratingStringChar             rune
	VibratingStringOverFretwireChar rune
	VibratingStringOverBoarderChar  rune
	InlayChar                       rune
	PressedFretSignChar             rune
	NotPluckedStringChar            rune
	PluckedStringSignChar           rune
}

type rawTheme struct {
	FretboardBorderColor   string `yaml:"fretboard_border_color"`
	FretboardBgColor       string `yaml:"fretboard_bg_color"`
	NutBorderColor         string `yaml:"nut_border_color"`
	NutBgColor             string `yaml:"nut_bg_color"`
	FretwireColor          string `yaml:"fretwire_color"`
	StringColor            string `yaml:"string_color"`
	InlayColor             string `yaml:"inlay_color"`
	PressedFretSignColor   string `yaml:"pressed_fret_sign_color"`
	PluckedStringSignColor string `yaml:"plucked_string_sign_color"`
	BaseNoteNameFgColor    string `yaml:"base_note_name_fg_color"`
	BaseNoteNameBgColor    string `yaml:"base_note_name_bg_color"`
	InlayShape             string `yaml:"inlay_shape"`

	FretboardVBorderChar            string `yaml:"fretboard_v_border_char"`
	FretboardHBorderChar            string `yaml:"fretboard_h_border_char"`
	NutVBorderChar                  string `yaml:"nut_v_border_char"`
	NutHBorderChar                  string `yaml:"nut_h_border_char"`
	NutULCornerChar                 string `yaml:"nut_ul_corner_char"`
	NutLLCornerChar                 string `yaml:"nut_ll_corner_char"`
	NutURCornerChar                 string `yaml:"nut_ur_corner_char"`
	NutLRCornerChar                 string `yaml:"nut_lr_corner_char"`
	FretwireChar                    string `yaml:"fretwire_char"`
	FretwireTopBorderChar           string `yaml:"fretwire_top_border_char"`
	FretwireBottomBorderChar        string `yaml:"fretwire_bottom_border_char"`
	StringChar                      string `yaml:"string_char"`
	StringOverFretwireChar          string `yaml:"string_over_fretwire_char"`
	StringOverBoarderChar           string `yaml:"string_over_boarder_char"`
	VibratingStringChar             string `yaml:"vibrating_string_char"`
	VibratingStringOverFretwireChar string `yaml:"vibrating_string_over_fretwire_char"`
	VibratingStringOverBoarderChar  string `yaml:"vibrating_string_over_boarder_char"`
	InlayChar                       string `yaml:"inlay_char"`
	PressedFretSignChar             string `yaml:"pressed_fret_sign_char"`
	NotPluckedStringChar            string `yaml:"not_plucked_string_char"`
	PluckedStringSignChar           string `yaml:"plucked_string_sign_char"`
}

func LoadTheme(name string) (*Theme, error) {
	name = strings.TrimSuffix(strings.TrimSuffix(name, ".yaml"), ".yml")
	r, err := loadRawTheme(name)
	if err != nil {
		return nil, err
	}
	theme := r.convert()
	return &theme, nil
}

func loadRawTheme(name string) (*rawTheme, error) {
	v := viper.New()
	v.AddConfigPath(ThemeDir())
	v.SetConfigName(name)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading theme '%s': %w", name, err)
	}

	var rt rawTheme
	if err := v.Unmarshal(&rt, func(c *mapstructure.DecoderConfig) {
		c.TagName = "yaml"
	}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal theme '%s': %w", name, err)
	}
	return &rt, nil
}

// convert a `rawTheme` to a `Theme`
func (r *rawTheme) convert() Theme {
	// `inlayShape`
	var inlayShape InlayShape
	switch strings.ToLower(strings.TrimSpace(r.InlayShape)) {
	case "none", "":
		inlayShape = NoneInlayShape
	case "dot", "dots":
		inlayShape = DotInlayShape
	case "block", "blocks":
		inlayShape = BlockInlayShape
	default:
		inlayShape = DotInlayShape
	}

	// `getRune` is a helper function to extract the first rune from a string
	// or return a space if the string is empty
	getRune := func(s string) rune {
		runes := []rune(s)
		if len(runes) > 0 {
			return runes[0]
		}
		return ' '
	}

	return Theme{
		FretboardBorderColor:   lipgloss.Color(r.FretboardBorderColor),
		FretboardBgColor:       lipgloss.Color(r.FretboardBgColor),
		NutBorderColor:         lipgloss.Color(r.NutBorderColor),
		NutBgColor:             lipgloss.Color(r.NutBgColor),
		FretwireColor:          lipgloss.Color(r.FretwireColor),
		StringColor:            lipgloss.Color(r.StringColor),
		InlayColor:             lipgloss.Color(r.InlayColor),
		PressedFretSignColor:   lipgloss.Color(r.PressedFretSignColor),
		PluckedStringSignColor: lipgloss.Color(r.PluckedStringSignColor),
		BaseNoteNameFgColor:    lipgloss.Color(r.BaseNoteNameFgColor),
		BaseNoteNameBgColor:    lipgloss.Color(r.BaseNoteNameBgColor),

		InlayShape: inlayShape,

		FretboardVBorderChar:            getRune(r.FretboardVBorderChar),
		FretboardHBorderChar:            getRune(r.FretboardHBorderChar),
		NutVBorderChar:                  getRune(r.NutVBorderChar),
		NutHBorderChar:                  getRune(r.NutHBorderChar),
		NutULCornerChar:                 getRune(r.NutULCornerChar),
		NutLLCornerChar:                 getRune(r.NutLLCornerChar),
		NutURCornerChar:                 getRune(r.NutURCornerChar),
		NutLRCornerChar:                 getRune(r.NutLRCornerChar),
		FretwireChar:                    getRune(r.FretwireChar),
		FretwireTopBorderChar:           getRune(r.FretwireTopBorderChar),
		FretwireBottomBorderChar:        getRune(r.FretwireBottomBorderChar),
		StringChar:                      getRune(r.StringChar),
		StringOverFretwireChar:          getRune(r.StringOverFretwireChar),
		StringOverBoarderChar:           getRune(r.StringOverBoarderChar),
		VibratingStringChar:             getRune(r.VibratingStringChar),
		VibratingStringOverFretwireChar: getRune(r.VibratingStringOverFretwireChar),
		VibratingStringOverBoarderChar:  getRune(r.VibratingStringOverBoarderChar),
		InlayChar:                       getRune(r.InlayChar),
		PressedFretSignChar:             getRune(r.PressedFretSignChar),
		NotPluckedStringChar:            getRune(r.NotPluckedStringChar),
		PluckedStringSignChar:           getRune(r.PluckedStringSignChar),
	}
}
