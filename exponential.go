package backoff

import (
	"math"
	"math/rand"
	"time"
)

const (
	DefaultMultiplier  = 1.25
	DefaultRandFactor  = 0.25
	DefaultMinInterval = 10 * time.Millisecond
	DefaultMaxInterval = 10 * time.Minute
)

type exponentialBackOff struct {
	multiplier  float64
	randFactor  float64
	minInterval time.Duration
	maxInterval time.Duration

	attempts    uint
	maxAttempts uint
}

func (b *exponentialBackOff) Reset() {
	b.attempts = 0
}

func (b *exponentialBackOff) NextInterval() time.Duration {
	if b.attempts >= b.maxAttempts {
		return b.maxInterval
	}

	n := float64(b.attempts)
	b.attempts += 1

	return time.Duration(randomNear(float64(b.minInterval)*math.Pow(b.multiplier, n), b.randFactor))
}

func randomNear(value, ratio float64) float64 {
	min := value - (value * ratio)
	max := value + (value * ratio)

	return min + (max-min+1)*rand.Float64()
}

func NewDefaultExponentialBackoff() BackOff {
	return NewExponentialBackOff(
		DefaultMultiplier,
		DefaultRandFactor,
		DefaultMinInterval,
		DefaultMaxInterval,
	)
}

// A back-off interval generator which returns exponentially increasing
// intervals for each unsuccessful retry. The base interval is given by
// the function (MinInterval * Multiplier ^ n) where n is the number of
// previous failed attempts in the current sequence. The value returned
// is given by min(MaxInterval, base +/- (base * RandFactor)). A random
// factor of zero will make the generator deterministic.
func NewExponentialBackOff(multiplier, randFactor float64, minInterval, maxInterval time.Duration) BackOff {
	// min * mult ^ n <     max
	//       mult ^ n <     max / min
	//              n < log(max / min) / log(mult)
	maxAttempts := math.Log(float64(maxInterval/minInterval)) / math.Log(multiplier)

	b := &exponentialBackOff{
		multiplier:  multiplier,
		randFactor:  randFactor,
		minInterval: minInterval,
		maxInterval: maxInterval,

		attempts:    0,
		maxAttempts: uint(maxAttempts),
	}

	b.Reset()
	return b
}
