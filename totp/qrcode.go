package totp

import (
	"fmt"
)

const (
	maxDataLengthV1 = 64
	maxDataLengthV2 = 128

	moduleEmpty = -1
	moduleWhite = 0
	moduleBlack = 1

	quietZone = 4

	eccLevelL = 1
)

type QR struct {
	matrix [][]int
	size   int
	mask   int
}

type QRVersion struct {
	size         int
	dataCapacity int
}

var version1 = QRVersion{21, 64}
var version2 = QRVersion{25, 128}

func selectVersion(len int) QRVersion {
	if len <= maxDataLengthV1 {
		return version1
	}
	return version2
}

func newQR(size int) *QR {
	matrix := make([][]int, size)

	for i := range matrix {
		matrix[i] = make([]int, size)
	}

	return &QR{
		matrix: matrix,
		size:   size,
	}
}

func GenerateQRCode(data string) (*QR, error) {
	dataBytes := []byte(data)

	version := selectVersion(len(dataBytes))

	if !version.fits(dataBytes) {
		return nil, fmt.Errorf("data too large for QR version")
	}

	qr := newQR(version.size)

	qr.initMatrix()
	qr.addFinderPatterns()
	qr.addTimingPatterns()

	bits := encodeData(data)

	// ECC now works on BYTE STREAM
	final := applyECC(bits)

	qr.placeData(final)

	best := qr.selectBestMask()

	best.applyFormatInfo(best.mask)

	return best, nil
}

func (q *QR) selectBestMask() *QR {
	var best *QR
	bestScore := int(^uint(0) >> 1)

	for mask := 0; mask < 8; mask++ {
		candidate := q.clone()
		candidate.mask = mask

		candidate.applyMask(mask)

		// format info must be applied before scoring
		candidate.applyFormatInfo(mask)

		score := candidate.score()

		if score < bestScore {
			bestScore = score
			best = candidate
		}
	}

	return best
}

func (q *QR) initMatrix() {
	for y := 0; y < q.size; y++ {
		for x := 0; x < q.size; x++ {
			q.matrix[y][x] = moduleEmpty
		}
	}
}

func (q *QR) Print() {
	for y := -quietZone; y < q.size+quietZone; y++ {
		for x := -quietZone; x < q.size+quietZone; x++ {
			if x < 0 || y < 0 || x >= q.size || y >= q.size {
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

func (v QRVersion) fits(data []byte) bool {
	return len(data) <= v.dataCapacity
}

func (q *QR) Size() int {
	return q.size
}

func (q *QR) Mask() int {
	return q.mask
}
