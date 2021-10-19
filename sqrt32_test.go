package math

import (
	"math"
	"math/rand"
	"testing"
)

func TestSqrtS(t *testing.T) {
	const dom = 10000
	rng := rand.New(rand.NewSource(2))
	assert := toleranceAsserter{
		t,
		1e-5,
	}
	for i := 0; i < 100; i++ {
		x := floatInRange(rng, -dom, dom)
		got := float64(SqrtS(float32(x)))
		want := math.Sqrt(x)
		assert.scalar(x, got, want)
	}
}
