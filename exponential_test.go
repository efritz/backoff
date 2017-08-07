package backoff

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

type ExponentialSuite struct{}

func (s *ExponentialSuite) TestNonRandom(t *testing.T) {
	b := NewExponentialBackoff(
		2,
		0,
		time.Millisecond,
		time.Minute,
	)

	testSequence(b, time.Millisecond, []uint{1, 2, 4, 8, 16, 32})
	b.Reset()
	testSequence(b, time.Millisecond, []uint{1, 2, 4, 8, 16, 32})
}

func (s *ExponentialSuite) TestMax(t *testing.T) {
	b := NewExponentialBackoff(
		2,
		0,
		time.Millisecond,
		time.Millisecond*4,
	)

	testSequence(b, time.Millisecond, []uint{1, 2, 4, 4, 4, 4})
	b.Reset()
	testSequence(b, time.Millisecond, []uint{1, 2, 4, 4, 4, 4})
}

func (s *ExponentialSuite) TestRandomized(t *testing.T) {
	b := NewExponentialBackoff(
		2,
		.25,
		time.Millisecond,
		time.Minute,
	)

	testRandomizedSequence(b, time.Millisecond, .25, []uint{1, 2, 4, 8, 16, 32})
	b.Reset()
	testRandomizedSequence(b, time.Millisecond, .25, []uint{1, 2, 4, 8, 16, 32})
}

func (s *ExponentialSuite) TestOverflowLimit(t *testing.T) {
	b := NewExponentialBackoff(
		2,
		0,
		time.Millisecond,
		time.Minute,
	)

	for i := 0; i < 100; i++ {
		Expect(b.NextInterval()).To(BeNumerically(">=", 1000))
	}
}
