package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type BackoffSuite struct{}

func TestZeroBackoff(t *testing.T) {
	b := NewZeroBackoff()
	testSequenceNew(t, b, time.Millisecond, []uint{0, 0, 0, 0})
	b.Reset()
	testSequenceNew(t, b, time.Millisecond, []uint{0, 0, 0, 0})
}

func TestConstantBackoff(t *testing.T) {
	b1 := NewConstantBackoff(25 * time.Second)
	b2 := NewConstantBackoff(50 * time.Minute)

	testSequenceNew(t, b1, time.Second, []uint{25, 25, 25, 25})
	b2.Reset()
	testSequenceNew(t, b2, time.Minute, []uint{50, 50, 50, 50})

	testSequenceNew(t, b1, time.Second, []uint{25, 25, 25, 25})
	b1.Reset()
	testSequenceNew(t, b2, time.Minute, []uint{50, 50, 50, 50})

}

func TestLinearBackoff(t *testing.T) {
	t.Run("max", func(t *testing.T) {
		b := NewLinearBackoff(time.Millisecond, time.Millisecond, time.Millisecond*4)
		testSequenceNew(t, b, time.Millisecond, []uint{1, 2, 3, 4, 4, 4})
		b.Reset()
		testSequenceNew(t, b, time.Millisecond, []uint{1, 2, 3, 4, 4, 4})
	})
}

func TestExponentialBackoff(t *testing.T) {
	t.Run("non-random", func(t *testing.T) {
		b := NewExponentialBackoff(time.Millisecond, time.Minute)
		testSequenceNew(t, b, time.Millisecond, []uint{1, 2, 4, 8, 16, 32})
		b.Reset()
		testSequenceNew(t, b, time.Millisecond, []uint{1, 2, 4, 8, 16, 32})
	})
	t.Run("zero time minimum", func(t *testing.T) {
		assert.NotPanics(t, func() {
			NewExponentialBackoff(0*time.Millisecond, time.Millisecond*4)
		})
	})
	t.Run("max", func(t *testing.T) {
		b := NewExponentialBackoff(time.Millisecond, time.Millisecond*4)
		testSequenceNew(t, b, time.Millisecond, []uint{1, 2, 4, 4, 4, 4})
		b.Reset()
		testSequenceNew(t, b, time.Millisecond, []uint{1, 2, 4, 4, 4, 4})
	})
	t.Run("randomized", func(t *testing.T) {
		b := NewExponentialBackoff(
			time.Millisecond,
			time.Minute,
			WithRandomFactor(0.25),
		)

		testRandomizedSequenceNew(t, b, time.Millisecond, .25, []uint{1, 2, 4, 8, 16, 32})
		b.Reset()
		testRandomizedSequenceNew(t, b, time.Millisecond, .25, []uint{1, 2, 4, 8, 16, 32})
	})
	t.Run("overflow limit", func(t *testing.T) {
		b := NewExponentialBackoff(time.Millisecond, time.Minute)

		for i := 0; i < 100; i++ {
			next := int(b.NextInterval())
			assert.GreaterOrEqual(t, next, 1000)
		}
	})
}

//
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
