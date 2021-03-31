package sketch

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"math/rand"
	"time"

	"github.com/fogleman/gg"

	"github.com/ajagnic/gogenart/funcs"
)

// Params represents the configuration of a sketch.
type Params struct {
	Iterations       int
	PolygonSidesMin  int
	PolygonSidesMax  int
	PolygonFill      float64 // effective range: 0.0-1.0
	PolygonColor     float64 // effective range: 0.0-1.0
	PolygonSizeRatio float64 // percentage of width
	PixelShake       float64 // percentage of width
	PixelSpin        int     // degrees of rotation
	NewWidth         float64
	NewHeight        float64
	Greyscale        bool
	InvertScaling    bool
}

// Sketch draws onto a destination image from a source image.
type Sketch struct {
	Params
	Source  image.Image
	CenterX float64
	CenterY float64
	Stroke  float64
	Shake   int
	dc      *gg.Context
	width   float64
	height  float64
}

// Source decodes a JPEG or PNG image from an input source.
// If input can't be decoded, returns a 100x100 blank image.
func Source(in io.Reader) (img image.Image, enc string) {
	img, enc, err := image.Decode(in)
	if err != nil {
		img = image.Rect(0, 0, 100, 100)
	}
	return
}

// Encode writes img to out in either JPEG or PNG format. Defaults to JPEG.
func Encode(out io.Writer, img image.Image, enc string) {
	switch enc {
	case "png":
		png.Encode(out, img)
	default:
		jpeg.Encode(out, img, nil)
	}
}

// NewSketch returns a blank Sketch based on the source image.
// Seeds the math/rand pkg.
func NewSketch(source image.Image, config Params) *Sketch {
	rand.Seed(time.Now().Unix())

	if min, max := config.PolygonSidesMin, config.PolygonSidesMax; min > max {
		config.PolygonSidesMin, config.PolygonSidesMax = max, min
	}

	max := source.Bounds().Max
	x, y := float64(max.X), float64(max.Y)
	if config.NewWidth == 0 {
		config.NewWidth = x
	}
	if config.NewHeight == 0 {
		config.NewHeight = y
	}
	w, h := config.NewWidth, config.NewHeight

	canvas := gg.NewContext(int(w), int(h))
	canvas.SetColor(color.Black)
	canvas.DrawRectangle(0, 0, w, h)
	canvas.FillPreserve()

	return &Sketch{
		Params:  config,
		Source:  source,
		CenterX: w / 2,
		CenterY: h / 2,
		Stroke:  config.PolygonSizeRatio * w,
		Shake:   int(config.PixelShake * w),
		dc:      canvas,
		width:   x,
		height:  y,
	}
}

// Draw iterates over the source image, creating the destination image.
func (s *Sketch) Draw() image.Image {
	for i := 0; i < s.Iterations; i++ {
		s.DrawOnce()
	}
	return s.Image()
}

// DrawOnce picks a random pixel, and draws a polygon at that pixels position.
// The polygons size is determinant on the pixels luminance.
// Polygon shape, size, color and position can be modified by the Params struct.
func (s *Sketch) DrawOnce() {
	rx, ry := s.Pixel()
	r, g, b := funcs.ColorToRGB(s.Source.At(int(rx), int(ry)))

	l := funcs.Luminance(r, g, b)

	stroke := s.Stroke
	if s.InvertScaling {
		stroke = 0
		if invL := math.Round(l * 100); invL != 0 {
			stroke = s.Stroke / invL
		}
	} else {
		stroke *= l
	}

	sides := rand.Intn((s.PolygonSidesMax - s.PolygonSidesMin) + 1)
	sides += s.PolygonSidesMin

	x := rx * s.NewWidth / s.width
	y := ry * s.NewHeight / s.height
	if max := s.Shake; max > 0 {
		x += float64(rand.Intn(2*max) - max)
		y += float64(rand.Intn(2*max) - max)
	}

	if s.PixelSpin > 0 {
		x, y = funcs.RotateAround(x, y, s.CenterX, s.CenterY, s.PixelSpin)
	}

	if s.Greyscale {
		grey := int(l * 255)
		r, g, b = grey, grey, grey
	} else if l > 0.1 && funcs.RandomChance(s.PolygonColor) {
		r, g, b = rand.Intn(256), rand.Intn(256), rand.Intn(256)
	}

	s.DrawAt(x, y, stroke, rand.Float64(), sides, r, g, b, rand.Intn(256))
}

// DrawAt draws a n-sided polygon at (x,y), colored with RGBA values.
func (s *Sketch) DrawAt(x, y, stroke, rotation float64, n, r, g, b, a int) {
	s.dc.SetRGBA255(r, g, b, a)
	s.dc.DrawRegularPolygon(n, x, y, stroke, rotation)
	if funcs.RandomChance(s.PolygonFill) {
		s.dc.FillPreserve()
	}
	s.dc.Stroke()
}

// Pixel returns a random point from the source images coordinate space.
func (s *Sketch) Pixel() (x, y float64) {
	x = rand.Float64() * s.width
	y = rand.Float64() * s.height
	return
}

// Image returns the destination image.
func (s *Sketch) Image() image.Image {
	return s.dc.Image()
}
