package backoff

import "time"

// BackOff is the interface to a back-off interval generator.
type BackOff interface {
	// Mark the next call to NextInterval as the "first" retry in a sequence.
	// If the generated intervals are dependent on the number of consecutive
	// (unsuccessful) retries, previous retries should be forgotten here.
	Reset()

	// Generate the next back-off interval.
	NextInterval() time.Duration
}

//
//

// A back-off interval generator which always returns a zero interval.
func NewZeroBackOff() BackOff {
	return NewConstantBackOff(0)
}

// A back-off interval generator which always returns the same interval.
func NewConstantBackOff(interval time.Duration) BackOff {
	return NewLinearBackOff(interval, 0, interval)
}

// A back-off interval generator which increases by a constant amount on
// each unsuccessful retry.
func NewLinearBackOff(minInterval, addInterval, maxInterval time.Duration) BackOff {
	b := &linearBackOff{
		minInterval: minInterval,
		addInterval: addInterval,
		maxInterval: maxInterval,
	}

	b.Reset()
	return b
}

//
//

type linearBackOff struct {
	minInterval time.Duration
	addInterval time.Duration
	maxInterval time.Duration
	current     time.Duration
}

func (b *linearBackOff) Reset() {
	b.current = b.minInterval
}

func (b *linearBackOff) NextInterval() time.Duration {
	current := b.current

	if current <= b.maxInterval-b.addInterval {
		b.current += b.addInterval
	} else {
		b.current = b.maxInterval
	}

	return current
}
