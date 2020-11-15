package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinearBackoff(t *testing.T) {
	t.Run("max", func(t *testing.T) {
		b := NewLinearBackoff(time.Millisecond, time.Millisecond, 4*time.Millisecond)
		testSequenceNew(t, b, time.Millisecond, []uint{1, 2, 3, 4, 4, 4})
		b.Reset()
		testSequenceNew(t, b, time.Millisecond, []uint{1, 2, 3, 4, 4, 4})
	})
	t.Run("min interval", func(t *testing.T) {
		b := NewLinearBackoff(time.Millisecond, 2*time.Millisecond, 4*time.Millisecond)
		assert.Equal(t, time.Millisecond, b.MinInterval())
	})
	t.Run("max interval", func(t *testing.T) {
		b := NewLinearBackoff(time.Millisecond, 2*time.Millisecond, 4*time.Millisecond)
		assert.Equal(t, 4*time.Millisecond, b.MaxInterval())
	})
	t.Run("add interval", func(t *testing.T) {
		b := NewLinearBackoff(time.Millisecond, 2*time.Millisecond, 4*time.Millisecond)
		assert.Equal(t, 2*time.Millisecond, b.AddInterval())
	})
	t.Run("clone", func(t *testing.T) {
		b := NewLinearBackoff(time.Millisecond, 2*time.Millisecond, 4*time.Millisecond)

		c := b.Clone()
		assert.NotSame(t, c, b)
		require.IsType(t, &LinearBackoff{}, c)
		lbc := c.(*LinearBackoff)
		assert.Equal(t, time.Millisecond, lbc.MinInterval())
		assert.Equal(t, 4*time.Millisecond, lbc.MaxInterval())
		assert.Equal(t, 2*time.Millisecond, lbc.AddInterval())
	})
}

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
