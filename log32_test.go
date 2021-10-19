package math

import (
	"math"
	"math/rand"
	"testing"
)

func TestLogS(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	assert := toleranceAsserter{
		t,
		1e-6,
	}
	for i := 0; i < 1000; i++ {
		x := floatInRange(rng, 1e-4, 1e4)
		got := float64(LogS(float32(x)))
		want := math.Log(x)
		assert.scalarN(x, got, want)
	}
}
