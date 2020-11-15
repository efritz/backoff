package backoff

import (
	"time"
)

// Backoff is the interface to a backoff interval generator.
type Backoff interface {
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
