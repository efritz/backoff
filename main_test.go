package backoff

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type BackoffSuite struct{}

var _ = Suite(&BackoffSuite{})

//
//

func testSequence(c *C, b Backoff, base time.Duration, durations []uint) {
	testRandomizedSequence(c, b, base, 0, durations)
}

func testRandomizedSequence(c *C, b Backoff, base time.Duration, ratio float64, durations []uint) {
	for _, duration := range durations {
		v := b.NextInterval()

		c.Assert(v >= time.Duration(float64(base)*float64(duration)*(1-ratio)), Equals, true)
		c.Assert(v <= time.Duration(float64(base)*float64(duration)*(1+ratio)), Equals, true)
	}
}

//
//

type MockBackoff struct {
	f1 func()
	f2 func() time.Duration
}

func NewMockBackoff(f1 func(), f2 func() time.Duration) Backoff {
	return &MockBackoff{
		f1: f1,
		f2: f2,
	}
}

func (m *MockBackoff) Reset() {
	m.f1()
}

func (m *MockBackoff) NextInterval() time.Duration {
	return m.f2()
}
