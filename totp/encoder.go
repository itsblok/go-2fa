package totp

// encodeData converts raw string data into QR bit stream.
// Current implementation only supports byte mode encoding.
func encodeData(data string) []int {
	var bits []int

	// byte mode indicator
	bits = append(bits, 0, 1, 0, 0)

	// length (8 bit simplified)
	length := len(data)
	for i := 7; i >= 0; i-- {
		bits = append(bits, (length>>i)&1)
	}

	// data bytes
	for _, b := range []byte(data) {
		for i := 7; i >= 0; i-- {
			bits = append(bits, int((b>>i)&1))
		}
	}

	// terminator (QR spec: up to 4 zeros)
	for i := 0; i < 4; i++ {
		bits = append(bits, 0)
	}

	// pad to byte alignment
	for len(bits)%8 != 0 {
		bits = append(bits, 0)
	}

	return bits
}

func bitsToBytes(bits []int) []int {
	var out []int

	for i := 0; i < len(bits); i += 8 {
		val := 0

		for j := 0; j < 8 && i+j < len(bits); j++ {
			val = (val << 1) | bits[i+j]
		}

		out = append(out, val)
	}

	return out
}
