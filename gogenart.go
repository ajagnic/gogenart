package main

import (
	"flag"
	"image/jpeg"
	"log"
	"os"

	"github.com/ajagnic/go-generative-art/sketch"
)

var (
	i   int
	min int
	max int
	s   float64
)

func main() {
	flag.IntVar(&i, "i", 10000, "number of iterations")
	flag.IntVar(&min, "min", 3, "minimum number of polygon sides")
	flag.IntVar(&max, "max", 5, "maximum number of polygon sides")
	flag.Float64Var(&s, "s", 0.1, "polygon size (percentage of width)")
	flag.Parse()

	img, err := jpeg.Decode(os.Stdin)
	if err != nil {
		log.Fatalf("decoding error: %v", err)
	}

	canvas := sketch.NewSketch(img, sketch.Params{
		Iterations:       i,
		PolygonSidesMin:  min,
		PolygonSidesMax:  max,
		PolygonSizeRatio: s,
	})
	canvas.Draw()

	jpeg.Encode(os.Stdout, canvas.Image(), nil)
}
