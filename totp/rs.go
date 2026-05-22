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
			ec[i] ^= gfMul(gen[i], factor)
		}
	}

	return ec
}

func applyECC(data []int) []int {
	bytes := bitsToBytes(data)

	dataBytes := bytes

	ecLen := 10 // QR-L simplified

	ecc := rsComputeECC(dataBytes, ecLen)

	result := make([]int, 0, len(dataBytes)+len(ecc))

	result = append(result, dataBytes...)
	result = append(result, ecc...)

	return result
}
