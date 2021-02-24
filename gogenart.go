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
	var err error

	output := flag.String("o", "stdout", "file to use as output")
	i := flag.Int("i", 10000, "number of iterations")
	min := flag.Uint("min", 3, "minimum number of polygon sides")
	max := flag.Uint("max", 5, "maximum number of polygon sides")
	fill := flag.Int("fill", 1, "1 in N chance to fill polygon")
	s := flag.Float64("s", 0.1, "polygon size (percentage of width)")
	flag.Parse()

	in := os.Stdin
	if args := flag.Args(); len(args) > 0 {
		in, err = os.Open(args[0])
		if err != nil {
			log.Fatalln(err)
		}
		defer in.Close()
	}

	img, enc, err := image.Decode(in)
	if err != nil {
		log.Fatalf("could not decode: %v\n", err)
	}

	if *max < *min {
		min, max = max, min
	}
	canvas := sketch.NewSketch(img, sketch.Params{
		Iterations:        *i,
		PolygonSidesMin:   int(*min),
		PolygonSidesMax:   int(*max),
		PolygonFillChance: *fill,
		PolygonSizeRatio:  *s,
	})
	canvas.Draw()

	out := os.Stdout
	if f := *output; f != "stdout" {
		out, err = os.Create(f)
		if err != nil {
			log.Printf("file error: %v: saved as result.%s instead", err, enc)
			out, err = os.Create("result." + enc)
			if err != nil {
				log.Fatalln(err)
			}
		}
		defer out.Close()
	}

	switch enc {
	case "png":
		png.Encode(out, canvas.Image())
	default:
		jpeg.Encode(out, canvas.Image(), nil)
	}
}
