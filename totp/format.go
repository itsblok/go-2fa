package totp

// applyFormatInfo writes mask + ECC metadata into fixed positions
func (q *QR) applyFormatInfo(mask int) {
	formatBits := encodeFormatBits(eccLevelL, mask)

	for i := 0; i < 6; i++ {
		q.matrix[8][i] = formatBits[i]
	}
	q.matrix[8][7] = formatBits[6]
	q.matrix[8][8] = formatBits[7]

	q.matrix[7][8] = formatBits[8]
	for i := 0; i < 6; i++ {
		q.matrix[5-i][8] = formatBits[9+i]
	}

	for i := 0; i < 7; i++ {
		q.matrix[q.size-1-i][8] = formatBits[i]
	}

	q.matrix[q.size-8][8] = moduleBlack

	for i := 0; i < 8; i++ {
		q.matrix[8][q.size-8+i] = formatBits[7+i]
	}
}

func encodeFormatBits(ecc int, mask int) []int {
	data := (ecc << 3) | mask

	// shift data left 10 bits to make room for BCH remainder
	remainder := data << 10
	g := 0b10100110111

	for i := 14; i >= 10; i-- {
		if (remainder>>i)&1 != 0 {
			remainder ^= g << (i - 10)
		}
	}

	// combine data + BCH remainder, then XOR with QR format mask
	format := (data<<10 | remainder) ^ 0x5412

	bits := make([]int, 15)
	for i := 14; i >= 0; i-- {
		bits[14-i] = (format >> i) & 1
	}

	return bits
}
