package backoff

import (
	"time"

	. "gopkg.in/check.v1"
)

func (s *BackoffSuite) TestNonRandom(c *C) {
	b := NewExponentialBackoff(
		2,
		0,
		time.Millisecond,
		time.Minute,
	)

	testSequence(c, b, time.Millisecond, []uint{1, 2, 4, 8, 16, 32})
	b.Reset()
	testSequence(c, b, time.Millisecond, []uint{1, 2, 4, 8, 16, 32})
}

func (s *BackoffSuite) TestMax(c *C) {
	b := NewExponentialBackoff(
		2,
		0,
		time.Millisecond,
		time.Millisecond*4,
	)

	testSequence(c, b, time.Millisecond, []uint{1, 2, 4, 4, 4, 4})
	b.Reset()
	testSequence(c, b, time.Millisecond, []uint{1, 2, 4, 4, 4, 4})
}

func (s *BackoffSuite) TestRandomized(c *C) {
	b := NewExponentialBackoff(
		2,
		.25,
		time.Millisecond,
		time.Minute,
	)

	testRandomizedSequence(c, b, time.Millisecond, .25, []uint{1, 2, 4, 8, 16, 32})
	b.Reset()
	testRandomizedSequence(c, b, time.Millisecond, .25, []uint{1, 2, 4, 8, 16, 32})
}

func (s *BackoffSuite) TestOverflowLimit(c *C) {
	b := NewExponentialBackoff(
		2,
		0,
		time.Millisecond,
		time.Minute,
	)

	for i := 0; i < 100; i++ {
		c.Assert(b.NextInterval() >= 1000, Equals, true)
	}
}
