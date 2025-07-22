package audio

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Golevka2001/bassit/bass"
	"github.com/Golevka2001/bassit/config"

	"github.com/ebitengine/oto/v3"
	"github.com/stretchr/testify/assert"
)

func TestBassNotePlayer(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T)
	}{
		{
			name: "BassNotePlayer struct creation",
			validate: func(t *testing.T) {
				player := BassNotePlayerGroup{}

				assert.Equal(t, config.PluckTypeCount, len(player.players))

				for i := range player.players {
					assert.Nil(t, player.players[i])
				}
			},
		},
		{
			name: "BassNotePlayer players array type",
			validate: func(t *testing.T) {
				player := BassNotePlayerGroup{}

				assert.IsType(t, [config.PluckTypeCount]*oto.Player{}, player.players)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t)
		})
	}
}

func TestAudioManager(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T)
	}{
		{
			name: "AudioManager struct fields",
			validate: func(t *testing.T) {
				am := &AudioManager{
					otoCtx:       nil,
					playerGroups: make(map[bass.FretboardPosition]BassNotePlayerGroup),
				}

				assert.NotNil(t, am.playerGroups)
				assert.IsType(t, map[bass.FretboardPosition]BassNotePlayerGroup{}, am.playerGroups)
			},
		},
		{
			name: "AudioManager players map initialization",
			validate: func(t *testing.T) {
				am := &AudioManager{
					playerGroups: make(map[bass.FretboardPosition]BassNotePlayerGroup),
				}

				assert.Equal(t, 0, len(am.playerGroups))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t)
		})
	}
}

func TestNewAudioManager(t *testing.T) {
	originalSampleRate := config.SampleRate
	defer func() { config.SampleRate = originalSampleRate }()

	tests := []struct {
		name     string
		setup    func()
		expected error
		validate func(t *testing.T, am *AudioManager, err error)
	}{
		{
			name: "successful AudioManager creation",
			setup: func() {
				config.SampleRate = 44100
			},
			expected: nil,
			validate: func(t *testing.T, am *AudioManager, err error) {
				if err != nil {
					t.Skip("Skipping test due to audio system unavailability")
					return
				}

				assert.NotNil(t, am)
				assert.NotNil(t, am.otoCtx)
				assert.NotNil(t, am.playerGroups)
				assert.Equal(t, 0, len(am.playerGroups))
			},
		},
		{
			name: "AudioManager creation with different sample rate",
			setup: func() {
				config.SampleRate = 48000
			},
			expected: nil,
			validate: func(t *testing.T, am *AudioManager, err error) {
				if err != nil {
					t.Skip("Skipping test due to audio system unavailability")
					return
				}

				assert.NotNil(t, am)
				assert.NotNil(t, am.otoCtx)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			am, err := NewAudioManager()

			if tt.expected != nil {
				assert.Error(t, err)
			} else {
				if err != nil {
					t.Skipf("Skipping test due to audio system error: %v", err)
					return
				}
				assert.NoError(t, err)
			}

			if tt.validate != nil {
				tt.validate(t, am, err)
			}
		})
	}
}

func TestAudioManagerGetPlayer(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *AudioManager
		pos       bass.FretboardPosition
		pluckType bass.PluckType
		expected  bool
	}{
		{
			name: "get player for non-existent position",
			setup: func() *AudioManager {
				return &AudioManager{
					playerGroups: make(map[bass.FretboardPosition]BassNotePlayerGroup),
				}
			},
			pos:       bass.FretboardPosition{StringIdx: 0, FretIdx: 0},
			pluckType: bass.PluckTypeNormal1,
			expected:  false,
		},
		{
			name: "get player for existing position with nil player",
			setup: func() *AudioManager {
				am := &AudioManager{
					playerGroups: make(map[bass.FretboardPosition]BassNotePlayerGroup),
				}
				pos := bass.FretboardPosition{StringIdx: 0, FretIdx: 0}
				am.playerGroups[pos] = BassNotePlayerGroup{}
				return am
			},
			pos:       bass.FretboardPosition{StringIdx: 0, FretIdx: 0},
			pluckType: bass.PluckTypeNormal1,
			expected:  false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := tt.setup()

			player := am.getPlayer(tt.pos, tt.pluckType)

			if tt.expected {
				assert.NotNil(t, player)
			} else {
				assert.Nil(t, player)
			}
		})
	}
}

func TestAudioManagerStopBassNote(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *AudioManager
		pos      bass.FretboardPosition
		validate func(t *testing.T, am *AudioManager)
	}{
		{
			name: "stop bass note for non-existent position",
			setup: func() *AudioManager {
				return &AudioManager{
					playerGroups: make(map[bass.FretboardPosition]BassNotePlayerGroup),
				}
			},
			pos: bass.FretboardPosition{StringIdx: 0, FretIdx: 0},
			validate: func(t *testing.T, am *AudioManager) {
				assert.NotNil(t, am)
			},
		},
		{
			name: "stop bass note for existing position with nil players",
			setup: func() *AudioManager {
				am := &AudioManager{
					playerGroups: make(map[bass.FretboardPosition]BassNotePlayerGroup),
				}
				pos := bass.FretboardPosition{StringIdx: 0, FretIdx: 0}
				am.playerGroups[pos] = BassNotePlayerGroup{}
				return am
			},
			pos: bass.FretboardPosition{StringIdx: 0, FretIdx: 0},
			validate: func(t *testing.T, am *AudioManager) {
				assert.NotNil(t, am)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := tt.setup()

			assert.NotPanics(t, func() {
				am.StopBassNote(tt.pos)
			})

			if tt.validate != nil {
				tt.validate(t, am)
			}
		})
	}
}

func TestAudioManagerLoadSoundpackSamples(t *testing.T) {
	originalSoundpackDir := config.SoundpackDir
	originalSoundpackName := config.SoundpackName
	originalDisplayedFretCount := config.DisplayedFretCount
	defer func() {
		config.SoundpackDir = originalSoundpackDir
		config.SoundpackName = originalSoundpackName
		config.DisplayedFretCount = originalDisplayedFretCount
	}()

	tests := []struct {
		name     string
		setup    func() (*AudioManager, *bass.BassModel, string)
		validate func(t *testing.T, am *AudioManager, tempDir string)
	}{
		{
			name: "load soundpack samples with non-existent directory",
			setup: func() (*AudioManager, *bass.BassModel, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_audio_test_*")

				config.SoundpackDir = func() string { return tempDir }
				config.SoundpackName = "nonexistent"
				config.DisplayedFretCount = 2
				am := &AudioManager{
					playerGroups: make(map[bass.FretboardPosition]BassNotePlayerGroup),
				}

				tuning := [config.StringCnt]string{"G2", "D2", "A1", "E1"}
				bassModel, _ := bass.NewBass(tuning)

				return am, bassModel, tempDir
			},
			validate: func(t *testing.T, am *AudioManager, tempDir string) {
				assert.Equal(t, 0, len(am.playerGroups))
			},
		},
		{
			name: "load soundpack samples with empty directory",
			setup: func() (*AudioManager, *bass.BassModel, string) {
				tempDir, _ := os.MkdirTemp("", "bassit_audio_test_*")
				soundpackDir := filepath.Join(tempDir, "test_soundpack")
				os.MkdirAll(soundpackDir, 0755)

				for t := range bass.PluckType(config.PluckTypeCount) {
					os.MkdirAll(filepath.Join(soundpackDir, t.String()), 0755)
				}

				config.SoundpackDir = func() string { return tempDir }
				config.SoundpackName = "test_soundpack"
				config.DisplayedFretCount = 1
				am := &AudioManager{
					playerGroups: make(map[bass.FretboardPosition]BassNotePlayerGroup),
				}

				tuning := [config.StringCnt]string{"G2", "D2", "A1", "E1"}
				bassModel, _ := bass.NewBass(tuning)

				return am, bassModel, tempDir
			},
			validate: func(t *testing.T, am *AudioManager, tempDir string) {
				assert.Equal(t, 0, len(am.playerGroups))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am, bassModel, tempDir := tt.setup()
			defer os.RemoveAll(tempDir)

			assert.NotPanics(t, func() {
				am.LoadSoundpackSamples(bassModel)
			})

			if tt.validate != nil {
				tt.validate(t, am, tempDir)
			}
		})
	}
}

func TestAudioManagerPlayBassNote(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *AudioManager
		pos       bass.FretboardPosition
		pluckType bass.PluckType
		validate  func(t *testing.T, am *AudioManager)
	}{
		{
			name: "play bass note for non-existent position",
			setup: func() *AudioManager {
				return &AudioManager{
					playerGroups: make(map[bass.FretboardPosition]BassNotePlayerGroup),
				}
			},
			pos:       bass.FretboardPosition{StringIdx: 0, FretIdx: 0},
			pluckType: bass.PluckTypeNormal1,
			validate: func(t *testing.T, am *AudioManager) {
				assert.NotNil(t, am)
			},
		},
		{
			name: "play bass note for existing position with nil player",
			setup: func() *AudioManager {
				am := &AudioManager{
					playerGroups: make(map[bass.FretboardPosition]BassNotePlayerGroup),
				}
				pos := bass.FretboardPosition{StringIdx: 0, FretIdx: 0}
				am.playerGroups[pos] = BassNotePlayerGroup{}
				return am
			},
			pos:       bass.FretboardPosition{StringIdx: 0, FretIdx: 0},
			pluckType: bass.PluckTypeNormal1,
			validate: func(t *testing.T, am *AudioManager) {
				assert.NotNil(t, am)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := tt.setup()

			assert.NotPanics(t, func() {
				am.PlayBassNote(tt.pos, tt.pluckType)
			})

			if tt.validate != nil {
				tt.validate(t, am)
			}
		})
	}
}
