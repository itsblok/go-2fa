package totp

func (q *QR) applyMask(mask int) {
	for y := 0; y < q.size; y++ {
		for x := 0; x < q.size; x++ {

			if q.isReserved(x, y) {
				continue
			}

			if q.matrix[y][x] == moduleBlack || q.matrix[y][x] == moduleWhite {
				if shouldFlip(mask, x, y) {
					if q.matrix[y][x] == moduleBlack {
						q.matrix[y][x] = moduleWhite
					} else {
						q.matrix[y][x] = moduleBlack
					}
				}
			}
		}
	}
}

func shouldFlip(mask, x, y int) bool {
	switch mask {
	case 0:
		return (x+y)%2 == 0
	case 1:
		return y%2 == 0
	case 2:
		return x%3 == 0
	case 3:
		return (x+y)%3 == 0
	case 4:
		return (x/3+y/2)%2 == 0
	case 5:
		return (x*y)%2+x*y%3 == 0
	case 6:
		return ((x*y)%2+(x*y)%3)%2 == 0
	case 7:
		return ((x+y)%2+(x*y)%3)%2 == 0
	default:
		return false
	}
}
