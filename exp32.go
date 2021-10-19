// Reference: Go standard library.

package math

func LdexpS(frac float32, exp int) float32 {
	x := Float32bits(frac)
	// special cases
	if isZeroInfNaN(x) {
		return frac
	}

	frac, e := normalizeS(frac)
	exp += e

	exp += int(x>>shiftS)&maskS - biasS

	if exp < -biasS-shiftS {
		return CopysignS(0, frac) // underflow
	}
	if exp > biasS { // overflow
		if frac < 0 {
			return InfS(-1)
		}
		return InfS(1)
	}
	var m float32 = 1

	if exp < -biasS+1 { // denormal
		exp += shiftS + 1
		m = 1.0 / (1 << (shiftS + 1)) // 2**-53
	}
	x &^= maskS << shiftS
	x |= uint32(exp+biasS) << shiftS
	return m * Float32frombits(x)
}

// Frexp breaks f into a normalized fraction
// and an integral power of two.
// It returns frac and exp satisfying f == frac × 2**exp,
// with the absolute value of frac in the interval [½, 1).
//
// Special cases are:
//	Frexp(±0) = ±0, 0
//	Frexp(±Inf) = ±Inf, 0
//	Frexp(NaN) = NaN, 0
func FrexpS(f float32) (frac float32, exp int) {
	// special cases
	if isZeroInfNaN(Float32bits(f)) {
		return f, 0
	}
	f, exp = normalizeS(f)
	x := Float32bits(f)
	exp += int((x>>shiftS)&maskS) - biasS + 1
	x &^= maskS << shiftS
	x |= (-1 + biasS) << shiftS
	frac = Float32frombits(x)
	return
}

// Exp returns e**x, the base-e exponential of x.
//
// Special cases are:
//	Exp(+Inf) = +Inf
//	Exp(NaN) = NaN
// Very large values overflow to 0 or +Inf.
// Very small values underflow to 1.
func ExpS(x float32) float32 {
	const (
		Ln2Hi = 6.93147180369123816490e-01
		Ln2Lo = 1.90821492927058770002e-10
		Log2e = 1.44269504088896338700e+00

		// TODO fix hard coded values below.
		Overflow  = 7.09782712893383973096e+02
		Underflow = -7.45133219101941108420e+02
		NearZero  = 1.0 / (1 << 28) // 2**-28
	)

	// special cases
	switch {
	case IsNaNS(x) || IsInfS(x, 1):
		return x
	case IsInfS(x, -1):
		return 0
	case x > Overflow:
		return InfS(1)
	case x < Underflow:
		return 0
	case -NearZero < x && x < NearZero:
		return 1 + x
	}

	// reduce; computed as r = hi - lo for extra precision.
	var k int
	switch {
	case x < 0:
		k = int(Log2e*x - 0.5)
	case x > 0:
		k = int(Log2e*x + 0.5)
	}
	hi := x - float32(k)*Ln2Hi
	lo := float32(k) * Ln2Lo

	// compute
	return expmultiS(hi, lo, k)
}

// exp1 returns e**r × 2**k where r = hi - lo and |r| ≤ ln(2)/2.
func expmultiS(hi, lo float32, k int) float32 {
	const (
		P1 = 1.66666666666666657415e-01  /* 0x3FC55555; 0x55555555 */
		P2 = -2.77777777770155933842e-03 /* 0xBF66C16C; 0x16BEBD93 */
		P3 = 6.61375632143793436117e-05  /* 0x3F11566A; 0xAF25DE2C */
		P4 = -1.65339022054652515390e-06 /* 0xBEBBBD41; 0xC5D26BF1 */
		P5 = 4.13813679705723846039e-08  /* 0x3E663769; 0x72BEA4D0 */
	)

	r := hi - lo
	t := r * r
	c := r - t*(P1+t*(P2+t*(P3+t*(P4+t*P5))))
	y := 1 - ((lo - (r*c)/(2-c)) - hi)
	// TODO(rsc): make sure Ldexp can handle boundary k
	return LdexpS(y, k)
}

const (
	exp2TableBits = 5
	exp2PolyOrder = 3
	nExp2         = 1 << exp2TableBits
)

type exp2data struct {
	tab           [1 << exp2TableBits]uint64
	shiftScaled   float32
	poly          [exp2PolyOrder]float32
	shift         float32
	invln2_scaled float32
	polyScaled    [exp2PolyOrder]float32
}

var _exp2Data = exp2data{
	tab: [1 << exp2TableBits]uint64{
		0x3ff0000000000000, 0x3fefd9b0d3158574, 0x3fefb5586cf9890f, 0x3fef9301d0125b51,
		0x3fef72b83c7d517b, 0x3fef54873168b9aa, 0x3fef387a6e756238, 0x3fef1e9df51fdee1,
		0x3fef06fe0a31b715, 0x3feef1a7373aa9cb, 0x3feedea64c123422, 0x3feece086061892d,
		0x3feebfdad5362a27, 0x3feeb42b569d4f82, 0x3feeab07dd485429, 0x3feea47eb03a5585,
		0x3feea09e667f3bcd, 0x3fee9f75e8ec5f74, 0x3feea11473eb0187, 0x3feea589994cce13,
		0x3feeace5422aa0db, 0x3feeb737b0cdc5e5, 0x3feec49182a3f090, 0x3feed503b23e255d,
		0x3feee89f995ad3ad, 0x3feeff76f2fb5e47, 0x3fef199bdd85529c, 0x3fef3720dcef9069,
		0x3fef5818dcfba487, 0x3fef7c97337b9b5f, 0x3fefa4afa2a490da, 0x3fefd0765b6e4540,
	},
	shiftScaled:   0x1.8p+52 / nExp2,
	poly:          [exp2PolyOrder]float32{0x1.c6af84b912394p-5, 0x1.ebfce50fac4f3p-3, 0x1.62e42ff0c52d6p-1},
	shift:         0x1.8p+52,
	invln2_scaled: 0x1.71547652b82fep+0 * nExp2,
	polyScaled: [exp2PolyOrder]float32{
		0x1.c6af84b912394p-5 / nExp2 / nExp2 / nExp2, 0x1.ebfce50fac4f3p-3 / nExp2 / nExp2, 0x1.62e42ff0c52d6p-1 / nExp2,
	},
}
