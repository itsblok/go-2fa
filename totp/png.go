package totp

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func (q *QR) SavePNG(filename string, scale int) error {
	size := q.size

	img := image.NewRGBA(image.Rect(0, 0, size*scale, size*scale))

	black := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {

			var c color.Color
			if q.matrix[y][x] == moduleBlack {
				c = black
			} else {
				c = white
			}

			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					img.Set(x*scale+dx, y*scale+dy, c)
				}
			}
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
