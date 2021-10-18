package math

import (
	"math"
	"math/rand"
	"testing"
)

func floatInRange(r *rand.Rand, x1, x2 float64) float64 {
	return x1 + r.Float64()*(x2-x1)
}

type toleranceAsserter struct {
	t   testing.TB
	tol float64
}

func (as *toleranceAsserter) scalar(input, got, want float64) {
	if math.Abs(got-want) > as.tol {
		as.t.Errorf("with %g: got %g, want %g", input, got, want)
	}
}

func naiveAbs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
