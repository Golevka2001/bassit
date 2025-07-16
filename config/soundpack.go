package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-music-theory/music-theory/note"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

var (
	RubberbandCommand = "rubberband-r3"
)

var (
	SampleRate int
	BasePitch  string
	Normal     [2]SampleItem
	Slap       [2]SampleItem
	Mute       [2]SampleItem
)

type SampleItem struct {
	Name string `mapstructure:"name"`
	File string `mapstructure:"file"`
}

type SoundpackInfo struct {
	// Name of the sound pack
	Name string `yaml:"name"`

	// Description of the sound pack
	Description string `yaml:"description"`

	// Author of the sound pack
	Author string `yaml:"author"`

	// SampleRate is the sample rate for all WAV files in this sound pack
	SampleRate int `yaml:"sample_rate"`

	// BasePitch indicates the pitch of all samples in this sound pack
	BasePitch string `yaml:"base_pitch"`

	Normal [2]SampleItem `mapstructure:"normal"`
	Slap   [2]SampleItem `mapstructure:"slap"`
	Mute   [2]SampleItem `mapstructure:"mute"`
}

func LoadSoundpackInfo(name string) (*SoundpackInfo, error) {
	v := viper.New()
	v.SetConfigFile(filepath.Join(SoundpackDir(), name, "info.yaml"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading soundpack info '%s': %w", name, err)
	}

	// Load the soundpack info
	var info SoundpackInfo
	if err := v.Unmarshal(&info, func(c *mapstructure.DecoderConfig) {
		c.TagName = "yaml"
	}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal soundpack info '%s': %w", name, err)
	}

	// Validate the soundpack info
	if err := validateSoundpackInfo(&info, name); err != nil {
		return nil, fmt.Errorf("invalid soundpack info '%s': %w", name, err)
	}

	SampleRate = info.SampleRate
	BasePitch = info.BasePitch
	Normal = info.Normal
	Slap = info.Slap
	Mute = info.Mute

	return &info, nil
}

func validateSoundpackInfo(info *SoundpackInfo, name string) error {
	// `Name`
	if info.Name != name {
		return fmt.Errorf("soundpack name mismatch: %s != %s", info.Name, name)
	}

	// `BasePitch`
	n := note.Named(info.BasePitch)
	if n.Class == note.Nil || n.Octave < 0 {
		return fmt.Errorf("invalid base pitch: %s", info.BasePitch)
	}

	// `Normal`, `Slap`, `Mute`
	for _, items := range [][]SampleItem{info.Normal[:], info.Slap[:], info.Mute[:]} {
		for _, item := range items {
			// `Name`
			if err := validateName(item.Name); err != nil {
				return err
			}
			// `File`
			filePath := filepath.Join(SoundpackDir(), name, item.File)
			if _, err := os.Stat(filePath); err != nil {
				return fmt.Errorf("file not found: %s", filePath)
			}
		}
	}

	return nil
}

func validateName(name string) error {
	if strings.ContainsAny(name, "/\\:*?\"<>|") {
		return fmt.Errorf("invalid name: %s, \"/\\:*?\"<>|\" are not allowed", name)
	}

	return nil
}
