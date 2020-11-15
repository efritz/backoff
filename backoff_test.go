package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Sequence Assertion Helpers
func testSequenceNew(t *testing.T, b Backoff, base time.Duration, durations []uint) {
	testRandomizedSequenceNew(t, b, base, 0, durations)
}

func testRandomizedSequenceNew(t *testing.T, b Backoff, base time.Duration, ratio float64, durations []uint) {
	for _, duration := range durations {
		lo := int64(float64(base) * (1 - ratio) * float64(duration))
		hi := int64(float64(base) * (1 + ratio) * float64(duration))

		val := int64(b.NextInterval())
		assert.GreaterOrEqual(t, val, lo)
		assert.LessOrEqual(t, val, hi)
	}
}
