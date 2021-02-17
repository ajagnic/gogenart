package sketch

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

// Sketch draws onto a destination image from a source image.
type Sketch struct {
	dc     *gg.Context
	src    image.Image
	width  int
	height int
}

// NewSketch returns a Sketch based on the source image.
func NewSketch(source image.Image) (s *Sketch) {
	max := source.Bounds().Max
	w, h := max.X, max.Y

	canvas := gg.NewContext(w, h)
	canvas.SetColor(color.Black)
	canvas.DrawRectangle(0, 0, float64(w), float64(h))
	canvas.FillPreserve()

	s.dc = canvas
	s.src = source
	s.width, s.height = w, h
	return
}
