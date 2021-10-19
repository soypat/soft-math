package math

const (
	signBias = 1 << (exp2TableBits + 11)
)

func PowS(x, y float32) float32 {
	return pow(x, y)
}

// checkint Returns 0 if not int, 1 if odd int, 2 if even int.  The argument is
// the bit representation of a non-zero finite floating-point value.
func checkint(ix uint32) int {
	e := ix >> 23 & 0xff
	switch {
	case e < 0x7f:
		return 0
	case e > 0x7f+23:
		return 2
	case ix&((1<<(0x7f+23-e))-1) != 0:
		return 0
	case ix&(1<<(0x7f+23-e)) != 0:
		return 1
	default:
		return 2
	}
}

func pow(x, y float32) float32 {
	iy := Float32bits(y)
	switch {
	case y == 0 || x == 1:
		return 1
	case y == 1:
		return x
	case IsNaNS(x) || IsNaNS(y):
		return NaNS()
	case x == 0:
		switch {
		case y < 0:
			if checkint(iy) == 1 {
				return CopysignS(InfS(1), x)
			}
			return InfS(1)
		case y > 0:
			if checkint(iy) == 1 {
				return x
			}
			return 0
		}
	case IsInfS(y, 0):
		switch {
		case x == -1:
			return 1
		case (AbsS(x) < 1) == IsInfS(y, 1):
			return 0
		default:
			return InfS(1)
		}
	case IsInfS(x, 0):
		if IsInfS(x, -1) {
			return PowS(1/x, -y) // Pow(-0, -y)
		}
		switch {
		case y < 0:
			return 0
		case y > 0:
			return InfS(1)
		}
	case y == 0.5:
		return SqrtS(x)
	case y == -0.5:
		return 1 / SqrtS(x)
	}

	yi, yf := ModfS(AbsS(y))
	if yf != 0 && x < 0 {
		return NaNS()
	}
	if yi >= 1<<31 {
		// yi is a large even int that will lead to overflow (or underflow to 0)
		// for all x except -1 (x == 1 was handled earlier)
		switch {
		case x == -1:
			return 1
		case (AbsS(x) < 1) == (y > 0):
			return 0
		default:
			return InfS(1)
		}
	}

	// ans = a1 * 2**ae (= 1 for now).
	var a1 float32 = 1
	ae := 0

	// ans *= x**yf
	if yf != 0 {
		if yf > 0.5 {
			yf--
			yi++
		}
		a1 = ExpS(yf * LogS(x))
	}

	// ans *= x**yi
	// by multiplying in successive squarings
	// of x according to bits of yi.
	// accumulate powers of two into exp.
	const intpart = 32 - shiftS
	x1, xe := FrexpS(x)
	for i := int64(yi); i != 0; i >>= 1 {
		if xe < -1<<intpart || 1<<intpart < xe {
			// catch xe before it overflows the left shift below
			// Since i !=0 it has at least one bit still set, so ae will accumulate xe
			// on at least one more iteration, ae += xe is a lower bound on ae
			// the lower bound on ae exceeds the size of a float64 exp
			// so the final call to Ldexp will produce under/overflow (0/Inf)
			ae += xe
			break
		}
		if i&1 == 1 {
			a1 *= x1
			ae += xe
		}
		x1 *= x1
		xe <<= 1
		if x1 < .5 {
			x1 += x1
			xe--
		}
	}

	// ans = a1*2**ae
	// if y < 0 { ans = 1 / ans }
	// but in the opposite order
	if y < 0 {
		a1 = 1 / a1
		ae = -ae
	}
	return LdexpS(a1, ae)
}

func powglibc(x, y float32) float32 {
	ix := Float32bits(x)
	iy := Float32bits(y)
	var sgn uint32
	if ix-0x00800000 >= 0x7f800000-0x00800000 || isZeroInfNaN(iy) {
		if isZeroInfNaN(iy) {
			switch {
			case (2*iy == 0):
				if IsNaNS(x) {
					return x + y
				}
				return 1
			case (ix == uvoneS):
				if IsNaNS(y) {
					return x + y
				}
				return 1
			case (2*ix > 2*0x7f800000 || 2*iy > 2*0x7f800000):
				return x + y
			case (2*ix == 2*0x3f800000):
				return 1
			case (2*ix < 2*0x3f800000) == (iy&0x80000000 == 0):
				return 0 // |x|<1 && y==inf or |x|>1 && y==-inf.
			default:
				return y * y
			}
		}

		if isZeroInfNaN(ix) {
			x2 := x * x
			if ix&0x80000000 != 0 && checkint(iy) == 1 {
				x2 = -x2
				sgn = 1
			}
			// Divide by zero
			if 2*ix == 0 && iy&0x80000000 != 0 {
				return NaNS() // TODO __math_divzerof (sign_bias);
			}
			if iy&0x80000000 != 0 {
				return 1 / x2
			}
			return x2
		}

		// x and y are non-zero finite
		if ix&0x80000000 != 0 {
			/* Finite x < 0.  */
			yint := checkint(iy)

			if yint == 0 {
				return NaNS() // TODO __math_invalidf (x);
			}

			if yint == 1 {
				sgn = signBias
			}
			ix &= 0x7fffffff
		}
		if ix < 0x00800000 {
			/* Normalize subnormal x so exponent becomes negative.  */
			ix = Float32bits(x * 0x1p23)
			ix &= 0x7fffffff
			ix -= 23 << 23
		}
	}
	logx := iLog2S(ix)
	ylogx := y * logx // Note: cannot overflow, y is single prec.
	if Float32bits(ylogx)>>47&0xffff >= Float32bits(126*powScale)>>47 {
		switch {
		case ylogx > 0x1.fffffffd1d571p+6*powScale:
			// |y*log(x)| >= 126.
			return InfS(int(sgn)) // TODO verify this line?
		case ylogx > 0x1.fffffffa3aae2p+6*powScale:
			// |x^y| > 0x1.fffffep127, check if we round away from 0.
			// TODO powf overflow handling in non-nearest rounding mode
		case ylogx <= -150*powScale:
			return 0 // TODO __math_uflowf (sign_bias);
		}
	}
	return 0
}

func iLog2S(ix uint32) float32 {
	// TODO: original routine uses double for greater precision.
	const (
		off          = 0x3f330000
		powScaleBits = 0
		N            = 1 << powLog2TableBits
	)
	// var z, r, r2, r4, p, q, y, y0, invc, logc float32
	T := &_powlog2table.tab
	A := &_powlog2table.poly
	// x = 2^k z; where z is in range [OFF,2*OFF] and exact.
	// The range is split into N subintervals.
	// The ith subinterval contains z and c is near its center.
	tmp := ix - off
	i := (tmp >> (23 - powLog2TableBits)) % N
	top := tmp & 0xff800000
	iz := ix - top
	k := top >> (23 - powScaleBits) // Arithmetic shift.
	invc := T[i].invc
	logc := T[i].logc
	z := Float32frombits(iz)

	// log2(x) = log1p(z/c-1)/ln2 + log2(c) + k
	r := z*invc - 1
	y0 := logc + float32(k)

	// Pipelined polynomial evaluation to approximate log1p(r)/ln2.
	r2 := r * r
	y := A[0]*r + A[1]
	p := A[2]*r + A[3]
	r4 := r2 * r2
	q := A[4]*r + y0
	q = p*r2 + q
	y = y*r4 + q
	return y
}

const (
	powPolyOrder     = 5
	powLog2TableBits = 4
	powScaleBits     = 0
	powScale         = 1 << powScaleBits
)

type powRecord struct {
	invc, logc float32
}
type powLog2Data struct {
	tab  [1 << powLog2TableBits]powRecord
	poly [powPolyOrder]float32
}

var _powlog2table = powLog2Data{
	tab: [1 << powLog2TableBits]powRecord{
		{0x1.661ec79f8f3bep+0, -0x1.efec65b963019p-2 * powScale},
		{0x1.571ed4aaf883dp+0, -0x1.b0b6832d4fca4p-2 * powScale},
		{0x1.49539f0f010bp+0, -0x1.7418b0a1fb77bp-2 * powScale},
		{0x1.3c995b0b80385p+0, -0x1.39de91a6dcf7bp-2 * powScale},
		{0x1.30d190c8864a5p+0, -0x1.01d9bf3f2b631p-2 * powScale},
		{0x1.25e227b0b8eap+0, -0x1.97c1d1b3b7afp-3 * powScale},
		{0x1.1bb4a4a1a343fp+0, -0x1.2f9e393af3c9fp-3 * powScale},
		{0x1.12358f08ae5bap+0, -0x1.960cbbf788d5cp-4 * powScale},
		{0x1.0953f419900a7p+0, -0x1.a6f9db6475fcep-5 * powScale},
		{0x1p+0, 0x0p+0 * powScale},
		{0x1.e608cfd9a47acp-1, 0x1.338ca9f24f53dp-4 * powScale},
		{0x1.ca4b31f026aap-1, 0x1.476a9543891bap-3 * powScale},
		{0x1.b2036576afce6p-1, 0x1.e840b4ac4e4d2p-3 * powScale},
		{0x1.9c2d163a1aa2dp-1, 0x1.40645f0c6651cp-2 * powScale},
		{0x1.886e6037841edp-1, 0x1.88e9c2c1b9ff8p-2 * powScale},
		{0x1.767dcf5534862p-1, 0x1.ce0a44eb17bccp-2 * powScale},
	},
	poly: [powPolyOrder]float32{
		0x1.27616c9496e0bp-2 * powScale,
		-0x1.71969a075c67ap-2 * powScale,
		0x1.ec70a6ca7baddp-2 * powScale,
		-0x1.7154748bef6c8p-1 * powScale,
		0x1.71547652ab82bp0 * powScale,
	},
}
