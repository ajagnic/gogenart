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

	var in *os.File
	var out *os.File
	var err error
	if args := flag.Args(); len(args) >= 1 {
		in, err = os.Open(args[0])
		if err != nil {
			log.Fatalf("file error: %v", err)
		}
		defer in.Close()
		if len(args) == 2 {
			out, err = os.Create(args[1])
			if err != nil {
				log.Fatalf("file error: %v", err)
			}
			defer out.Close()
		} else {
			out = os.Stdout
		}
	} else {
		in = os.Stdin
	}

	img, enc, err := image.Decode(in)
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
		png.Encode(out, canvas.Image())
	default:
		jpeg.Encode(out, canvas.Image(), nil)
	}
}
