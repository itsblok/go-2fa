package totp

import (
	"fmt"
)

const (
	// Real byte-mode capacity for ECC Level-L
	// Formula: (dataCW*8 - 4 - 8 - 4) / 8
	maxDataLengthV1 = 17  // Version 1-L: 19 data CW
	maxDataLengthV2 = 32  // Version 2-L: 34 data CW
	maxDataLengthV3 = 53  // Version 3-L: 55 data CW
	maxDataLengthV4 = 78  // Version 4-L: 80 data CW
	maxDataLengthV5 = 106 // Version 5-L: 108 data CW
	maxDataLengthV6 = 134 // Version 6-L: 136 data CW

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
	maxDataBytes int // maximum bytes encodable in byte mode (ECC-L)
	dataCW       int // data codewords
	ecCW         int // error-correction codewords
}

var (
	version1 = QRVersion{21, maxDataLengthV1, 19, 7}
	version2 = QRVersion{25, maxDataLengthV2, 34, 10}
	version3 = QRVersion{29, maxDataLengthV3, 55, 15}
	version4 = QRVersion{33, maxDataLengthV4, 80, 20}
	version5 = QRVersion{37, maxDataLengthV5, 108, 26}
	version6 = QRVersion{41, maxDataLengthV6, 136, 36}
)

func selectVersion(dataLen int) (QRVersion, error) {
	switch {
	case dataLen <= maxDataLengthV1:
		return version1, nil
	case dataLen <= maxDataLengthV2:
		return version2, nil
	case dataLen <= maxDataLengthV3:
		return version3, nil
	case dataLen <= maxDataLengthV4:
		return version4, nil
	case dataLen <= maxDataLengthV5:
		return version5, nil
	case dataLen <= maxDataLengthV6:
		return version6, nil
	default:
		return QRVersion{}, fmt.Errorf("data too large: %d bytes (max %d for supported versions)", dataLen, maxDataLengthV6)
	}
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

	version, err := selectVersion(len(dataBytes))
	if err != nil {
		return nil, err
	}

	qr := newQR(version.size)

	qr.initMatrix()
	qr.addFinderPatterns()
	qr.addTimingPatterns()
	qr.addAlignmentPatterns() // required for Version 2+

	bits := encodeData(data)
	final := applyECC(bits, version.dataCW, version.ecCW)

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

func (q *QR) Size() int {
	return q.size
}

func (q *QR) Mask() int {
	return q.mask
}

// CountCells returns counts of empty, black, and white cells (for debugging)
func CountCells(q *QR) (empty, black, white int) {
	for y := 0; y < q.size; y++ {
		for x := 0; x < q.size; x++ {
			switch q.matrix[y][x] {
			case moduleEmpty:
				empty++
			case moduleBlack:
				black++
			case moduleWhite:
				white++
			}
		}
	}
	return
}
