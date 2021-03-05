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
	Iterations         int
	Width              int
	Height             int
	PolygonSidesMin    int
	PolygonSidesMax    int
	PolygonFillChance  float64
	PolygonColorChance float64
	PolygonSizeRatio   float64
	PixelShake         float64
	Greyscale          bool
}

// Sketch draws onto a destination image from a source image.
type Sketch struct {
	Params
	dc     *gg.Context
	src    image.Image
	width  float64
	height float64
	stroke float64
	shake  int
}

// NewSketch returns a blank Sketch based on the source image.
func NewSketch(source image.Image, config Params) *Sketch {
	rand.Seed(time.Now().Unix())

	max := source.Bounds().Max
	if config.Width == 0 {
		config.Width = max.X
	}
	if config.Height == 0 {
		config.Height = max.Y
	}
	w, h := config.Width, config.Height

	canvas := gg.NewContext(w, h)
	canvas.SetColor(color.Black)
	canvas.DrawRectangle(0, 0, float64(w), float64(h))
	canvas.FillPreserve()

	return &Sketch{
		Params: config,
		dc:     canvas,
		src:    source,
		width:  float64(max.X),
		height: float64(max.Y),
		stroke: config.PolygonSizeRatio * float64(w),
		shake:  int(config.PixelShake * float64(w)),
	}
}

// Draw iterates over the source image, creating the destination image.
func (s *Sketch) Draw() {
	for i := 0; i <= s.Iterations; i++ {
		rx := rand.Float64() * s.width
		ry := rand.Float64() * s.height
		r, g, b := colorToRGB(s.src.At(int(rx), int(ry)))

		l := luminance(r, g, b)
		stroke := s.stroke * l

		sides := rand.Intn((s.PolygonSidesMax - s.PolygonSidesMin) + 1)
		sides += s.PolygonSidesMin

		x := rx * float64(s.Width) / s.width
		y := ry * float64(s.Height) / s.height
		if max := s.shake; max > 0 {
			x += float64(rand.Intn(2*max) - max)
			y += float64(rand.Intn(2*max) - max)
		}

		if s.Greyscale {
			grey := int(l * 255)
			r, g, b = grey, grey, grey
		} else if l > 0.1 && randomChance(s.PolygonColorChance) {
			r, g, b = rand.Intn(256), rand.Intn(256), rand.Intn(256)
		}
		s.dc.SetRGBA255(r, g, b, rand.Intn(256))
		s.dc.DrawRegularPolygon(sides, x, y, stroke, rand.Float64())
		if randomChance(s.PolygonFillChance) {
			s.dc.FillPreserve()
		}
		s.dc.Stroke()
	}
}

// Image returns the destination image.
func (s *Sketch) Image() image.Image {
	return s.dc.Image()
}

func randomChance(odds float64) bool {
	if r := rand.Intn(100) + 1; r < int(odds*100) {
		return true
	}
	return false
}

func colorToRGB(c color.Color) (r, g, b int) {
	rr, gg, bb, _ := c.RGBA()
	r, g, b = int(rr/255), int(gg/255), int(bb/255)
	return
}

func luminance(r, g, b int) float64 {
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
