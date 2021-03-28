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
)

// Params represents the configuration of a sketch.
type Params struct {
	Iterations         int
	PolygonSidesMin    int
	PolygonSidesMax    int
	PolygonFillChance  float64
	PolygonColorChance float64
	PolygonSizeRatio   float64
	PixelShake         float64
	PixelSpin          int
	NewWidth           float64
	NewHeight          float64
	Greyscale          bool
	InvertScaling      bool
}

// Sketch draws onto a destination image from a source image.
type Sketch struct {
	Params
	dc      *gg.Context
	src     image.Image
	width   float64
	height  float64
	centerX float64
	centerY float64
	stroke  float64
	shake   int
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

// NewSketch returns a blank Sketch based on the source image.
// Seeds the math/rand pkg.
func NewSketch(source image.Image, config Params) *Sketch {
	rand.Seed(time.Now().Unix())

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
		dc:      canvas,
		src:     source,
		width:   x,
		height:  y,
		centerX: w / 2,
		centerY: h / 2,
		stroke:  config.PolygonSizeRatio * w,
		shake:   int(config.PixelShake * w),
	}
}

// Draw iterates over the source image, creating the destination image.
func (s *Sketch) Draw() image.Image {
	for i := 0; i < s.Iterations; i++ {
		rx := rand.Float64() * s.width
		ry := rand.Float64() * s.height
		r, g, b := colorToRGB(s.src.At(int(rx), int(ry)))

		l := luminance(r, g, b)

		stroke := s.stroke
		if s.InvertScaling {
			stroke = 0
			if invL := math.Round(l * 100); invL != 0 {
				stroke = s.stroke / invL
			}
		} else {
			stroke *= l
		}

		sides := rand.Intn((s.PolygonSidesMax - s.PolygonSidesMin) + 1)
		sides += s.PolygonSidesMin

		x := rx * s.NewWidth / s.width
		y := ry * s.NewHeight / s.height
		if max := s.shake; max > 0 {
			x += float64(rand.Intn(2*max) - max)
			y += float64(rand.Intn(2*max) - max)
		}

		if s.PixelSpin > 0 {
			x, y = rotateAround(x, y, s.centerX, s.centerY, s.PixelSpin)
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
	return s.dc.Image()
}

// Image returns the destination image.
func (s *Sketch) Image() image.Image {
	return s.dc.Image()
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

func rotateAround(x, y, cx, cy float64, angle int) (float64, float64) {
	rangle := rand.Intn(angle)
	theta := float64(rangle) * (math.Pi / 180)
	x1 := math.Cos(theta)*(x-cx) - math.Sin(theta)*(y-cy) + cx
	y1 := math.Sin(theta)*(x-cx) + math.Cos(theta)*(y-cy) + cy
	return x1, y1
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
