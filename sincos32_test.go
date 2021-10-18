package math

import (
	"math"
	"math/rand"
	"testing"
)

func TestCosS(t *testing.T) {
	const dom = 10000
	rng := rand.New(rand.NewSource(2))
	assert := toleranceAsserter{
		t,
		1e-3,
	}

	for i := 0; i < 100; i++ {
		x := floatInRange(rng, -dom, dom)
		got := float64(CosS(float32(x)))
		want := math.Cos(x)
		assert.scalar(x, got, want)
	}
}
