package backoff

import (
	"math"
	"math/rand"
	"time"
)

// ExponentialConfigFunc is a function used to initialize a new exponential backoff.
type ExponentialConfigFunc func(*ExponentialBackoff)

// WithMultiplier sets the base of the exponential function (default is 2).
func WithMultiplier(multiplier float64) ExponentialConfigFunc {
	return func(b *ExponentialBackoff) { b.multiplier = multiplier }
}

// WithRandomFactor sets the random factor (default is 0, no randomness).
func WithRandomFactor(randFactor float64) ExponentialConfigFunc {
	return func(b *ExponentialBackoff) { b.randFactor = randFactor }
}

// ExponentialBackoff is an exponential backoff interval generator using
// a minimum and maximum interval. The base interval is given by the
// following function where n is the number of previous failed attempts in the
// current sequence.
//
//    `MinInterval * Multiplier ^ n`
//
// The value returned on each update is given by the following, where base is the
// value calculated above. A random factor of zero makes the generator deterministic.
// Some random jitter tends to work well in practice to avoid issues around a  thundering herd.
//
//     `min(MaxInterval, base +/- (base * RandFactor))`.
type ExponentialBackoff struct {
	minInterval time.Duration
	maxInterval time.Duration
	multiplier  float64
	randFactor  float64
	attempts    uint
	maxAttempts uint
}

// NewExponentialBackoff creates an ExponentialBackoff using the given minimum and maximum interval.
func NewExponentialBackoff(minInterval, maxInterval time.Duration, configs ...ExponentialConfigFunc) *ExponentialBackoff {
	if minInterval == 0 {
		// To avoid a divide by zero we set the minimum interval to something
		// that makes sense.
		minInterval = 1
	}

	backoff := &ExponentialBackoff{
		minInterval: minInterval,
		maxInterval: maxInterval,
		multiplier:  2,
		randFactor:  0,
		attempts:    0,
	}

	for _, config := range configs {
		config(backoff)
	}

	// Calculate and stash the maximum number of attempts now. This may be
	// expensive as it involves logs. We need to know the max number of
	// attempts now so we can shield ourselves from overflow when dealing
	// with larger intervals.

	// min * mult ^ n <     max
	//       mult ^ n <     max / min
	//              n < log(max / min) / log(mult)

	var (
		num         = math.Log(float64(maxInterval / minInterval))
		denom       = math.Log(backoff.multiplier)
		maxAttempts = uint(num / denom)
	)

	backoff.maxAttempts = maxAttempts
	return backoff
}

func (b *ExponentialBackoff) Reset() {
	b.attempts = 0
}

func (b *ExponentialBackoff) NextInterval() time.Duration {
	if b.attempts >= b.maxAttempts {
		return b.maxInterval
	}

	n := float64(b.attempts)
	b.attempts++

	return time.Duration(jitter(float64(b.minInterval)*math.Pow(b.multiplier, n), b.randFactor))
}

func (b *ExponentialBackoff) MinInterval() time.Duration {
	return b.minInterval
}

func (b *ExponentialBackoff) MaxInterval() time.Duration {
	return b.maxInterval
}

func (b *ExponentialBackoff) RandomFactor() float64 {
	return b.randFactor
}

func (b *ExponentialBackoff) Multiplier() float64 {
	return b.multiplier
}

func (b *ExponentialBackoff) Clone() Backoff {
	return &ExponentialBackoff{
		minInterval: b.minInterval,
		maxInterval: b.maxInterval,
		multiplier:  b.multiplier,
		randFactor:  b.randFactor,
		maxAttempts: b.maxAttempts,
	}
}

func jitter(value, ratio float64) float64 {
	min := value - (value * ratio)
	max := value + (value * ratio)

	return min + (max-min+1)*rand.Float64()
}
