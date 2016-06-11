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

type exponentialBackoff struct {
	multiplier  float64
	randFactor  float64
	minInterval time.Duration
	maxInterval time.Duration

	attempts    uint
	maxAttempts uint
}

func (b *exponentialBackoff) Reset() {
	b.attempts = 0
}

func (b *exponentialBackoff) NextInterval() time.Duration {
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

// Creates an exponential backoff interval generator using the default
// values for multipler, random factor, minimum, and maximum intervals.
func NewDefaultExponentialBackoff() Backoff {
	return NewExponentialBackoff(
		DefaultMultiplier,
		DefaultRandFactor,
		DefaultMinInterval,
		DefaultMaxInterval,
	)
}

// A backoff interval generator which returns exponentially increasing
// intervals for each unsuccessful retry. The base interval is given by
// the function (MinInterval * Multiplier ^ n) where n is the number of
// previous failed attempts in the current sequence. The value returned
// is given by min(MaxInterval, base +/- (base * RandFactor)). A random
// factor of zero will make the generator deterministic.
func NewExponentialBackoff(multiplier, randFactor float64, minInterval, maxInterval time.Duration) Backoff {
	// min * mult ^ n <     max
	//       mult ^ n <     max / min
	//              n < log(max / min) / log(mult)
	maxAttempts := math.Log(float64(maxInterval/minInterval)) / math.Log(multiplier)

	b := &exponentialBackoff{
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
