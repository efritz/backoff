package backoff

import (
	"testing"
	"time"

	"github.com/aphistic/sweet"
	junit "github.com/aphistic/sweet-junit"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&BackoffSuite{})
	})
}

type BackoffSuite struct{}

func (s *BackoffSuite) TestZeroBackoff(t sweet.T) {
	b := NewZeroBackoff()
	testSequence(b, time.Millisecond, []uint{0, 0, 0, 0})
	b.Reset()
	testSequence(b, time.Millisecond, []uint{0, 0, 0, 0})
}

func (s *BackoffSuite) TestConstantBackoff(t sweet.T) {
	b1 := NewConstantBackoff(25 * time.Second)
	b2 := NewConstantBackoff(50 * time.Minute)

	testSequence(b1, time.Second, []uint{25, 25, 25, 25})
	b2.Reset()
	testSequence(b2, time.Minute, []uint{50, 50, 50, 50})

	testSequence(b1, time.Second, []uint{25, 25, 25, 25})
	b1.Reset()
	testSequence(b2, time.Minute, []uint{50, 50, 50, 50})
}

func (s *BackoffSuite) TestLinearBackoffMax(t sweet.T) {
	b := NewLinearBackoff(time.Millisecond, time.Millisecond, time.Millisecond*4)
	testSequence(b, time.Millisecond, []uint{1, 2, 3, 4, 4, 4})
	b.Reset()
	testSequence(b, time.Millisecond, []uint{1, 2, 3, 4, 4, 4})
}

func (s *BackoffSuite) TestNonRandom(t sweet.T) {
	b := NewExponentialBackoff(time.Millisecond, time.Minute)
	testSequence(b, time.Millisecond, []uint{1, 2, 4, 8, 16, 32})
	b.Reset()
	testSequence(b, time.Millisecond, []uint{1, 2, 4, 8, 16, 32})
}

func (s *BackoffSuite) TestZeroTimeMinimum(t sweet.T) {
	Expect(func() {
		NewExponentialBackoff(0*time.Millisecond, time.Millisecond*4)
	}).ToNot(Panic())
}

func (s *BackoffSuite) TestMax(t sweet.T) {
	b := NewExponentialBackoff(time.Millisecond, time.Millisecond*4)
	testSequence(b, time.Millisecond, []uint{1, 2, 4, 4, 4, 4})
	b.Reset()
	testSequence(b, time.Millisecond, []uint{1, 2, 4, 4, 4, 4})
}

func (s *BackoffSuite) TestRandomized(t sweet.T) {
	b := NewExponentialBackoff(
		time.Millisecond,
		time.Minute,
		WithRandomFactor(0.25),
	)

	testRandomizedSequence(b, time.Millisecond, .25, []uint{1, 2, 4, 8, 16, 32})
	b.Reset()
	testRandomizedSequence(b, time.Millisecond, .25, []uint{1, 2, 4, 8, 16, 32})
}

func (s *BackoffSuite) TestOverflowLimit(t sweet.T) {
	b := NewExponentialBackoff(time.Millisecond, time.Minute)

	for i := 0; i < 100; i++ {
		Expect(b.NextInterval()).To(BeNumerically(">=", 1000))
	}
}

//
// Sequence Assertion Helpers

func testSequence(b Backoff, base time.Duration, durations []uint) {
	testRandomizedSequence(b, base, 0, durations)
}

func testRandomizedSequence(b Backoff, base time.Duration, ratio float64, durations []uint) {
	for _, duration := range durations {
		lo := time.Duration(float64(base) * (1 - ratio) * float64(duration))
		hi := time.Duration(float64(base) * (1 + ratio) * float64(duration))

		val := b.NextInterval()
		Expect(val).To(BeNumerically(">=", lo))
		Expect(val).To(BeNumerically("<=", hi))
	}
}
