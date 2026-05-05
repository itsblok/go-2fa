package totp

import (
	"fmt"
)

const qrSize = 21

type QR struct {
	matrix [qrSize][qrSize]int
}

// GenerateQRCode builds a minimal QR matrix.
// Only supports very small payloads.
func GenerateQRCode(data string) (*QR, error) {
	if len(data) > 64 {
		// TODO: Support larger QR versions
		// - Select QR version dynamically based on data length
		// - Update matrix size (e.g. 25x25, 29x29, ...)
		// - Adjust alignment patterns for higher versions
		return nil, fmt.Errorf("data too long for current QR version")
	}

	qr := &QR{}
	qr.initMatrix()
	qr.addFinderPatterns()
	qr.addTimingPatterns()

	bits := encodeData(data)

	// Right now we are placing raw data bits directly into the matrix.
	// This makes the QR fragile and difficult to scan.
	// TODO: Add Reed-Solomon error correction
	// - Split data into codewords (8-bit blocks)
	// - Generate error correction codewords using Reed-Solomon
	// - Interleave data and EC codewords as per QR spec
	// - Append final bit stream before placement

	qr.placeData(bits)

	// TODO: Apply masking patterns
	// - Apply all mask patterns (0–7)
	// - Score each result using penalty rules
	// - Select the best one

	return qr, nil
}

// Print renders QR to terminal
func (q *QR) Print() {
	for y := 0; y < qrSize; y++ {
		for x := 0; x < qrSize; x++ {
			if q.matrix[y][x] == 1 {
				fmt.Print("██")
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Println()
	}
}

func (q *QR) initMatrix() {
	for y := 0; y < qrSize; y++ {
		for x := 0; x < qrSize; x++ {
			q.matrix[y][x] = -1
		}
	}
}

func (q *QR) addFinderPatterns() {
	addFinder := func(x, y int) {
		for dy := -1; dy <= 7; dy++ {
			for dx := -1; dx <= 7; dx++ {
				xx := x + dx
				yy := y + dy

				if xx < 0 || yy < 0 || xx >= qrSize || yy >= qrSize {
					continue
				}

				if dx >= 0 && dx <= 6 && dy >= 0 && dy <= 6 &&
					(dx == 0 || dx == 6 || dy == 0 || dy == 6 ||
						(dx >= 2 && dx <= 4 && dy >= 2 && dy <= 4)) {
					q.matrix[yy][xx] = 1
				} else {
					q.matrix[yy][xx] = 0
				}
			}
		}
	}

	addFinder(0, 0)
	addFinder(qrSize-7, 0)
	addFinder(0, qrSize-7)
}

func (q *QR) addTimingPatterns() {
	for i := 8; i < qrSize-8; i++ {
		val := i % 2
		q.matrix[6][i] = val
		q.matrix[i][6] = val
	}
}

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

func (q *QR) placeData(bits []int) {
	x := qrSize - 1
	y := qrSize - 1
	dir := -1

	i := 0

	for x > 0 {
		if x == 6 {
			x--
		}

		for {
			for dx := 0; dx < 2; dx++ {
				xx := x - dx
				if q.matrix[y][xx] != -1 {
					continue
				}

				if i < len(bits) {
					q.matrix[y][xx] = bits[i]
					i++
				} else {
					q.matrix[y][xx] = 0
				}
			}

			y += dir
			if y < 0 || y >= qrSize {
				y -= dir
				dir = -dir
				break
			}
		}

		x -= 2
	}
}
