package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestZeroBackoff(t *testing.T) {
	b := NewZeroBackoff()
	testSequence(t, b, time.Millisecond, []uint{0, 0, 0, 0})
	b.Reset()
	testSequence(t, b, time.Millisecond, []uint{0, 0, 0, 0})
}

func TestConstantBackoff(t *testing.T) {
	b1 := NewConstantBackoff(25 * time.Second)
	b2 := NewConstantBackoff(50 * time.Minute)

	testSequence(t, b1, time.Second, []uint{25, 25, 25, 25})
	b2.Reset()
	testSequence(t, b2, time.Minute, []uint{50, 50, 50, 50})

	testSequence(t, b1, time.Second, []uint{25, 25, 25, 25})
	b1.Reset()
	testSequence(t, b2, time.Minute, []uint{50, 50, 50, 50})
}

func TestLinearBackoffMax(t *testing.T) {
	b := NewLinearBackoff(time.Millisecond, time.Millisecond, time.Millisecond*4)
	testSequence(t, b, time.Millisecond, []uint{1, 2, 3, 4, 4, 4})
	b.Reset()
	testSequence(t, b, time.Millisecond, []uint{1, 2, 3, 4, 4, 4})
}

func TestNonRandom(t *testing.T) {
	b := NewExponentialBackoff(time.Millisecond, time.Minute)
	testSequence(t, b, time.Millisecond, []uint{1, 2, 4, 8, 16, 32})
	b.Reset()
	testSequence(t, b, time.Millisecond, []uint{1, 2, 4, 8, 16, 32})
}

func TestMax(t *testing.T) {
	b := NewExponentialBackoff(time.Second, time.Minute*15)
	testSequence(t, b, time.Second, []uint{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 900, 900})
	b.Reset()
	testSequence(t, b, time.Second, []uint{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 900, 900})
}

func TestRandomized(t *testing.T) {
	b := NewExponentialBackoff(
		time.Millisecond,
		time.Minute,
		WithRandomFactor(0.25),
	)

	testRandomizedSequence(t, b, time.Millisecond, .25, []uint{1, 2, 4, 8, 16, 32})
	b.Reset()
	testRandomizedSequence(t, b, time.Millisecond, .25, []uint{1, 2, 4, 8, 16, 32})
}

func TestOverflowLimit(t *testing.T) {
	b := NewExponentialBackoff(time.Millisecond, time.Minute)

	for i := 0; i < 100; i++ {
		interval := b.NextInterval()
		assert.GreaterOrEqual(t, interval, time.Nanosecond*1000)
	}
}

//
// Sequence Assertion Helpers

func testSequence(t *testing.T, b Backoff, base time.Duration, durations []uint) {
	testRandomizedSequence(t, b, base, 0, durations)
}

func testRandomizedSequence(t *testing.T, b Backoff, base time.Duration, ratio float64, durations []uint) {
	for _, duration := range durations {
		lo := time.Duration(float64(base) * (1 - ratio) * float64(duration))
		hi := time.Duration(float64(base) * (1 + ratio) * float64(duration))

		val := b.NextInterval()
		assert.GreaterOrEqual(t, val, lo)
		assert.LessOrEqual(t, val, hi)
	}
}
