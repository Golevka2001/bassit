package bass

import (
	"testing"

	"github.com/Golevka2001/bassit/config"
	
	"github.com/stretchr/testify/assert"
)

func TestFretboardPosition(t *testing.T) {
	tests := []struct {
		name     string
		pos      FretboardPosition
		validate func(t *testing.T, pos FretboardPosition)
	}{
		{
			name: "valid fretboard position - open string",
			pos:  FretboardPosition{StringIdx: 0, FretIdx: 0},
			validate: func(t *testing.T, pos FretboardPosition) {
				assert.Equal(t, 0, pos.StringIdx)
				assert.Equal(t, 0, pos.FretIdx)
			},
		},
		{
			name: "valid fretboard position - fretted note",
			pos:  FretboardPosition{StringIdx: 2, FretIdx: 5},
			validate: func(t *testing.T, pos FretboardPosition) {
				assert.Equal(t, 2, pos.StringIdx)
				assert.Equal(t, 5, pos.FretIdx)
			},
		},
		{
			name: "fretboard position with maximum values",
			pos:  FretboardPosition{StringIdx: config.StringCnt - 1, FretIdx: config.MaxFretCnt},
			validate: func(t *testing.T, pos FretboardPosition) {
				assert.Equal(t, config.StringCnt-1, pos.StringIdx)
				assert.Equal(t, config.MaxFretCnt, pos.FretIdx)
			},
		},
		{
			name: "fretboard position as map key",
			pos:  FretboardPosition{StringIdx: 1, FretIdx: 3},
			validate: func(t *testing.T, pos FretboardPosition) {
				testMap := make(map[FretboardPosition]string)
				testMap[pos] = "test_value"

				assert.Equal(t, "test_value", testMap[pos])
				assert.Equal(t, 1, len(testMap))

				pos2 := FretboardPosition{StringIdx: 1, FretIdx: 4}
				testMap[pos2] = "another_value"
				assert.Equal(t, 2, len(testMap))
				assert.Equal(t, "test_value", testMap[pos])
				assert.Equal(t, "another_value", testMap[pos2])
			},
		},
		{
			name: "fretboard position equality",
			pos:  FretboardPosition{StringIdx: 2, FretIdx: 7},
			validate: func(t *testing.T, pos FretboardPosition) {
				pos2 := FretboardPosition{StringIdx: 2, FretIdx: 7}
				pos3 := FretboardPosition{StringIdx: 2, FretIdx: 8}

				assert.Equal(t, pos, pos2)
				assert.NotEqual(t, pos, pos3)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t, tt.pos)
		})
	}
}

func TestFretboardPositionBounds(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T)
	}{
		{
			name: "string index bounds",
			validate: func(t *testing.T) {
				validMin := FretboardPosition{StringIdx: 0, FretIdx: 0}
				validMax := FretboardPosition{StringIdx: config.StringCnt - 1, FretIdx: 0}

				assert.GreaterOrEqual(t, validMin.StringIdx, 0)
				assert.Less(t, validMax.StringIdx, config.StringCnt)
			},
		},
		{
			name: "fret index bounds",
			validate: func(t *testing.T) {
				validMin := FretboardPosition{StringIdx: 0, FretIdx: 0}
				validMax := FretboardPosition{StringIdx: 0, FretIdx: config.MaxFretCnt}

				assert.GreaterOrEqual(t, validMin.FretIdx, 0)
				assert.LessOrEqual(t, validMax.FretIdx, config.MaxFretCnt)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t)
		})
	}
}

func TestPluckType(t *testing.T) {
	tests := []struct {
		name      string
		pluckType PluckType
		expected  string
	}{
		{
			name:      "PluckTypeNormal1",
			pluckType: PluckTypeNormal1,
			expected:  "normal1",
		},
		{
			name:      "PluckTypeNormal2",
			pluckType: PluckTypeNormal2,
			expected:  "normal2",
		},
		{
			name:      "PluckTypeSlap1",
			pluckType: PluckTypeSlap1,
			expected:  "slap1",
		},
		{
			name:      "PluckTypeSlap2",
			pluckType: PluckTypeSlap2,
			expected:  "slap2",
		},
		{
			name:      "PluckTypeMute1",
			pluckType: PluckTypeMute1,
			expected:  "mute1",
		},
		{
			name:      "PluckTypeMute2",
			pluckType: PluckTypeMute2,
			expected:  "mute2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.pluckType.String()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestPluckTypeValues(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T)
	}{
		{
			name: "PluckType values are within PluckTypeCount",
			validate: func(t *testing.T) {
				assert.Less(t, int(PluckTypeNormal1), config.PluckTypeCount)
				assert.Less(t, int(PluckTypeNormal2), config.PluckTypeCount)
				assert.Less(t, int(PluckTypeSlap1), config.PluckTypeCount)
				assert.Less(t, int(PluckTypeSlap2), config.PluckTypeCount)
				assert.Less(t, int(PluckTypeMute1), config.PluckTypeCount)
				assert.Less(t, int(PluckTypeMute2), config.PluckTypeCount)
			},
		},
		{
			name: "PluckType values are sequential",
			validate: func(t *testing.T) {
				assert.Equal(t, 0, int(PluckTypeNormal1))
				assert.Equal(t, 1, int(PluckTypeNormal2))
				assert.Equal(t, 2, int(PluckTypeSlap1))
				assert.Equal(t, 3, int(PluckTypeSlap2))
				assert.Equal(t, 4, int(PluckTypeMute1))
				assert.Equal(t, 5, int(PluckTypeMute2))
			},
		},
		{
			name: "PluckType count matches config",
			validate: func(t *testing.T) {
				assert.Equal(t, 6, config.PluckTypeCount)

				for i := 0; i < config.PluckTypeCount; i++ {
					pluckType := PluckType(i)
					stringValue := pluckType.String()
					assert.NotEmpty(t, stringValue)
				}
			},
		},
		{
			name: "PluckType string representations are unique",
			validate: func(t *testing.T) {
				stringValues := make(map[string]bool)

				for i := 0; i < config.PluckTypeCount; i++ {
					pluckType := PluckType(i)
					stringValue := pluckType.String()

					assert.False(t, stringValues[stringValue], "Duplicate string value: %s", stringValue)
					stringValues[stringValue] = true
				}

				assert.Equal(t, config.PluckTypeCount, len(stringValues))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t)
		})
	}
}

func TestPluckTypeCategories(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T)
	}{
		{
			name: "normal pluck types",
			validate: func(t *testing.T) {
				normalTypes := []PluckType{PluckTypeNormal1, PluckTypeNormal2}

				for _, pt := range normalTypes {
					stringValue := pt.String()
					assert.Contains(t, stringValue, "normal")
				}
			},
		},
		{
			name: "slap pluck types",
			validate: func(t *testing.T) {
				slapTypes := []PluckType{PluckTypeSlap1, PluckTypeSlap2}

				for _, pt := range slapTypes {
					stringValue := pt.String()
					assert.Contains(t, stringValue, "slap")
				}
			},
		},
		{
			name: "mute pluck types",
			validate: func(t *testing.T) {
				muteTypes := []PluckType{PluckTypeMute1, PluckTypeMute2}

				for _, pt := range muteTypes {
					stringValue := pt.String()
					assert.Contains(t, stringValue, "mute")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t)
		})
	}
}
