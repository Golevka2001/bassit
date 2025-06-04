package util

import (
	"testing"
	"time"

	C "github.com/Golevka2001/bassit/constant"

	"github.com/stretchr/testify/assert"
)

func TestGetVibDuration(t *testing.T) {
	// Test that the function returns a duration within expected bounds
	for i := range 100 {
		_ = i
		duration := GetVibDuration()
		assert.GreaterOrEqual(t, duration, C.StringVibDurMin*time.Millisecond, "Duration should be at least %d ms", C.StringVibDurMin)
		assert.LessOrEqual(t, duration, C.StringVibDurMax*time.Millisecond, "Duration should be at most %d ms", C.StringVibDurMax)
	}
}

func TestNormalRand(t *testing.T) {
	mean := 100.0
	stdDev := 15.0

	// Generate many samples to test distribution properties
	samples := make([]float64, 10000)
	for i := range samples {
		samples[i] = normalRand(mean, stdDev)
	}

	// Calculate sample mean
	var sum float64
	for _, sample := range samples {
		sum += sample
	}
	sampleMean := sum / float64(len(samples))

	// Test that sample mean is close to expected mean
	assert.InDelta(t, mean, sampleMean, 1.0, "Sample mean should be close to expected mean")
}
