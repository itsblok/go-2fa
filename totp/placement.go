package totp

func (q *QR) placeData(bytes []int) {
	x := q.size - 1
	y := q.size - 1
	dir := -1

	i := 0

	for x > 0 {
		if x == 6 {
			x--
		}

		for {
			for dx := 0; dx < 2; dx++ {
				xx := x - dx

				if q.isReserved(xx, y) {
					continue
				}

				if q.matrix[y][xx] != moduleEmpty {
					continue
				}

				var bit int

				if i < len(bytes)*8 {
					byteIndex := i / 8
					bitIndex := 7 - (i % 8)

					bit = (bytes[byteIndex] >> bitIndex) & 1
					i++
				} else {
					bit = 0
				}

				if bit == 1 {
					q.matrix[y][xx] = moduleBlack
				} else {
					q.matrix[y][xx] = moduleWhite
				}
			}

			y += dir

			if y < 0 || y >= q.size {
				y -= dir
				dir = -dir
				break
			}
		}

		x -= 2
	}
}

func (q *QR) clone() *QR {
	newMatrix := make([][]int, q.size)

	for i := range q.matrix {
		newMatrix[i] = append([]int(nil), q.matrix[i]...)
	}

	return &QR{
		matrix: newMatrix,
		size:   q.size,
		mask:   q.mask,
	}
}

func (q *QR) isReserved(x, y int) bool {
	// finder patterns (top-left, top-right, bottom-left) + separators
	if (x < 9 && y < 9) ||
		(x > q.size-9 && y < 9) ||
		(x < 9 && y > q.size-9) {
		return true
	}

	// timing patterns
	if x == 6 || y == 6 {
		return true
	}

	// alignment pattern (Version 2+): 5×5 block centered at alignmentCenter
	cx := alignmentCenter(q.size)
	if cx >= 0 && abs(x-cx) <= 2 && abs(y-cx) <= 2 {
		return true
	}

	return false
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
