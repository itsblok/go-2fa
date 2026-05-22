package totp

func (q *QR) addFinderPatterns() {
	addFinder := func(x, y int) {
		for dy := -1; dy <= 7; dy++ {
			for dx := -1; dx <= 7; dx++ {
				xx := x + dx
				yy := y + dy

				if xx < 0 || yy < 0 || xx >= qrVersion1Size || yy >= qrVersion1Size {
					continue
				}

				if dx >= 0 && dx <= 6 && dy >= 0 && dy <= 6 &&
					(dx == 0 || dx == 6 || dy == 0 || dy == 6 ||
						(dx >= 2 && dx <= 4 && dy >= 2 && dy <= 4)) {
					q.matrix[yy][xx] = moduleBlack
				} else {
					q.matrix[yy][xx] = moduleWhite
				}
			}
		}
	}

	addFinder(0, 0)
	addFinder(qrVersion1Size-7, 0)
	addFinder(0, qrVersion1Size-7)
}

func (q *QR) addTimingPatterns() {
	for i := 8; i < qrVersion1Size-8; i++ {
		val := moduleWhite

		if i%2 == 1 {
			val = moduleBlack
		}

		q.matrix[6][i] = val
		q.matrix[i][6] = val
	}
}
