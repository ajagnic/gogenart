package sketch

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
)

// Params represents the configuration of a sketch.
type Params struct {
	Iterations   int
	MinSides     int
	MaxSides     int
	MaxSizeRatio float64
}

// Sketch draws onto a destination image from a source image.
type Sketch struct {
	Params
	dc     *gg.Context
	src    image.Image
	width  float64
	height float64
	stroke float64
}

// NewSketch returns a blank Sketch based on the source image.
func NewSketch(source image.Image, config Params) *Sketch {
	rand.Seed(time.Now().Unix())

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
		stroke: config.MaxSizeRatio * w,
	}
	return s
}

// Draw iterates over the source image, creating the destination image.
func (s *Sketch) Draw() {
	for i := 0; i <= s.Iterations; i++ {
		rx := rand.Float64() * s.width
		ry := rand.Float64() * s.height
		r, g, b := colorToRGB(s.src.At(int(rx), int(ry)))

		l := computeLuminance(r, g, b)
		stroke := s.stroke * l

		sides := rand.Intn((s.MaxSides - s.MinSides) + 1)
		sides += s.MinSides

		s.dc.SetRGBA255(r, g, b, rand.Intn(256))
		s.dc.DrawRegularPolygon(sides, rx, ry, stroke, rand.Float64())
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

func computeLuminance(r, g, b int) float64 {
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
