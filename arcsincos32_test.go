package math

import (
	"math"
	"math/rand"
	"testing"
)

func TestAsinS(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	assert := toleranceAsserter{
		t,
		1e-7,
	}
	for i := 0; i < 10000; i++ {
		x := floatInRange(rng, -1e10, 1e10)
		got := float64(AsinS(float32(x)))
		want := math.Asin(x)
		assert.scalar(x, got, want)
	}
}

func TestAcosS(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	assert := toleranceAsserter{
		t,
		1e-7,
	}
	for i := 0; i < 10000; i++ {
		x := floatInRange(rng, -1e10, 1e10)
		got := float64(AcosS(float32(x)))
		want := math.Acos(x)
		assert.scalar(x, got, want)
	}
}
