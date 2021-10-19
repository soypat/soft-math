package math

import (
	"math"
	"math/rand"
	"testing"
)

func TestPowS(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	assert := toleranceAsserter{
		t,
		1e-5,
	}
	for i := 0; i < 1000; i++ {
		x := floatInRange(rng, -1e2, 1e2)
		y := floatInRange(rng, 1e-4, 20)
		got := float64(PowS(float32(x), float32(y)))
		want := math.Pow(x, y)
		assert.scalarN(x, got, want)
	}
}
