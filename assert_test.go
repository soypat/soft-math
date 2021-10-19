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

// Absolute tolerance   math.Abs(got-want)
func (as *toleranceAsserter) scalar(input, got, want float64) {
	if math.Abs(got-want) > as.tol {
		as.t.Errorf("with %g: got %g, want %g (diff %+.2g)", input, got, want, want-got)
	}
}

// Fractional tolerance  math.Abs(got-want)/want
func (as *toleranceAsserter) scalarN(input, got, want float64) {
	if math.Abs(got-want)/want > as.tol {
		as.t.Errorf("with %g: got %g, want %g (diff %+.2g%%)", input, got, want, (got-want)/want*100)
	}
}

// Fractional tolerance  math.Abs(got-want)/want
func (as *toleranceAsserter) int(input float64, got, want int) {
	if got != want {
		as.t.Errorf("with %g: got %d, want %d (diff %+d)", input, got, want, got-want)
	}
}

func naiveAbs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
