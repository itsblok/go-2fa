package totp

// encodeData converts raw string data into QR bit stream.
// Current implementation only supports byte mode encoding.
func encodeData(data string) []int {
	var bits []int

	// Byte mode
	bits = append(bits, 0, 1, 0, 0)

	length := len(data)

	for i := 7; i >= 0; i-- {
		bits = append(bits, (length>>i)&1)
	}

	for _, b := range []byte(data) {
		for i := 7; i >= 0; i-- {
			bits = append(bits, int((b>>i)&1))
		}
	}

	return bits
}
