// Reference: glibc

package math

const (
	hpiInv = 0x1.45F306DC9C883p-1
	hpiinv = (Pi / 2)
)

var sincosTable0 = &sincos{
	sign:   [4]float32{1, -1, -1, 1},
	hpiInv: 0x1.45F306DC9C883p+23, // TOINTINTRINSICS=0
	hpi:    0x1.921FB54442D18p0,
	c0:     0x1p0,
	c1:     -0x1.ffffffd0c621cp-2,
	c2:     0x1.55553e1068f19p-5,
	c3:     -0x1.6c087e89a359dp-10,
	c4:     0x1.99343027bf8c3p-16,
	s1:     -0x1.555545995a603p-3,
	s2:     0x1.1107605230bc4p-7,
	s3:     -0x1.994eb3774cf24p-13,
}

var sincosTable1 = &sincos{
	sign:   [4]float32{1, -1, -1, 1},
	hpiInv: 0x1.45F306DC9C883p+23, // TOINTINTRINSICS=0
	hpi:    0x1.921FB54442D18p0,
	c0:     -0x1p0,
	c1:     0x1.ffffffd0c621cp-2,
	c2:     -0x1.55553e1068f19p-5,
	c3:     0x1.6c087e89a359dp-10,
	c4:     -0x1.99343027bf8c3p-16,
	s1:     -0x1.555545995a603p-3,
	s2:     0x1.1107605230bc4p-7,
	s3:     -0x1.994eb3774cf24p-13,
}

func CosS(y float32) float32 {
	const pio4 = Pi / 4

	x := y
	sincost := sincosTable0
	switch {
	case absToP12s(y) < absToP12s(pio4):

		x2 := x * x
		if absToP12s(y) < absToP12s(0x1p-12) {
			return 1
		}
		return sinPolyS(x, x2, sincost, 1)

	case absToP12s(y) < absToP12s(120):
		var n int
		x, n = reducePi(x, sincost)

		// Setup the signs for sin and cos.
		s := sincost.sign[n&3]
		if n&2 != 0 {
			sincost = sincosTable1
		}
		return sinPolyS(x*s, x*x, sincost, n^1)

	case absToP12s(y) < absToP12s(InfS(1)):
		// xi := Float32bits(y)
		// sign := xi >> 31
		return 0
	default:
		return InfS(1)
	}
}

type sincos struct {
	sign               [4]float32 // Sign of sine in quadrants 0..3.
	hpiInv             float32
	hpi                float32
	c0, c1, c2, c3, c4 float32 // Cosine polynomial
	s1, s2, s3         float32 // Sine polynomial
}

// Return the sine of inputs X and X2 (X squared) using the polynomial P.
// N is the quadrant, and if odd the cosine polynomial is used.
func sinPolyS(x, x2 float32, p *sincos, n int) float32 {
	if n&1 == 0 {
		x3 := x * x2
		s1 := p.s2 + x2*p.s3
		x7 := x3 * x2
		s := x + x3*p.s1
		return s + x7*s1
	}

	x4 := x2 * x2
	c2 := p.c3 + x2*p.c4
	c1 := p.c0 + x2*p.c1
	x6 := x4 * x2
	c := c1 + x4*p.c2
	return c + x6*c2
}

func reducePi(x float32, p *sincos) (float32, int) {
	// 	Use scaled float to int conversion with explicit rounding.
	// hpi_inv is prescaled by 2^24 so the quadrant ends up in bits 24..31.
	// This avoids inaccuracies introduced by truncating negative values.
	r := x * p.hpiInv
	n := int((int32(r) + 0x800000) >> 24)
	return x - float32(n)*p.hpi, n
}

func reduceLargePi(xi uint32) (float32, int) {

}
