package totp

var gfExp [512]int
var gfLog [256]int

func init() {
	x := 1

	for i := 0; i < 255; i++ {
		gfExp[i] = x
		gfLog[x] = i

		x <<= 1
		if x&0x100 != 0 {
			x ^= 0x11d // QR polynomial
		}
	}

	for i := 255; i < 512; i++ {
		gfExp[i] = gfExp[i-255]
	}
}

func gfMul(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return gfExp[gfLog[a]+gfLog[b]]
}
