package math

func AcosS(x float32) float32 {
	const (
		pio2_hi = Pi / 2
		pio2_lo = -4.37113900018624283e-8

		pS0 = 1.6666667163e-01  // 0x3e2aaaab
		pS1 = -3.2556581497e-01 // 0xbea6b090
		pS2 = 2.0121252537e-01  // 0x3e4e0aa8
		pS3 = -4.0055535734e-02 // 0xbd241146
		pS4 = 7.9153501429e-04  // 0x3a4f7f04
		pS5 = 3.4793309169e-05  // 0x3811ef08
		qS1 = -2.4033949375e+00 // 0xc019d139
		qS2 = 2.0209457874e+00  // 0x4001572d
		qS3 = -6.8828397989e-01 // 0xbf303361
		qS4 = 7.7038154006e-02  // 0x3d9dc62e
	)
	var z, p, q, r, w, s, c, df float32
	hx := floati32bits(x)
	ix := hx & 0x7fffffff

	if ix == 0x3f800000 { // |x|==1
		if hx > 0 {
			return 0.0 // acos(1) = 0
		} else {
			return Pi + 2*pio2_lo // acos(-1)= pi
		}
	} else if ix > 0x3f800000 { //|x| >= 1
		return (x - x) / (x - x) // acos(|x|>1) is NaN
	}

	if ix < 0x3f000000 { // |x| < 0.5
		if ix <= 0x32800000 {
			return pio2_hi + pio2_lo //if|x|<=2**-26
		}
		z = x * x
		p = z * (pS0 + z*(pS1+z*(pS2+z*(pS3+z*(pS4+z*pS5)))))
		q = 1 + z*(qS1+z*(qS2+z*(qS3+z*qS4)))
		r = p / q
		return pio2_hi - (x - (pio2_lo - x*r))
	} else if hx < 0 { // x < -0.5
		z = (1 + x) * 0.5
		p = z * (pS0 + z*(pS1+z*(pS2+z*(pS3+z*(pS4+z*pS5)))))
		q = 1 + z*(qS1+z*(qS2+z*(qS3+z*qS4)))
		s = SqrtS(z)
		r = p / q
		w = r*s - pio2_lo
		return Pi - 2*(s+w)
	} else { // x > 0.5
		z = (1 - x) * 0.5
		s = SqrtS(z)
		df = s
		idf := floati32bits(df)
		msk := 0xfffff000
		df = floati32frombits(idf & int32(msk))
		c = (z - df*df) / (s + df)
		p = z * (pS0 + z*(pS1+z*(pS2+z*(pS3+z*(pS4+z*pS5)))))
		q = 1 + z*(qS1+z*(qS2+z*(qS3+z*qS4)))
		r = p / q
		w = r*s + c
		return 2 * (df + w)
	}
}

func AsinS(x float32) float32 {
	const (
		huge    = 1e30
		pio2_hi = Pi / 2
		pio2_lo = -4.37113900018624283e-8
		pio4_hi = Pi / 4
		p0      = 1.666675248e-1
		p1      = 7.495297643e-2
		p2      = 4.547037598e-2
		p3      = 2.417951451e-2
		p4      = 4.216630880e-2
	)

	hx := floati32bits(x)
	ix := hx & 0x7fffffff
	if ix == 0x3f800000 {
		// asin(1)=+-pi/2 with inexact
		return x*pio2_hi + x*pio2_lo
	} else if ix > 0x3f800000 { // |x|>= 1
		// asin(|x|>1) is NaN
		return (x - x) / (x - x)
	} else if ix < 0x3f000000 { //|x|<0.5
		if ix < 0x32000000 { // if |x| < 2**-27
			// math_check_force_underflow(x)
			if huge+x > 1 {
				// return x with inexact if x!=0
				return x
			}
		} else {
			t := x * x
			w := t * (p0 + t*(p1+t*(p2+t*(p3+t*p4))))
			return x + x*w
		}
	}

	// 1> |x|>= 0.5
	w := 1 - AbsS(x)
	t := w * 0.5
	p := t * (p0 + t*(p1+t*(p2+t*(p3+t*p4))))
	s := SqrtS(t)
	if ix >= 0x3F79999A { /* if |x| > 0.975 */
		t = pio2_hi - (2.0*(s+s*p) - pio2_lo)
	} else {
		w = s
		iw := floati32bits(w)
		msk := 0xfffff000
		w = floati32frombits(iw & int32(msk))
		c := (t - w*w) / (s + w)
		r := p
		p = 2*s*r - (pio2_lo - 2*c)
		q := pio4_hi - 2*w
		t = pio4_hi - (p - q)
	}
	if hx > 0 {
		return t
	}
	return -t
}
