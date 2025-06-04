package config

import (
	"fmt"
	"path/filepath"
	"strings"

	C "github.com/Golevka2001/bassit/constant"
	
	"github.com/gdamore/tcell/v2"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type InlayShape int

const (
	NoneInlayShape InlayShape = iota
	DotInlayShape
	BlockInlayShape
)

var (
	// tViper is a viper instance for theme configuration
	tViper   = viper.New()
	RawTheme = rawTheme{}
	Theme    = theme{}
)

type rawTheme struct {
	// Color properties
	TitleFgColor           string `yaml:"title_fg_color"`
	TitleBgColor           string `yaml:"title_bg_color"`
	FretboardBorderColor   string `yaml:"fretboard_border_color"`
	FretboardBgColor       string `yaml:"fretboard_bg_color"`
	NutBorderColor         string `yaml:"nut_border_color"`
	NutBgColor             string `yaml:"nut_bg_color"`
	FretWireColor          string `yaml:"fret_wire_color"`
	StringColor            string `yaml:"string_color"`
	InlayColor             string `yaml:"inlay_color"`
	PressedFretSignColor   string `yaml:"pressed_fret_sign_color"`
	PluckedStringSignColor string `yaml:"plucked_string_sign_color"`
	BaseNoteNameFgColor    string `yaml:"base_note_name_fg_color"`
	BaseNoteNameBgColor    string `yaml:"base_note_name_bg_color"`

	InlayShape string `yaml:"inlay_shape"`

	// Characters for display
	FretboardVBorderChar            string `yaml:"fretboard_v_border_char"`
	FretboardHBorderChar            string `yaml:"fretboard_h_border_char"`
	NutVBorderChar                  string `yaml:"nut_v_border_char"`
	NutHBorderChar                  string `yaml:"nut_h_border_char"`
	NutULCornerChar                 string `yaml:"nut_ul_corner_char"`
	NutLLCornerChar                 string `yaml:"nut_ll_corner_char"`
	NutURCornerChar                 string `yaml:"nut_ur_corner_char"`
	NutLRCornerChar                 string `yaml:"nut_lr_corner_char"`
	FretWireChar                    string `yaml:"fret_wire_char"`
	FretWireUpperChar               string `yaml:"fret_wire_upper_char"`
	FretWireLowerChar               string `yaml:"fret_wire_lower_char"`
	StringChar                      string `yaml:"string_char"`
	StringOverFretWireChar          string `yaml:"string_over_fret_char"`
	StringOverBoarderChar           string `yaml:"string_over_boarder_char"`
	VibratingStringChar             string `yaml:"vibrating_string_char"`
	VibratingStringOverFretWireChar string `yaml:"vibrating_string_over_fret_char"`
	VibratingStringOverBoarderChar  string `yaml:"vibrating_string_over_boarder_char"`
	InlayChar                       string `yaml:"inlay_char"`
	PressedFretSignChar             string `yaml:"pressed_fret_sign_char"`
	NotPluckedStringChar            string `yaml:"not_plucked_string_char"`
	PluckedStringSignChar           string `yaml:"plucked_string_sign_char"`
}

type theme struct {
	// Color properties
	TitleFgColor           tcell.Color
	TitleBgColor           tcell.Color
	FretboardBorderColor   tcell.Color
	FretboardBgColor       tcell.Color
	NutBorderColor         tcell.Color
	NutBgColor             tcell.Color
	FretWireColor          tcell.Color
	StringColor            tcell.Color
	InlayColor             tcell.Color
	PressedFretSignColor   tcell.Color
	PluckedStringSignColor tcell.Color
	BaseNoteNameFgColor    tcell.Color
	BaseNoteNameBgColor    tcell.Color

	InlayShape InlayShape

	// Characters for display
	FretboardVBorderChar            rune
	FretboardHBorderChar            rune
	NutVBorderChar                  rune
	NutHBorderChar                  rune
	NutULCornerChar                 rune
	NutLLCornerChar                 rune
	NutURCornerChar                 rune
	NutLRCornerChar                 rune
	FretWireChar                    rune
	FretWireUpperChar               rune
	FretWireLowerChar               rune
	StringChar                      rune
	StringOverFretWireChar          rune
	StringOverBoarderChar           rune
	VibratingStringChar             rune
	VibratingStringOverFretWireChar rune
	VibratingStringOverBoarderChar  rune
	InlayChar                       rune
	PressedFretSignChar             rune
	NotPluckedStringChar            rune
	PluckedStringSignChar           rune
}

func loadRawTheme(name string) error {
	themeFilePath := filepath.Join(C.BaseDir, "themes", name+".yaml")
	tViper.SetConfigFile(themeFilePath)
	tViper.SetConfigType("yaml")

	if err := tViper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading theme '%s': %w", name, err)
	}
	if err := tViper.Unmarshal(&RawTheme, func(c *mapstructure.DecoderConfig) { c.TagName = "yaml" }); err != nil {
		return fmt.Errorf("failed to unmarshal theme '%s': %w", name, err)
	}
	return nil
}

func LoadTheme(name string) error {
	if err := loadRawTheme(name); err != nil {
		return err
	}

	Theme = theme{
		TitleFgColor:           tcell.GetColor(RawTheme.TitleFgColor),
		TitleBgColor:           tcell.GetColor(RawTheme.TitleBgColor),
		FretboardBorderColor:   tcell.GetColor(RawTheme.FretboardBorderColor),
		FretboardBgColor:       tcell.GetColor(RawTheme.FretboardBgColor),
		NutBorderColor:         tcell.GetColor(RawTheme.NutBorderColor),
		NutBgColor:             tcell.GetColor(RawTheme.NutBgColor),
		FretWireColor:          tcell.GetColor(RawTheme.FretWireColor),
		StringColor:            tcell.GetColor(RawTheme.StringColor),
		InlayColor:             tcell.GetColor(RawTheme.InlayColor),
		PressedFretSignColor:   tcell.GetColor(RawTheme.PressedFretSignColor),
		PluckedStringSignColor: tcell.GetColor(RawTheme.PluckedStringSignColor),
		BaseNoteNameFgColor:    tcell.GetColor(RawTheme.BaseNoteNameFgColor),
		BaseNoteNameBgColor:    tcell.GetColor(RawTheme.BaseNoteNameBgColor),

		InlayShape: getInlayShape(RawTheme.InlayShape),

		FretboardVBorderChar:            stringToRune(RawTheme.FretboardVBorderChar),
		FretboardHBorderChar:            stringToRune(RawTheme.FretboardHBorderChar),
		NutVBorderChar:                  stringToRune(RawTheme.NutVBorderChar),
		NutHBorderChar:                  stringToRune(RawTheme.NutHBorderChar),
		NutULCornerChar:                 stringToRune(RawTheme.NutULCornerChar),
		NutLLCornerChar:                 stringToRune(RawTheme.NutLLCornerChar),
		NutURCornerChar:                 stringToRune(RawTheme.NutURCornerChar),
		NutLRCornerChar:                 stringToRune(RawTheme.NutLRCornerChar),
		FretWireChar:                    stringToRune(RawTheme.FretWireChar),
		FretWireUpperChar:               stringToRune(RawTheme.FretWireUpperChar),
		FretWireLowerChar:               stringToRune(RawTheme.FretWireLowerChar),
		StringChar:                      stringToRune(RawTheme.StringChar),
		StringOverFretWireChar:          stringToRune(RawTheme.StringOverFretWireChar),
		StringOverBoarderChar:           stringToRune(RawTheme.StringOverBoarderChar),
		VibratingStringChar:             stringToRune(RawTheme.VibratingStringChar),
		VibratingStringOverFretWireChar: stringToRune(RawTheme.VibratingStringOverFretWireChar),
		VibratingStringOverBoarderChar:  stringToRune(RawTheme.VibratingStringOverBoarderChar),
		InlayChar:                       stringToRune(RawTheme.InlayChar),
		PressedFretSignChar:             stringToRune(RawTheme.PressedFretSignChar),
		NotPluckedStringChar:            stringToRune(RawTheme.NotPluckedStringChar),
		PluckedStringSignChar:           stringToRune(RawTheme.PluckedStringSignChar),
	}

	return nil
}

func stringToRune(s string) rune {
	runes := []rune(s)
	if len(runes) > 0 {
		return runes[0]
	}
	return ' '
}

func getInlayShape(shape string) InlayShape {
	shape = strings.TrimSpace(shape)
	shape = strings.ToLower(shape)
	switch shape {
	case "none", "":
		return NoneInlayShape
	case "dot", "dots":
		return DotInlayShape
	case "block", "blocks":
		return BlockInlayShape
	default:
		return DotInlayShape
	}
}
