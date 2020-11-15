package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExponentialBackoff(t *testing.T) {
	t.Run("non-random", func(t *testing.T) {
		b := NewExponentialBackoff(time.Millisecond, time.Minute)
		testSequenceNew(t, b, time.Millisecond, []uint{1, 2, 4, 8, 16, 32})
		b.Reset()
		testSequenceNew(t, b, time.Millisecond, []uint{1, 2, 4, 8, 16, 32})
	})
	t.Run("zero time minimum", func(t *testing.T) {
		var eb *ExponentialBackoff
		assert.NotPanics(t, func() {
			eb = NewExponentialBackoff(0*time.Millisecond, time.Millisecond*4)
		})
		assert.Equal(t, 1*time.Nanosecond, eb.MinInterval())
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
	t.Run("min interval", func(t *testing.T) {
		b := NewExponentialBackoff(time.Millisecond, time.Minute)
		assert.Equal(t, 1*time.Millisecond, b.MinInterval())
	})
	t.Run("max interval", func(t *testing.T) {
		b := NewExponentialBackoff(time.Millisecond, time.Minute)
		assert.Equal(t, 1*time.Minute, b.MaxInterval())
	})
	t.Run("random factor", func(t *testing.T) {
		b := NewExponentialBackoff(
			time.Millisecond, time.Minute,
			WithRandomFactor(0.25),
		)
		assert.Equal(t, 0.25, b.RandomFactor())
	})
	t.Run("multiplier", func(t *testing.T) {
		b := NewExponentialBackoff(
			time.Millisecond, time.Minute,
			WithMultiplier(2),
		)
		assert.Equal(t, float64(2), b.Multiplier())
	})
	t.Run("clone", func(t *testing.T) {
		b := NewExponentialBackoff(
			time.Millisecond, time.Minute,
			WithRandomFactor(0.25),
			WithMultiplier(2),
		)

		c := b.Clone()
		assert.NotSame(t, b, c)
		require.IsType(t, &ExponentialBackoff{}, c)
		ebc := c.(*ExponentialBackoff)
		assert.Equal(t, time.Millisecond, ebc.MinInterval())
		assert.Equal(t, time.Minute, ebc.MaxInterval())
		assert.Equal(t, 0.25, ebc.RandomFactor())
		assert.Equal(t, float64(2), ebc.Multiplier())
	})
}
