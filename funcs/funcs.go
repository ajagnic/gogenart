package funcs

import (
	"image/color"
	"math"
	"math/rand"
)

func RotateAround(x, y, cx, cy float64, angle int) (float64, float64) {
	rangle := rand.Intn(angle)
	theta := float64(rangle) * (math.Pi / 180)
	x1 := math.Cos(theta)*(x-cx) - math.Sin(theta)*(y-cy) + cx
	y1 := math.Sin(theta)*(x-cx) + math.Cos(theta)*(y-cy) + cy
	return x1, y1
}

func RandomChance(odds float64) bool {
	if r := rand.Intn(100) + 1; r < int(odds*100) {
		return true
	}
	return false
}

func ColorToRGB(c color.Color) (r, g, b int) {
	rr, gg, bb, _ := c.RGBA()
	r, g, b = int(rr/255), int(gg/255), int(bb/255)
	return
}

func Luminance(r, g, b int) float64 {
	values := [3]float64{float64(r), float64(g), float64(b)}
	for i, c := range values {
		c = c / 255
		if c <= 0.03928 {
			c = c / 12.92
		} else {
			c = math.Pow((c+0.055)/1.055, 2.4)
		}
		values[i] = c
	}
	return 0.2126*values[0] + 0.7152*values[1] + 0.0722*values[2]
}
