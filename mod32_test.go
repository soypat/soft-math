package math

import (
	"math"
	"math/rand"
	"testing"
)

func TestModfS(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	assert := toleranceAsserter{
		t,
		1e-3,
	}
	for i := 0; i < 1000; i++ {
		x := floatInRange(rng, -1, 1)
		inte, frac := ModfS(float32(x))
		goti, gotf := float64(inte), float64(frac)
		wanti, wantf := math.Modf(x)
		assert.scalarN(x, goti, wanti)
		assert.scalarN(x, gotf, wantf)
	}
}
