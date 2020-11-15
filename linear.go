package backoff

import "time"

// LinearBackoff is a backoff interval generator which increases by a constant amount
// on each unsuccessful retry.
type LinearBackoff struct {
	minInterval time.Duration
	addInterval time.Duration
	maxInterval time.Duration
	current     time.Duration
}

// NewLinearBackoff creates a LinearBackoff with the given minimum and maximum interval which
// increases at the given additional interval.
func NewLinearBackoff(minInterval, addInterval, maxInterval time.Duration) *LinearBackoff {
	b := &LinearBackoff{
		minInterval: minInterval,
		addInterval: addInterval,
		maxInterval: maxInterval,
	}

	b.Reset()
	return b
}

// NewConstantBackoff creates a LinearBackoff which always returns the given interval.
func NewConstantBackoff(interval time.Duration) *LinearBackoff {
	return NewLinearBackoff(interval, 0, interval)
}

// NewZeroBackoff creates a LinearBackoff which always returns a zero interval.
func NewZeroBackoff() *LinearBackoff {
	return NewConstantBackoff(0)
}

func (b *LinearBackoff) Reset() {
	b.current = b.minInterval
}

func (b *LinearBackoff) NextInterval() time.Duration {
	current := b.current

	if current <= b.maxInterval-b.addInterval {
		b.current += b.addInterval
	} else {
		b.current = b.maxInterval
	}

	return current
}

func (b *LinearBackoff) MinInterval() time.Duration {
	return b.minInterval
}

func (b *LinearBackoff) AddInterval() time.Duration {
	return b.addInterval
}

func (b *LinearBackoff) MaxInterval() time.Duration {
	return b.maxInterval
}

func (b *LinearBackoff) Clone() Backoff {
	return NewLinearBackoff(b.minInterval, b.addInterval, b.maxInterval)
}
