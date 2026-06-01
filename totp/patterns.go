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

		if i%2 == 0 {
			val = moduleBlack
		}

		q.matrix[6][i] = val
		q.matrix[i][6] = val
	}
}

// addAlignmentPatterns places the 5×5 alignment pattern for Version 2+.
// All versions 2-4 (25-33) have a single alignment pattern whose centre
// is at (size-7, size-7).
func (q *QR) addAlignmentPatterns() {
	cx := alignmentCenter(q.size)
	if cx < 0 {
		return
	}

	for dy := -2; dy <= 2; dy++ {
		for dx := -2; dx <= 2; dx++ {
			ax := cx + dx
			ay := cx + dy

			// outer ring or centre module → black; interior ring → white
			if dx == -2 || dx == 2 || dy == -2 || dy == 2 || (dx == 0 && dy == 0) {
				q.matrix[ay][ax] = moduleBlack
			} else {
				q.matrix[ay][ax] = moduleWhite
			}
		}
	}
}

// alignmentCenter returns the single alignment-pattern centre coordinate
// (same for row and col) for a given matrix size, or -1 for Version 1.
func alignmentCenter(size int) int {
	switch size {
	case 25:
		return 18 // Version 2
	case 29:
		return 22 // Version 3
	case 33:
		return 26 // Version 4
	case 37:
		return 30 // Version 5
	case 41:
		return 34 // Version 6
	default:
		return -1 // Version 1 has no alignment pattern
	}
}
