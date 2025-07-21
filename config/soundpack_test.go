package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSampleItem(t *testing.T) {
	tests := []struct {
		name     string
		item     SampleItem
		validate func(t *testing.T, item SampleItem)
	}{
		{
			name: "valid sample item",
			item: SampleItem{
				Name: "normal1",
				File: "normal/sample1.wav",
			},
			validate: func(t *testing.T, item SampleItem) {
				assert.Equal(t, "normal1", item.Name)
				assert.Equal(t, "normal/sample1.wav", item.File)
			},
		},
		{
			name: "empty sample item",
			item: SampleItem{
				Name: "",
				File: "",
			},
			validate: func(t *testing.T, item SampleItem) {
				assert.Equal(t, "", item.Name)
				assert.Equal(t, "", item.File)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t, tt.item)
		})
	}
}

func TestSoundpackInfo(t *testing.T) {
	tests := []struct {
		name     string
		info     SoundpackInfo
		validate func(t *testing.T, info SoundpackInfo)
	}{
		{
			name: "valid soundpack info",
			info: SoundpackInfo{
				Name:        "test_soundpack",
				Description: "A test soundpack",
				Author:      "Test Author",
				SampleRate:  44100,
				BasePitch:   "E1",
				Normal: [2]SampleItem{
					{Name: "normal1", File: "normal/sample1.wav"},
					{Name: "normal2", File: "normal/sample2.wav"},
				},
				Slap: [2]SampleItem{
					{Name: "slap1", File: "slap/sample1.wav"},
					{Name: "slap2", File: "slap/sample2.wav"},
				},
				Mute: [2]SampleItem{
					{Name: "mute1", File: "mute/sample1.wav"},
					{Name: "mute2", File: "mute/sample2.wav"},
				},
			},
			validate: func(t *testing.T, info SoundpackInfo) {
				assert.Equal(t, "test_soundpack", info.Name)
				assert.Equal(t, "A test soundpack", info.Description)
				assert.Equal(t, "Test Author", info.Author)
				assert.Equal(t, 44100, info.SampleRate)
				assert.Equal(t, "E1", info.BasePitch)
				assert.Equal(t, "normal1", info.Normal[0].Name)
				assert.Equal(t, "normal/sample1.wav", info.Normal[0].File)
				assert.Equal(t, "slap1", info.Slap[0].Name)
				assert.Equal(t, "slap/sample1.wav", info.Slap[0].File)
				assert.Equal(t, "mute1", info.Mute[0].Name)
				assert.Equal(t, "mute/sample1.wav", info.Mute[0].File)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t, tt.info)
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected error
	}{
		{
			name:     "valid name",
			input:    "normal1",
			expected: nil,
		},
		{
			name:     "valid name with underscore",
			input:    "sample_1",
			expected: nil,
		},
		{
			name:     "valid name with dash",
			input:    "sample-1",
			expected: nil,
		},
		{
			name:     "invalid name with slash",
			input:    "normal/1",
			expected: assert.AnError,
		},
		{
			name:     "invalid name with backslash",
			input:    "normal\\1",
			expected: assert.AnError,
		},
		{
			name:     "invalid name with colon",
			input:    "normal:1",
			expected: assert.AnError,
		},
		{
			name:     "invalid name with asterisk",
			input:    "normal*1",
			expected: assert.AnError,
		},
		{
			name:     "invalid name with question mark",
			input:    "normal?1",
			expected: assert.AnError,
		},
		{
			name:     "invalid name with quotes",
			input:    "normal\"1",
			expected: assert.AnError,
		},
		{
			name:     "invalid name with angle brackets",
			input:    "normal<1>",
			expected: assert.AnError,
		},
		{
			name:     "invalid name with pipe",
			input:    "normal|1",
			expected: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateName(tt.input)

			if tt.expected != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid name")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSoundpackInfo(t *testing.T) {
	originalBaseDir := BaseDir
	defer func() { BaseDir = originalBaseDir }()

	tests := []struct {
		name     string
		setup    func() (SoundpackInfo, string, string)
		expected error
		validate func(t *testing.T, err error)
	}{
		{
			name: "valid soundpack info",
			setup: func() (SoundpackInfo, string, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_soundpack_test_*")

				soundpackDir := filepath.Join(tempDir, "soundpacks", "test_pack")
				normalDir := filepath.Join(soundpackDir, "normal")
				slapDir := filepath.Join(soundpackDir, "slap")
				muteDir := filepath.Join(soundpackDir, "mute")
				os.MkdirAll(normalDir, 0755)
				os.MkdirAll(slapDir, 0755)
				os.MkdirAll(muteDir, 0755)

				os.WriteFile(filepath.Join(normalDir, "sample1.wav"), []byte("fake wav"), 0644)
				os.WriteFile(filepath.Join(normalDir, "sample2.wav"), []byte("fake wav"), 0644)
				os.WriteFile(filepath.Join(slapDir, "sample1.wav"), []byte("fake wav"), 0644)
				os.WriteFile(filepath.Join(slapDir, "sample2.wav"), []byte("fake wav"), 0644)
				os.WriteFile(filepath.Join(muteDir, "sample1.wav"), []byte("fake wav"), 0644)
				os.WriteFile(filepath.Join(muteDir, "sample2.wav"), []byte("fake wav"), 0644)

				BaseDir = tempDir

				info := SoundpackInfo{
					Name:        "test_pack",
					Description: "Test soundpack",
					Author:      "Test Author",
					SampleRate:  44100,
					BasePitch:   "E1",
					Normal: [2]SampleItem{
						{Name: "normal1", File: "normal/sample1.wav"},
						{Name: "normal2", File: "normal/sample2.wav"},
					},
					Slap: [2]SampleItem{
						{Name: "slap1", File: "slap/sample1.wav"},
						{Name: "slap2", File: "slap/sample2.wav"},
					},
					Mute: [2]SampleItem{
						{Name: "mute1", File: "mute/sample1.wav"},
						{Name: "mute2", File: "mute/sample2.wav"},
					},
				}

				return info, "test_pack", tempDir
			},
			expected: nil,
			validate: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "name mismatch",
			setup: func() (SoundpackInfo, string, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_soundpack_test_*")
				BaseDir = tempDir

				info := SoundpackInfo{
					Name:       "wrong_name",
					SampleRate: 44100,
					BasePitch:  "E1",
				}

				return info, "test_pack", tempDir
			},
			expected: assert.AnError,
			validate: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "soundpack name mismatch")
			},
		},
		{
			name: "invalid base pitch",
			setup: func() (SoundpackInfo, string, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_soundpack_test_*")
				BaseDir = tempDir

				info := SoundpackInfo{
					Name:       "test_pack",
					SampleRate: 44100,
					BasePitch:  "X1"}

				return info, "test_pack", tempDir
			},
			expected: assert.AnError,
			validate: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid base pitch")
			},
		},
		{
			name: "invalid sample name",
			setup: func() (SoundpackInfo, string, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_soundpack_test_*")
				BaseDir = tempDir

				info := SoundpackInfo{
					Name:       "test_pack",
					SampleRate: 44100,
					BasePitch:  "E1",
					Normal: [2]SampleItem{
						{Name: "normal/1", File: "normal/sample1.wav"}, {Name: "normal2", File: "normal/sample2.wav"},
					},
				}

				return info, "test_pack", tempDir
			},
			expected: assert.AnError,
			validate: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid name")
			},
		},
		{
			name: "missing sample file",
			setup: func() (SoundpackInfo, string, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_soundpack_test_*")

				soundpackDir := filepath.Join(tempDir, "soundpacks", "test_pack")
				os.MkdirAll(soundpackDir, 0755)

				BaseDir = tempDir

				info := SoundpackInfo{
					Name:       "test_pack",
					SampleRate: 44100,
					BasePitch:  "E1",
					Normal: [2]SampleItem{
						{Name: "normal1", File: "normal/sample1.wav"}, {Name: "normal2", File: "normal/sample2.wav"},
					},
				}

				return info, "test_pack", tempDir
			},
			expected: assert.AnError,
			validate: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "file not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, soundpackName, tempDir := tt.setup()
			defer os.RemoveAll(tempDir)

			err := validateSoundpackInfo(&info, soundpackName)

			if tt.validate != nil {
				tt.validate(t, err)
			}
		})
	}
}

func TestLoadSoundpackInfo(t *testing.T) {
	originalBaseDir := BaseDir
	originalSampleRate := SampleRate
	originalBasePitch := BasePitch
	originalNormal := Normal
	originalSlap := Slap
	originalMute := Mute
	defer func() {
		BaseDir = originalBaseDir
		SampleRate = originalSampleRate
		BasePitch = originalBasePitch
		Normal = originalNormal
		Slap = originalSlap
		Mute = originalMute
	}()

	tests := []struct {
		name     string
		setup    func() (string, string)
		expected error
		validate func(t *testing.T, info *SoundpackInfo, err error, tempDir string)
	}{
		{
			name: "load valid soundpack info",
			setup: func() (string, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_soundpack_test_*")

				soundpackDir := filepath.Join(tempDir, "soundpacks", "test_pack")
				normalDir := filepath.Join(soundpackDir, "normal")
				slapDir := filepath.Join(soundpackDir, "slap")
				muteDir := filepath.Join(soundpackDir, "mute")
				os.MkdirAll(normalDir, 0755)
				os.MkdirAll(slapDir, 0755)
				os.MkdirAll(muteDir, 0755)

				os.WriteFile(filepath.Join(normalDir, "sample1.wav"), []byte("fake wav"), 0644)
				os.WriteFile(filepath.Join(normalDir, "sample2.wav"), []byte("fake wav"), 0644)
				os.WriteFile(filepath.Join(slapDir, "sample1.wav"), []byte("fake wav"), 0644)
				os.WriteFile(filepath.Join(slapDir, "sample2.wav"), []byte("fake wav"), 0644)
				os.WriteFile(filepath.Join(muteDir, "sample1.wav"), []byte("fake wav"), 0644)
				os.WriteFile(filepath.Join(muteDir, "sample2.wav"), []byte("fake wav"), 0644)

				infoContent := `
name: "test_pack"
description: "Test soundpack"
author: "Test Author"
sample_rate: 48000
base_pitch: "A1"
normal:
  - name: "normal1"
    file: "normal/sample1.wav"
  - name: "normal2"
    file: "normal/sample2.wav"
slap:
  - name: "slap1"
    file: "slap/sample1.wav"
  - name: "slap2"
    file: "slap/sample2.wav"
mute:
  - name: "mute1"
    file: "mute/sample1.wav"
  - name: "mute2"
    file: "mute/sample2.wav"
`
				infoFile := filepath.Join(soundpackDir, "info.yaml")
				os.WriteFile(infoFile, []byte(infoContent), 0644)

				BaseDir = tempDir
				return "test_pack", tempDir
			},
			expected: nil,
			validate: func(t *testing.T, info *SoundpackInfo, err error, tempDir string) {
				assert.NoError(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, "test_pack", info.Name)
				assert.Equal(t, "Test soundpack", info.Description)
				assert.Equal(t, "Test Author", info.Author)
				assert.Equal(t, 48000, info.SampleRate)
				assert.Equal(t, "A1", info.BasePitch)

				assert.Equal(t, 48000, SampleRate)
				assert.Equal(t, "A1", BasePitch)
				assert.Equal(t, "normal1", Normal[0].Name)
				assert.Equal(t, "slap1", Slap[0].Name)
				assert.Equal(t, "mute1", Mute[0].Name)
			},
		},
		{
			name: "soundpack info file not found",
			setup: func() (string, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_soundpack_test_*")

				soundpackDir := filepath.Join(tempDir, "soundpacks", "test_pack")
				os.MkdirAll(soundpackDir, 0755)

				BaseDir = tempDir
				return "test_pack", tempDir
			},
			expected: assert.AnError,
			validate: func(t *testing.T, info *SoundpackInfo, err error, tempDir string) {
				assert.Error(t, err)
				assert.Nil(t, info)
				assert.Contains(t, err.Error(), "error reading soundpack info")
			},
		},
		{
			name: "invalid YAML syntax",
			setup: func() (string, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_soundpack_test_*")

				soundpackDir := filepath.Join(tempDir, "soundpacks", "test_pack")
				os.MkdirAll(soundpackDir, 0755)

				infoContent := `
name: "test_pack"
invalid yaml syntax here
sample_rate: 44100
`
				infoFile := filepath.Join(soundpackDir, "info.yaml")
				os.WriteFile(infoFile, []byte(infoContent), 0644)

				BaseDir = tempDir
				return "test_pack", tempDir
			},
			expected: assert.AnError,
			validate: func(t *testing.T, info *SoundpackInfo, err error, tempDir string) {
				assert.Error(t, err)
				assert.Nil(t, info)
				assert.Contains(t, err.Error(), "error reading soundpack info")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			soundpackName, tempDir := tt.setup()
			defer os.RemoveAll(tempDir)

			info, err := LoadSoundpackInfo(soundpackName)

			if tt.validate != nil {
				tt.validate(t, info, err, tempDir)
			}
		})
	}
}
