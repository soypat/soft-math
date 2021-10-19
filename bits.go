// Taken from the Go standard Library. The Go Authors.

package math

const (
	uvnan    = 0x7FF8000000000001
	uvinf    = 0x7FF0000000000000
	uvneginf = 0xFFF0000000000000
	uvone    = 0x3FF0000000000000
	mask     = 0x7FF
	shift    = 64 - 11 - 1
	bias     = 1023
	signMask = 1 << 63
	fracMask = 1<<shift - 1
)

// 32bit float constants
const (
	uvnanS    uint32 = 0x7FE00000
	uvinfS    uint32 = 0x7F800000
	uvneginfS uint32 = 0xFF800000
	uvoneS    uint32 = 0x3F800000
	maskS            = 0xFF
	biasS            = 127
	shiftS           = 32 - 8 - 1
	fracMaskS        = 1<<shiftS - 1
)

// Abs returns the absolute value of x.
//
// Special cases are:
//	Abs(±Inf) = +Inf
//	Abs(NaN) = NaN
func Abs(x float64) float64 {
	return Float64frombits(Float64bits(x) &^ (1 << 63))
}

func AbsS(x float32) float32 {
	return Float32frombits(Float32bits(x) &^ (1 << 31))
}

// Top 12 bits of the float representation with the sign bit cleared.
func absToP12s(x float32) uint32 {
	return (Float32bits(x) >> 20) & mask
}

// Inf returns positive infinity if sign >= 0, negative infinity if sign < 0.
func Inf(sign int) float64 {
	var v uint64
	if sign >= 0 {
		v = uvinf
	} else {
		v = uvneginf
	}
	return Float64frombits(v)
}

// Inf returns positive infinity if sign >= 0, negative infinity if sign < 0.
func InfS(sign int) float32 {
	var v uint32
	if sign >= 0 {
		v = uvinfS
	} else {
		v = uvneginfS
	}
	return Float32frombits(v)
}

// NaN returns an IEEE 754 ``not-a-number'' value.
func NaN() float64 { return Float64frombits(uvnan) }

// NaN returns an IEEE 754 ``not-a-number'' value.
func NaNS() float32 { return Float32frombits(uvnanS) }

// IsNaN reports whether f is an IEEE 754 ``not-a-number'' value.
func IsNaN(f float64) (is bool) {
	// IEEE 754 says that only NaNs satisfy f != f.
	// To avoid the floating-point hardware, could use:
	//	x := Float64bits(f);
	//	return uint32(x>>shift)&mask == mask && x != uvinf && x != uvneginf
	return f != f
}

// IsNaNS reports whether f is an IEEE 754 ``not-a-number'' value.
func IsNaNS(f float32) (is bool) {
	x := Float32bits(f)
	return uint32(x>>shiftS)&maskS == maskS && x != uvinfS && x != uvneginfS
}

func isZeroInfNaN(ix uint32) bool {
	return 2*ix-1 >= 2*0x7f800000-1
}

// IsInf reports whether f is an infinity, according to sign.
// If sign > 0, IsInf reports whether f is positive infinity.
// If sign < 0, IsInf reports whether f is negative infinity.
// If sign == 0, IsInf reports whether f is either infinity.
func IsInf(f float64, sign int) bool {
	// Test for infinity by comparing against maximum float.
	// To avoid the floating-point hardware, could use:
	x := Float64bits(f)
	return sign >= 0 && x == uvinf || sign <= 0 && x == uvneginf
}

func IsInfS(f float32, sign int) bool {
	x := Float32bits(f)
	return sign >= 0 && x == uvinfS || sign <= 0 && x == uvneginfS
}

// Copysign returns a value with the magnitude
// of x and the sign of y.
func CopysignS(x, y float32) float32 {
	const sign = 1 << 31
	return Float32frombits(Float32bits(x)&^sign | Float32bits(y)&sign)
}

// normalize returns a normal number y and exponent exp
// satisfying x == y × 2**exp. It assumes x is finite and non-zero.
func normalize(x float64) (y float64, exp int) {
	const SmallestNormal = 2.2250738585072014e-308 // 2**-1022
	if Abs(x) < SmallestNormal {
		return x * (1 << shift), -shift
	}
	return x, 0
}

// normalize returns a normal number y and exponent exp
// satisfying x == y × 2**exp. It assumes x is finite and non-zero.
func normalizeS(x float32) (y float32, exp int) {
	const SmallestNormal = 0x1p-126 // 2**-126
	if AbsS(x) < SmallestNormal {
		return x * (1 << shiftS), -shiftS
	}
	return x, 0
}
