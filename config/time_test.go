package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeConstants(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{
			name:     "StringVibDurMean constant",
			input:    StringVibDurMean,
			expected: 850.0,
		},
		{
			name:     "StringVibDurStd constant",
			input:    StringVibDurStd,
			expected: 50.0,
		},
		{
			name:     "StringVibDurMin constant",
			input:    StringVibDurMin,
			expected: 700.0,
		},
		{
			name:     "StringVibDurMax constant",
			input:    StringVibDurMax,
			expected: 1000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input)
		})
	}
}

func TestTimeConstantTypes(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T)
	}{
		{
			name: "StringVibDurMean is float64",
			validate: func(t *testing.T) {
				assert.IsType(t, float64(0), StringVibDurMean)
			},
		},
		{
			name: "StringVibDurStd is float64",
			validate: func(t *testing.T) {
				assert.IsType(t, float64(0), StringVibDurStd)
			},
		},
		{
			name: "StringVibDurMin is float64",
			validate: func(t *testing.T) {
				assert.IsType(t, float64(0), StringVibDurMin)
			},
		},
		{
			name: "StringVibDurMax is float64",
			validate: func(t *testing.T) {
				assert.IsType(t, float64(0), StringVibDurMax)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t)
		})
	}
}

func TestTimeConstantRanges(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T)
	}{
		{
			name: "All time constants are positive",
			validate: func(t *testing.T) {
				assert.Greater(t, StringVibDurMean, 0.0)
				assert.Greater(t, StringVibDurStd, 0.0)
				assert.Greater(t, StringVibDurMin, 0.0)
				assert.Greater(t, StringVibDurMax, 0.0)
			},
		},
		{
			name: "StringVibDurMin is less than StringVibDurMax",
			validate: func(t *testing.T) {
				assert.Less(t, StringVibDurMin, StringVibDurMax)
			},
		},
		{
			name: "StringVibDurMean is between min and max",
			validate: func(t *testing.T) {
				assert.GreaterOrEqual(t, StringVibDurMean, StringVibDurMin)
				assert.LessOrEqual(t, StringVibDurMean, StringVibDurMax)
			},
		},
		{
			name: "StringVibDurStd is reasonable compared to mean",
			validate: func(t *testing.T) {
				assert.Less(t, StringVibDurStd, StringVibDurMean)
				assert.Greater(t, StringVibDurStd, StringVibDurMean*0.01)
			},
		},
		{
			name: "Time values are in reasonable millisecond range",
			validate: func(t *testing.T) {
				assert.GreaterOrEqual(t, StringVibDurMin, 100.0)
				assert.LessOrEqual(t, StringVibDurMax, 5000.0)
				assert.GreaterOrEqual(t, StringVibDurMean, 100.0)
				assert.LessOrEqual(t, StringVibDurMean, 5000.0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t)
		})
	}
}

func TestTimeConstantRelationships(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T)
	}{
		{
			name: "Mean should be approximately centered between min and max",
			validate: func(t *testing.T) {
				expectedMean := (StringVibDurMin + StringVibDurMax) / 2
				tolerance := (StringVibDurMax - StringVibDurMin) * 0.3
				assert.InDelta(t, expectedMean, StringVibDurMean, tolerance)
			},
		},
		{
			name: "Standard deviation should allow for reasonable distribution",
			validate: func(t *testing.T) {
				lowerBound := StringVibDurMean - 2*StringVibDurStd
				upperBound := StringVibDurMean + 2*StringVibDurStd

				assert.GreaterOrEqual(t, lowerBound, 0.0)

				assert.GreaterOrEqual(t, lowerBound, StringVibDurMin*0.8)
				assert.LessOrEqual(t, upperBound, StringVibDurMax*1.2)
			},
		},
		{
			name: "Range (max - min) should be reasonable",
			validate: func(t *testing.T) {
				vibrationRange := StringVibDurMax - StringVibDurMin

				assert.Greater(t, vibrationRange, 0.0)

				assert.GreaterOrEqual(t, vibrationRange, 100.0)
				assert.LessOrEqual(t, vibrationRange, 2000.0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t)
		})
	}
}
