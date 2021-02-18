package sketch

import (
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
)

// Params represents the configuration of a sketch.
type Params struct {
	Iterations   int
	PolygonEdges int
}

// Sketch draws onto a destination image from a source image.
type Sketch struct {
	Params
	dc     *gg.Context
	src    image.Image
	width  float64
	height float64
}

func init() {
	rand.Seed(time.Now().Unix())
}

// NewSketch returns a blank Sketch based on the source image.
func NewSketch(source image.Image, config Params) *Sketch {
	max := source.Bounds().Max
	w, h := float64(max.X), float64(max.Y)

	canvas := gg.NewContext(max.X, max.Y)
	canvas.SetColor(color.Black)
	canvas.DrawRectangle(0, 0, w, h)
	canvas.FillPreserve()

	s := &Sketch{
		Params: config,
		dc:     canvas,
		src:    source,
		width:  w,
		height: h,
	}
	return s
}

// Draw iterates over the source image, creating the destination image.
func (s *Sketch) Draw() {
	for i := 0; i <= s.Iterations; i++ {
		rx := rand.Float64() * s.width
		ry := rand.Float64() * s.height
		r, g, b := colorToRGB(s.src.At(int(rx), int(ry)))

		strokeRatio := 0.01 * s.width

		s.dc.SetRGBA255(r, g, b, rand.Intn(255))
		s.dc.DrawRegularPolygon(s.PolygonEdges, rx, ry, strokeRatio, 0)
		s.dc.FillPreserve()
		s.dc.Stroke()
	}
}

func (s *Sketch) SaveImage(file string) {
	s.dc.SavePNG(file)
}

func colorToRGB(c color.Color) (r, g, b int) {
	rr, gg, bb, _ := c.RGBA()
	r, g, b = int(rr/255), int(gg/255), int(bb/255)
	return
}
