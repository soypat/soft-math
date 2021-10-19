package math

// Modf returns integer and fractional floating-point numbers
// that sum to f. Both values have the same sign as f.
//
// Special cases are:
//	Modf(±Inf) = ±Inf, NaN
//	Modf(NaN) = NaN, NaN
func ModfS(f float32) (inte float32, frac float32) {
	if f < 1 {
		switch {
		case f < 0:
			inte, frac = ModfS(-f)
			return -inte, -frac
		case f == 0:
			return f, f // Return -0, -0 when f == -0
		}
		return 0, f
	}

	x := Float32bits(f)
	e := uint(x>>shiftS)&maskS - biasS
	const intpart = 32 - shiftS
	// Keep the top 12+e bits, the integer part; clear the rest.
	if e < 32-intpart {
		x &^= 1<<(32-intpart-e) - 1
	}
	inte = Float32frombits(x)
	frac = f - inte
	return
}
