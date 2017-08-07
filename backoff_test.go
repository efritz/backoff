package backoff

import (
	"testing"
	"time"
)

type BackoffSuite struct{}

func (s *BackoffSuite) TestZeroBackoff(t *testing.T) {
	b := NewZeroBackoff()
	testSequence(b, time.Millisecond, []uint{0, 0, 0, 0})
	b.Reset()
	testSequence(b, time.Millisecond, []uint{0, 0, 0, 0})
}

func (s *BackoffSuite) TestConstantBackoff(t *testing.T) {
	b1 := NewConstantBackoff(25 * time.Second)
	b2 := NewConstantBackoff(50 * time.Minute)

	testSequence(b1, time.Second, []uint{25, 25, 25, 25})
	b2.Reset()
	testSequence(b2, time.Minute, []uint{50, 50, 50, 50})

	testSequence(b1, time.Second, []uint{25, 25, 25, 25})
	b1.Reset()
	testSequence(b2, time.Minute, []uint{50, 50, 50, 50})
}

func (s *BackoffSuite) TestLinearBackoffMax(t *testing.T) {
	b := NewLinearBackoff(time.Millisecond, time.Millisecond, time.Millisecond*4)

	testSequence(b, time.Millisecond, []uint{1, 2, 3, 4, 4, 4})
	b.Reset()
	testSequence(b, time.Millisecond, []uint{1, 2, 3, 4, 4, 4})
}
