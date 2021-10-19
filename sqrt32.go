// Reference: glibc
// sysdeps/ieee754/flt-32/...

package math

func SqrtS(x float32) float32 {
	_sign := 0x80000000
	sign := int32(_sign)
	ix := floati32bits(x)

	if ix&0x7f800000 == 0x7f800000 {
		// sqrt(NaN)=NaN, sqrt(+inf)=+inf, sqrt(-inf)=sNaN
		return x*x + x
	}

	// Take care of zero.
	if ix <= 0 {
		if ix&^sign == 0 {
			return x // sqrt(+-0) = +-0
		} else if ix < 0 {
			return (x - x) / (x - x) // sqrt(-ve) = sNaN
		}
	}

	// Normalize x.
	var i int32
	m := ix >> 23
	if m == 0 {
		// Subnormal x.
		for i = 0; ix&0x00800000 == 0; i++ {
			ix <<= 1
		}
		m -= i - 1
	}
	m -= 127 // unbias exponent.

	ix = ix&0x007fffff | 0x00800000
	if m&1 != 0 {
		// Odd m, double x to make it even.
		ix += ix
	}
	m >>= 1

	// Generate sqrt(x) bit by bit.
	ix += ix
	var q, r, s int32
	r = 0x01000000

	for r != 0 {
		t := s + r
		if t <= ix {
			s = t + r
			ix -= t
			q += r
		}
		ix += ix
		r >>= 1
	}

	// Use floating add to find out rounding direction.
	var z float32
	if ix != 0 {
		z = 0x1p0 - 0x1.4484cp-100
		if z >= 0x1p0 {
			z = 0x1p0 + 0x1.4484cp-100
		}
		if z > 0x1p0 {
			q += 2
		} else {
			q += q & 1
		}
	}
	ix = q>>1 + 0x3f000000
	ix += m << 23
	return floati32frombits(ix)
}
