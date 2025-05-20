package util

import (
	"math"
	"math/rand"
	"time"

	C "bassit/constant"
)

// GetVibDuration returns a random duration for string vibration
func GetVibDuration() time.Duration {
	duration := normalRand(C.StringVibDurnMean, C.StringVibDurStd)

	if duration < C.StringVibDurnMin {
		duration = C.StringVibDurnMin
	} else if duration > C.StringVibDurMax {
		duration = C.StringVibDurMax
	}

	return time.Duration(duration) * time.Millisecond
}

// normalRand uses the Box-Muller transform to generate a normally distributed random number
func normalRand(mean, stdDev float64) float64 {
	u1 := rand.Float64()
	u2 := rand.Float64()
	z := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
	return mean + z*stdDev
}
