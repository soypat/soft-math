// Taken from the Go standard Library. The Go Authors.

package math

import "unsafe"

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
	uvnans    uint32 = 0x7fc00000
	uvinfs    uint32 = 0x7f800000
	uvneginfs uint32 = 0xff800000
	uvones    uint32 = 0x3f800000
)

// Abs returns the absolute value of x.
//
// Special cases are:
//	Abs(±Inf) = +Inf
//	Abs(NaN) = NaN
func Abs(x float64) float64 {
	return Float64frombits(Float64bits(x) &^ (1 << 63))
}

func Abss(x float32) float32 {
	return Float32frombits(Float32bits(x) &^ (1 << 31))
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

// NaN returns an IEEE 754 ``not-a-number'' value.
func NaN() float64 { return Float64frombits(uvnan) }

// IsNaN reports whether f is an IEEE 754 ``not-a-number'' value.
func IsNaN(f float64) (is bool) {
	// IEEE 754 says that only NaNs satisfy f != f.
	// To avoid the floating-point hardware, could use:
	//	x := Float64bits(f);
	//	return uint32(x>>shift)&mask == mask && x != uvinf && x != uvneginf
	return f != f
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

func IsInfs(f float32, sign int) bool {
	x := Float32bits(f)
	return sign >= 0 && x == uvinfs || sign <= 0 && x == uvneginfs
}

// normalize returns a normal number y and exponent exp
// satisfying x == y × 2**exp. It assumes x is finite and non-zero.
func normalize(x float64) (y float64, exp int) {
	const SmallestNormal = 2.2250738585072014e-308 // 2**-1022
	if Abs(x) < SmallestNormal {
		return x * (1 << 52), -52
	}
	return x, 0
}

func i32FromBits(b uint32) int32 {
	return *(*int32)(unsafe.Pointer(&b))
}
