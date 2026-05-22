package totp

func (q *QR) placeData(bits []int) {
	x := qrVersion1Size - 1
	y := qrVersion1Size - 1
	dir := -1

	i := 0

	for x > 0 {
		if x == 6 {
			x--
		}

		for {
			for dx := 0; dx < 2; dx++ {
				xx := x - dx

				if q.matrix[y][xx] != moduleEmpty {
					continue
				}

				if i < len(bits) {
					q.matrix[y][xx] = bits[i]
					i++
				} else {
					q.matrix[y][xx] = moduleWhite
				}
			}

			y += dir

			if y < 0 || y >= qrVersion1Size {
				y -= dir
				dir = -dir
				break
			}
		}

		x -= 2
	}
}
