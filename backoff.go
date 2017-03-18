package backoff

import "time"

type (
	// Backoff is the interface to a backoff interval generator.
	Backoff interface {
		// Mark the next call to NextInterval as the "first" retry in a sequence.
		// If the generated intervals are dependent on the number of consecutive
		// (unsuccessful) retries, previous retries should be forgotten here.
		Reset()

		// Generate the next backoff interval.
		NextInterval() time.Duration
	}

	linearBackoff struct {
		minInterval time.Duration
		addInterval time.Duration
		maxInterval time.Duration
		current     time.Duration
	}
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
