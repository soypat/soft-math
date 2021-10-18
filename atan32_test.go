package math

import (
	"math"
	"math/rand"
	"testing"
)

func TestAtanS(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	assert := toleranceAsserter{
		t,
		1e-7,
	}
	for i := 0; i < 10000; i++ {
		x := floatInRange(rng, -1e10, 1e10)
		got := float64(AtanS(float32(x)))
		want := math.Atan(x)
		assert.scalar(x, got, want)
	}
}

func TestAtan2S(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	assert := toleranceAsserter{
		t,
		1e-1,
	}
	for i := 0; i < 1000; i++ {
		x := floatInRange(rng, -1e10, 1e10)
		y := floatInRange(rng, -1e10, 1e10)
		got := float64(Atan2S(float32(y), float32(x)))
		want := math.Atan2(y, x)
		assert.scalar(y/x, got, want)
	}
}

func TestAbsS(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	assert := toleranceAsserter{
		t,
		1e-3,
	}
	for i := 0; i < 1000; i++ {
		x := float32(floatInRange(rng, -1e10, 1e10))
		got := AbsS(x)
		want := naiveAbs(x)
		assert.scalar(float64(x), float64(got), float64(want))
	}
}
