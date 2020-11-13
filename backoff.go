package backoff

import (
	"math"
	"math/rand"
	"time"
)

type (
	// Backoff is the interface to a backoff interval generator.
	Backoff interface {
		// Mark the next call to NextInterval as the "first" retry in a sequence.
		// If the generated intervals are dependent on the number of consecutive
		// (unsuccessful) retries, previous retries should be forgotten here.
		Reset()

		// Generate the next backoff interval.
		NextInterval() time.Duration

		// Clone creates a copy of the backoff with a nil-internal state. This
		// allows a backoff object to be used as a prototype factory.
		Clone() Backoff
	}

	linearBackoff struct {
		minInterval time.Duration
		addInterval time.Duration
		maxInterval time.Duration
		current     time.Duration
	}

	exponentialBackoff struct {
		minInterval time.Duration
		maxInterval time.Duration
		multiplier  float64
		randFactor  float64
		attempts    uint
		maxAttempts uint
	}

	// ExponentialConfigFunc is a function used to initialize a new exponential backoff.
	ExponentialConfigFunc func(*exponentialBackoff)
)

// NewZeroBackoff creates a backoff interval generator which always returns
// a zero interval.
func NewZeroBackoff() Backoff {
	return NewConstantBackoff(0)
}

// NewConstantBackoff creates a backoff interval generator which always returns
// the same interval.
func NewConstantBackoff(interval time.Duration) Backoff {
	return NewLinearBackoff(interval, 0, interval)
}

// NewLinearBackoff creates a backoff interval generator which increases by a
// constant amount on each unsuccessful retry.
func NewLinearBackoff(minInterval, addInterval, maxInterval time.Duration) Backoff {
	b := &linearBackoff{
		minInterval: minInterval,
		addInterval: addInterval,
		maxInterval: maxInterval,
	}

	b.Reset()
	return b
}

func (b *linearBackoff) Reset() {
	b.current = b.minInterval
}

func (b *linearBackoff) NextInterval() time.Duration {
	current := b.current

	if current <= b.maxInterval-b.addInterval {
		b.current += b.addInterval
	} else {
		b.current = b.maxInterval
	}

	return current
}

func (b *linearBackoff) Clone() Backoff {
	return NewLinearBackoff(b.minInterval, b.addInterval, b.maxInterval)
}

//
// Exponential

// NewExponentialBackoff creates an exponential backoff interval generator using
// the given minimum and maximum interval. The base interval is given by the
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
func NewExponentialBackoff(minInterval, maxInterval time.Duration, configs ...ExponentialConfigFunc) Backoff {
	if minInterval == 0 {
		// To avoid a divide by zero we set the minimum interval to something
		// that makes sense.
		minInterval = 1
	}

	backoff := &exponentialBackoff{
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

// WithMultiplier sets the base of the exponential function (default is 2).
func WithMultiplier(multiplier float64) ExponentialConfigFunc {
	return func(b *exponentialBackoff) { b.multiplier = multiplier }
}

// WithRandomFactor sets the random factor (default is 0, no randomness).
func WithRandomFactor(randFactor float64) ExponentialConfigFunc {
	return func(b *exponentialBackoff) { b.randFactor = randFactor }
}

func (b *exponentialBackoff) Reset() {
	b.attempts = 0
}

func (b *exponentialBackoff) NextInterval() time.Duration {
	if b.attempts >= b.maxAttempts {
		return b.maxInterval
	}

	n := float64(b.attempts)
	b.attempts++

	return time.Duration(jitter(float64(b.minInterval)*math.Pow(b.multiplier, n), b.randFactor))
}

func (b *exponentialBackoff) Clone() Backoff {
	return &exponentialBackoff{
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
