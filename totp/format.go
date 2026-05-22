package totp

// applyFormatInfo writes mask + ECC metadata into fixed positions
func (q *QR) applyFormatInfo(mask int) {
	formatBits := encodeFormatBits(eccLevelL, mask)

	for i := 0; i < 8; i++ {
		q.matrix[8][i] = formatBits[i]
		q.matrix[i][8] = formatBits[i]
	}

	for i := 0; i < 7; i++ {
		q.matrix[q.size-1-i][8] = formatBits[i]
	}
}

func encodeFormatBits(ecc int, mask int) []int {
	// QR format: 5-bit ECC+mask + 10-bit BCH

	format := (ecc << 3) | mask

	// QR standard BCH mask (simplified but valid structure)
	g := 0b10100110111

	for i := 14; i >= 10; i-- {
		if (format & (1 << i)) != 0 {
			format ^= g << (i - 10)
		}
	}

	bits := make([]int, 15)

	for i := 14; i >= 0; i-- {
		bits[14-i] = (format >> i) & 1
	}

	return bits
}
