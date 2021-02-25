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

func main() {
	w := flag.Uint("w", 1600, "desired width of image")
	h := flag.Uint("h", 1200, "desired height of image")
	i := flag.Int("i", 10000, "number of iterations")
	min := flag.Uint("min", 3, "minimum number of polygon sides")
	max := flag.Uint("max", 5, "maximum number of polygon sides")
	fill := flag.Int("fill", 1, "1 in N chance to fill polygon")
	s := flag.Float64("s", 0.1, "polygon size (percentage of width)")
	output := flag.String("o", "stdout", "file to use as output")
	flag.Parse()

	in := handleInput()
	defer in.Close()

	img, enc, err := image.Decode(in)
	if err != nil {
		log.Fatalf("could not decode: %v\n", err)
	}

	if *max < *min {
		min, max = max, min
	}
	canvas := sketch.NewSketch(img, sketch.Params{
		Width:             int(*w),
		Height:            int(*h),
		Iterations:        *i,
		PolygonSidesMin:   int(*min),
		PolygonSidesMax:   int(*max),
		PolygonFillChance: *fill,
		PolygonSizeRatio:  *s,
	})
	canvas.Draw()

	out := handleOutput(*output, enc)
	defer out.Close()

	switch enc {
	case "png":
		png.Encode(out, canvas.Image())
	default:
		jpeg.Encode(out, canvas.Image(), nil)
	}
}

func handleInput() (in *os.File) {
	var err error
	if args := flag.Args(); len(args) > 0 {
		in, err = os.Open(args[0])
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		in = os.Stdin
	}
	return
}

func handleOutput(file, enc string) (out *os.File) {
	var err error
	if file != "stdout" {
		out, err = os.Create(file)
		if err != nil {
			if err != nil {
				log.Fatalln(err)
			}
		}
	} else {
		out = os.Stdout
	}
	return
}
