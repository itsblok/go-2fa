package totp

// score evaluates QR quality (lower = better)
func (q *QR) score() int {
	return scoreAdjacent(q.matrix) + scoreBlocks(q.matrix)
}

func scoreAdjacent(m [][]int) int {
	penalty := 0

	for y := 0; y < len(m); y++ {
		for x := 0; x < len(m)-1; x++ {

			if !isFunctional(x, y, len(m)) {
				continue
			}

			if m[y][x] == m[y][x+1] {
				penalty++
			}
		}
	}

	for x := 0; x < len(m); x++ {
		for y := 0; y < len(m)-1; y++ {

			if !isFunctional(x, y, len(m)) {
				continue
			}

			if m[y][x] == m[y+1][x] {
				penalty++
			}
		}
	}

	return penalty
}

func scoreBlocks(m [][]int) int {
	penalty := 0

	for y := 0; y < len(m)-1; y++ {
		for x := 0; x < len(m)-1; x++ {

			if !isFunctional(x, y, len(m)) {
				continue
			}

			v := m[y][x]

			if m[y][x+1] == v &&
				m[y+1][x] == v &&
				m[y+1][x+1] == v {
				penalty += 3
			}
		}
	}

	return penalty
}

func isFunctional(x, y, size int) bool {
	if (x < 9 && y < 9) ||
		(x > size-9 && y < 9) ||
		(x < 9 && y > size-9) {
		return false
	}

	if x == 6 || y == 6 {
		return false
	}

	cx := alignmentCenter(size)
	if cx >= 0 && abs(x-cx) <= 2 && abs(y-cx) <= 2 {
		return false
	}

	return true
}
