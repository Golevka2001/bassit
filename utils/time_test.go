package utils

import (
	"math"
	"testing"
	"time"

	"github.com/Golevka2001/bassit/config"

	"github.com/stretchr/testify/assert"
)

func TestGetVibDuration(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T)
	}{
		{
			name: "duration within bounds",
			validate: func(t *testing.T) {
				for i := range 100 {
					_ = i
					duration := GetVibDuration()
					assert.GreaterOrEqual(t, duration, time.Duration(config.StringVibDurMin)*time.Millisecond)
					assert.LessOrEqual(t, duration, time.Duration(config.StringVibDurMax)*time.Millisecond)
				}
			},
		},
		{
			name: "duration type is correct",
			validate: func(t *testing.T) {
				duration := GetVibDuration()
				assert.IsType(t, time.Duration(0), duration)
			},
		},
		{
			name: "duration is positive",
			validate: func(t *testing.T) {
				for i := range 50 {
					_ = i
					duration := GetVibDuration()
					assert.Greater(t, duration, time.Duration(0))
				}
			},
		},
		{
			name: "duration distribution is reasonable",
			validate: func(t *testing.T) {
				samples := make([]time.Duration, 1000)
				for i := range samples {
					samples[i] = GetVibDuration()
				}

				var sum time.Duration
				for _, sample := range samples {
					sum += sample
				}
				meanDuration := sum / time.Duration(len(samples))

				expectedMean := time.Duration(config.StringVibDurMean) * time.Millisecond
				tolerance := time.Duration(config.StringVibDurStd*3) * time.Millisecond
				assert.InDelta(t, float64(expectedMean), float64(meanDuration), float64(tolerance))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t)
		})
	}
}

func TestNormalRand(t *testing.T) {
	tests := []struct {
		name     string
		mean     float64
		stdDev   float64
		validate func(t *testing.T, mean, stdDev float64)
	}{
		{
			name:   "standard normal distribution",
			mean:   0.0,
			stdDev: 1.0,
			validate: func(t *testing.T, mean, stdDev float64) {
				samples := make([]float64, 10000)
				for i := range samples {
					samples[i] = normalRand(mean, stdDev)
				}

				var sum float64
				for _, sample := range samples {
					sum += sample
				}
				sampleMean := sum / float64(len(samples))

				assert.InDelta(t, mean, sampleMean, 0.1)
			},
		},
		{
			name:   "custom distribution",
			mean:   100.0,
			stdDev: 15.0,
			validate: func(t *testing.T, mean, stdDev float64) {
				samples := make([]float64, 10000)
				for i := range samples {
					samples[i] = normalRand(mean, stdDev)
				}

				var sum float64
				for _, sample := range samples {
					sum += sample
				}
				sampleMean := sum / float64(len(samples))

				var sumSquaredDiff float64
				for _, sample := range samples {
					diff := sample - sampleMean
					sumSquaredDiff += diff * diff
				}
				sampleStdDev := math.Sqrt(sumSquaredDiff / float64(len(samples)-1))

				assert.InDelta(t, mean, sampleMean, 1.0)
				assert.InDelta(t, stdDev, sampleStdDev, 2.0)
			},
		},
		{
			name:   "zero standard deviation",
			mean:   50.0,
			stdDev: 0.0,
			validate: func(t *testing.T, mean, stdDev float64) {
				for i := range 100 {
					_ = i
					value := normalRand(mean, stdDev)
					assert.Equal(t, mean, value)
				}
			},
		},
		{
			name:   "negative mean",
			mean:   -25.0,
			stdDev: 5.0,
			validate: func(t *testing.T, mean, stdDev float64) {
				samples := make([]float64, 1000)
				for i := range samples {
					samples[i] = normalRand(mean, stdDev)
				}

				var sum float64
				for _, sample := range samples {
					sum += sample
				}
				sampleMean := sum / float64(len(samples))

				assert.InDelta(t, mean, sampleMean, 1.0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t, tt.mean, tt.stdDev)
		})
	}
}
