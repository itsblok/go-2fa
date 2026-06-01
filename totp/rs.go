package totp

func rsGeneratorPoly(degree int) []int {
	g := []int{1}

	for i := 0; i < degree; i++ {
		g = polyMul(g, []int{1, gfExp[i]})
	}

	return g
}

func polyMul(a, b []int) []int {
	result := make([]int, len(a)+len(b)-1)

	for i := range a {
		for j := range b {
			result[i+j] ^= gfMul(a[i], b[j])
		}
	}

	return result
}

func rsComputeECC(data []int, ecLen int) []int {
	ec := make([]int, ecLen)

	gen := rsGeneratorPoly(ecLen)

	for _, b := range data {
		factor := b ^ ec[0]

		copy(ec, ec[1:])
		ec[len(ec)-1] = 0

		for i := 0; i < ecLen; i++ {
			// fix: gen has ecLen+1 coefficients; index 0 is the leading 1
			// which is implicit, so we multiply against gen[i+1]
			ec[i] ^= gfMul(gen[i+1], factor)
		}
	}

	return ec
}

// applyECC pads data to dataCapacity codewords, computes ecLen ECC bytes,
// and returns the full interleaved codeword slice.
func applyECC(data []int, dataCapacity int, ecLen int) []int {
	bytes := bitsToBytes(data)

	// pad to exact data codeword capacity with alternating 0xEC / 0x11
	for i := 0; len(bytes) < dataCapacity; i++ {
		if i%2 == 0 {
			bytes = append(bytes, 0xEC)
		} else {
			bytes = append(bytes, 0x11)
		}
	}

	ecc := rsComputeECC(bytes, ecLen)

	result := make([]int, 0, len(bytes)+len(ecc))
	result = append(result, bytes...)
	result = append(result, ecc...)

	return result
}
