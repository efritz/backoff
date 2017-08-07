package backoff

import (
	"testing"
	"time"

	"github.com/aphistic/sweet"
	"github.com/aphistic/sweet-junit"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&BackoffSuite{})
		s.AddSuite(&ExponentialSuite{})
	})
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
