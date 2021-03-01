# Generative Art in Go
My take on the work presented in [Generative Art in Go by Preslav Rachev](https://preslav.me/generative-art-in-golang/).

<p align="center" width="100%">
<img width="30%" src="examples/abstract.jpeg">
<img width="30%" src="examples/woman.jpeg">
<img width="30%" src="examples/lines.jpeg">
</p>

## Libraries
- [fogleman/gg](https://github.com/fogleman/gg)

## Usage
```
Usage of ./gogenart:
  -fill int
        1 in N chance to fill polygon (default 1)
  -height uint
        desired height of image (default 1200)
  -i int
        number of iterations (default 10000)
  -max uint
        maximum number of polygon sides (default 5)
  -min uint
        minimum number of polygon sides (default 3)
  -o string
        file to use as output
  -s float
        polygon size (percentage of width) (default 0.1)
  -shake float
        amount to randomize pixel positions
  -width uint
        desired width of image (default 1600)
```

The command can receive input and output in various ways.
```bash
# All of the following are equivalent:
$ cat example.jpeg | gogenart > result.jpeg

$ gogenart example.jpeg > result.jpeg

$ gogenart -o=result.jpeg example.jpeg
```

```bash
# Both JPEG and PNG files can be used.
$ gogenart example.png > result.png
```

## The Drawing Algorithm
```go
// Random pixel is selected
rx := rand.Float64() * s.width
ry := rand.Float64() * s.height
...
// Pixels luminance is used to scale the polygon
l := computeLuminance(r, g, b)
stroke := s.stroke * l
...
// Polygon is drawn with the color of that pixel
s.dc.SetRGBA255(r, g, b, rand.Intn(256))
s.dc.DrawRegularPolygon(sides, x, y, stroke, rand.Float64())
...
```

# Authors
Adrian Agnic [ [Github](https://github.com/ajagnic) ]
