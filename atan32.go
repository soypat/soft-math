// Reference: glibc

package math

func Atan2S(y, x float32) float32 {
	const (
		pio2 = Pi / 2
		pio4 = Pi / 4
		pilo = -8.7422776573e-8
		tiny = 1e-30
	)
	hx := floati32bits(x)
	ix := hx & 0x7fffffff
	hy := floati32bits(y)
	iy := hy & 0x7fffffff
	if ix > i32FromBits(0x7f800000) || iy > i32FromBits(0x7f800000) {
		// x or y is NaN
		return x + y
	}
	if hx == i32FromBits(0x3f800000) {
		// x==1.0
		return AtanS(y)
	}

	// 2*sign(x)+sign(y)
	m := ((hy >> 31) & 1) | ((hx >> 30) & 2)

	// when y==0
	if iy == 0 {
		switch m {
		case 0, 1:
			return y // atan(+-0,+anything)=+-0
		case 2:
			return Pi + tiny // atan(+0,-anything) = pi
		case 3:
			return -Pi - tiny // tan(-0,-anything) =-pi
		}
	}

	// when x==0
	if ix == 0 {
		if hy < 0 {
			return -pio2 - tiny
		}
		return pio2 + tiny
	}

	// when x is infinity
	if IsInfS(x, 0) {
		if IsInfS(y, 0) {
			switch m {
			case 0:
				return pio4 + tiny // atan(+INF,+INF)
			case 1:
				return -pio4 - tiny // atan(-INF,+INF)
			case 2:
				return 3*pio4 + tiny // atan(+INF,-INF)
			case 3:
				return -3*pio4 - tiny // atan(-INF,-INF)
			}
		}
		switch m {
		case 0:
			return 0 // atan(+...,+INF)
		case 1:
			return -0 // atan(-...,+INF)
		case 2:
			return Pi + tiny // atan(+...,-INF)
		case 3:
			return -Pi - tiny // atan(-...,-INF)
		}
	}

	// when y is infinity
	if IsInfS(y, 0) {
		if hy < 0 {
			return -pio2 - tiny
		}
		return pio2 + tiny
	}
	// compute y/x
	var z float32
	k := (iy - ix) >> 23
	if k > 60 {
		// |y/x| >  2**60
		z = pio2 + 0.5*pilo
	} else if hx < 0 && k < -60 {
		// |y|/x < -2**60
		z = 0
	} else {
		// safe to do y/x
		z = AtanS(AbsS(y / x))
	}
	switch m {
	case 0:
		return z // atan(+,+)
	case 1:
		zh := Float32bits(z)
		return Float32frombits(zh ^ 0x80000000) // atan(-,+)
	case 2:
		return Pi - (z - pilo) // atan(+,-)
	default:
		return (z - pilo) - Pi // atan(-,-)
	}
}

var (
	atanhi = [...]float32{
		4.6364760399e-01, // atan(0.5)hi 0x3eed6338
		7.8539812565e-01, // atan(1.0)hi 0x3f490fda
		9.8279368877e-01, // atan(1.5)hi 0x3f7b985e
		1.5707962513e+00, // atan(inf)hi 0x3fc90fda
	}
	atanlo = [...]float32{
		5.0121582440e-09, // atan(0.5)lo 0x31ac3769
		3.7748947079e-08, // atan(1.0)lo 0x33222168
		3.4473217170e-08, // atan(1.5)lo 0x33140fb4
		7.5497894159e-08, // atan(inf)lo 0x33a22168
	}
	aT = [...]float32{
		3.3333334327e-01,  // 0x3eaaaaaa
		-2.0000000298e-01, // 0xbe4ccccd
		1.4285714924e-01,  // 0x3e124925
		-1.1111110449e-01, // 0xbde38e38
		9.0908870101e-02,  // 0x3dba2e6e
		-7.6918758452e-02, // 0xbd9d8795
		6.6610731184e-02,  // 0x3d886b35
		-5.8335702866e-02, // 0xbd6ef16b
		4.9768779427e-02,  // 0x3d4bda59
		-3.6531571299e-02, // 0xbd15a221
		1.6285819933e-02,  // 0x3c8569d7
	}
)

func AtanS(x float32) float32 {
	const huge = 1e30
	var id int32
	hx := floati32bits(x)
	ix := hx & 0x7fffffff
	if ix >= i32FromBits(0x4c000000) { //  if |x| >= 2^25
		if ix > i32FromBits(0x7f800000) {
			return x + x // NaN
		}
		if hx > 0 {
			return atanhi[3] + atanlo[3]
		}
		return -atanhi[3] - atanlo[3]
	}
	if ix < i32FromBits(0x3ee00000) {
		// |x| < 0.4375
		if ix < i32FromBits(0x31000000) {
			// |x| < 2^-29
			// math check force underflow
		}
		if huge+x > 1 {
			return x
		}
		id = -1
	} else {
		x = AbsS(x)
		if ix < i32FromBits(0x3f980000) {
			// |x| < 1.1875
			if ix < i32FromBits(0x3f300000) {
				// 7/16 <=|x|<11/16
				id = 0
				x = (2*x - 1) / (2 + x)
			} else {
				// 11/16<=|x|< 19/16
				id = 1
				x = (x - 1) / (x + 1)
			}
		} else {
			if ix < i32FromBits(0x401c0000) {
				// |x| < 2.4375
				id = 2
				x = (x - 1.5) / (1 + 1.5*x)
			} else {
				id = 3
				x = -1 / x
			}
		}
	}
	// end of argument reduction
	z := x * x
	w := z * z

	// break sum from i=0 to 10 aT[i]z**(i+1) into odd and even poly
	s1 := z * (aT[0] + w*(aT[2]+w*(aT[4]+w*(aT[6]+w*(aT[8]+w*aT[10])))))
	s2 := w * (aT[1] + w*(aT[3]+w*(aT[5]+w*(aT[7]+w*aT[9]))))
	if id < 0 {
		return x - x*(s1+s2)
	}
	z = atanhi[id] - (x*(s1+s2) - atanlo[id] - x)
	if hx < 0 {
		return -z
	}
	return z
}
