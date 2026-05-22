package totp

import (
	"fmt"
)

const (
	qrVersion1Size = 21
	maxDataLength  = 64

	moduleEmpty = -1
	moduleWhite = 0
	moduleBlack = 1

	quietZone = 0
)

type QR struct {
	matrix [qrVersion1Size][qrVersion1Size]int
}

// GenerateQRCode builds a minimal QR matrix.
// Only supports very small payloads.
func GenerateQRCode(data string) (*QR, error) {
	if len(data) > maxDataLength {
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

// Print renders QR to terminal.
func (q *QR) Print() {
	for y := -quietZone; y < qrVersion1Size+quietZone; y++ {
		for x := -quietZone; x < qrVersion1Size+quietZone; x++ {
			if x < 0 || y < 0 || x >= qrVersion1Size || y >= qrVersion1Size {
				fmt.Print("  ")
				continue
			}

			if q.matrix[y][x] == moduleBlack {
				fmt.Print("██")
			} else {
				fmt.Print("  ")
			}
		}

		fmt.Println()
	}
}

func (q *QR) initMatrix() {
	for y := 0; y < qrVersion1Size; y++ {
		for x := 0; x < qrVersion1Size; x++ {
			q.matrix[y][x] = moduleEmpty
		}
	}
}
