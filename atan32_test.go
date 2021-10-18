package math

import (
	"math"
	"math/rand"
	"testing"
)

func TestAtans(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	assert := toleranceAsserter{
		t,
		1e-7,
	}
	for i := 0; i < 10000; i++ {
		x := floatInRange(rng, -1e10, 1e10)
		got := float64(Atans(float32(x)))
		want := math.Atan(x)
		assert.scalar(got, want)
	}
}

func TestAtan2s(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	assert := toleranceAsserter{
		t,
		1e-1,
	}
	for i := 0; i < 1000; i++ {
		x := floatInRange(rng, -1e10, 1e10)
		y := floatInRange(rng, -1e10, 1e10)
		got := float64(Atan2s(float32(y), float32(x)))
		want := math.Atan2(y, x)
		assert.scalar(got, want)
	}
}

func TestAbss(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	assert := toleranceAsserter{
		t,
		1e-3,
	}
	for i := 0; i < 1000; i++ {
		x := float32(floatInRange(rng, -1e10, 1e10))
		got := Abss(x)
		want := naiveAbs(x)
		assert.scalar(float64(got), float64(want))
	}
}

func floatInRange(r *rand.Rand, x1, x2 float64) float64 {
	return x1 + r.Float64()*(x2-x1)
}

type toleranceAsserter struct {
	t   testing.TB
	tol float64
}

func (as *toleranceAsserter) scalar(got, want float64) {
	if math.Abs(got-want) > as.tol {
		as.t.Errorf("got %g, want %g", got, want)
	}
}

func naiveAbs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
