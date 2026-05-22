package totp

func (q *QR) addFinderPatterns() {
	addFinder := func(x, y int) {
		for dy := -1; dy <= 7; dy++ {
			for dx := -1; dx <= 7; dx++ {
				xx := x + dx
				yy := y + dy

				if xx < 0 || yy < 0 || xx >= q.size || yy >= q.size {
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
	addFinder(q.size-7, 0)
	addFinder(0, q.size-7)
}

func (q *QR) addTimingPatterns() {
	for i := 8; i < q.size-8; i++ {
		val := moduleWhite

		if i%2 == 1 {
			val = moduleBlack
		}

		q.matrix[6][i] = val
		q.matrix[i][6] = val
	}
}
