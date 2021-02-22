package main

import (
	"flag"
	"image"
	"image/jpeg"
	"image/png"
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

	var f *os.File
	var err error
	switch args := flag.Args(); len(args) {
	case 1:
		f, err = os.Open(args[0])
		if err != nil {
			log.Fatalf("file error: %v", err)
		}
		defer f.Close()
	default:
		f = os.Stdin
	}

	img, enc, err := image.Decode(f)
	if err != nil {
		log.Fatalf("could not decode: %v", err)
	}

	canvas := sketch.NewSketch(img, sketch.Params{
		Iterations:       i,
		PolygonSidesMin:  min,
		PolygonSidesMax:  max,
		PolygonSizeRatio: s,
	})
	canvas.Draw()

	switch enc {
	case "png":
		png.Encode(os.Stdout, canvas.Image())
	default:
		jpeg.Encode(os.Stdout, canvas.Image(), nil)
	}
}
