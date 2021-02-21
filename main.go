package main

import (
	"image/jpeg"
	"log"
	"os"

	"github.com/ajagnic/go-generative-art/sketch"
)

func main() {
	img, err := jpeg.Decode(os.Stdin)
	if err != nil {
		log.Fatalf("decoding error: %v", err)
	}

	canvas := sketch.NewSketch(img, sketch.Params{
		Iterations:       10000,
		PolygonSidesMin:  3,
		PolygonSidesMax:  6,
		PolygonSizeRatio: 0.1,
	})
	canvas.Draw()

	jpeg.Encode(os.Stdout, canvas.Image(), nil)
}
