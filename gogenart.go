package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/ajagnic/gogenart/sketch"
)

func main() {
	i := flag.Int("i", 10000, "number of iterations")
	min := flag.Uint("min", 3, "minimum number of polygon sides")
	max := flag.Uint("max", 5, "maximum number of polygon sides")
	fill := flag.Uint("fill", 100, "percent chance to fill polygon")
	color := flag.Uint("color", 0, "percent chance to randomize polygon color")
	s := flag.Float64("s", 0.1, "polygon size (percentage of width)")
	shake := flag.Float64("shake", 0.0, "amount to randomize pixel positions")
	spin := flag.Float64("spin", 0.0, "")
	w := flag.Uint("width", 0, "desired width of image")
	h := flag.Uint("height", 0, "desired height of image")
	grey := flag.Bool("grey", false, "convert to greyscale")
	invert := flag.Bool("invert", false, "invert luminance scaling")
	output := flag.String("o", "", "file to use as output")
	flag.Parse()

	if *max < *min {
		min, max = max, min
	}

	in := handleInput()
	defer in.Close()

	img, enc := sketch.Source(in)
	newImg := sketch.NewSketch(img, sketch.Params{
		Iterations:         *i,
		PolygonSidesMin:    int(*min),
		PolygonSidesMax:    int(*max),
		PolygonFillChance:  float64(*fill) / 100.0,
		PolygonColorChance: float64(*color) / 100.0,
		PolygonSizeRatio:   *s,
		PixelShake:         *shake,
		PixelSpin:          *spin,
		NewWidth:           float64(*w),
		NewHeight:          float64(*h),
		Greyscale:          *grey,
		InvertScaling:      *invert,
	}).Draw()

	out, enc := handleOutput(*output, enc)
	defer out.Close()

	sketch.Encode(out, newImg, enc)
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

func handleOutput(file, enc string) (*os.File, string) {
	if file != "" {
		out, err := os.Create(file)
		if err != nil {
			log.Fatalln(err)
		}
		fSlc := strings.Split(file, ".")
		return out, fSlc[len(fSlc)-1]
	}
	return os.Stdout, enc
}
