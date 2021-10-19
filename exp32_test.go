package math

import (
	"math"
	"math/rand"
	"testing"
)

func TestLdexpS(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	assert := toleranceAsserter{
		t,
		1e-6,
	}
	for i := 0; i < 500; i++ {
		x := floatInRange(rng, -1, 1)
		for _, exp := range []int{0, 1, 2, 3, 5, 23, 41} {
			got := float64(LdexpS(float32(x), exp))
			want := math.Ldexp(x, exp)
			assert.scalarN(x, got, want)
		}
	}
}

func TestFrexpS(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	assert := toleranceAsserter{
		t,
		1e-7,
	}
	for i := 0; i < 1000; i++ {
		x := floatInRange(rng, -1, 1)
		got32, gote := FrexpS(float32(x))
		got := float64(got32)
		want, wante := math.Frexp(x)
		assert.scalarN(x, got, want)
		if gote != wante {
			t.Errorf("")
		}
	}
}

func TestExpS(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	assert := toleranceAsserter{
		t,
		1e-7,
	}
	for i := 0; i < 1000; i++ {
		x := floatInRange(rng, -1, 1)
		got := float64(ExpS(float32(x)))
		want := math.Exp(x)
		assert.scalarN(x, got, want)
	}
}
