package bitmapoptimizer

import (
	"image"
	"image/color"
)

type PerColourStrategy struct{}

var _ Strategy = &PerColourStrategy{}

func (p *PerColourStrategy) Process(i image.Image) map[string][]Point {
	pixels := make(map[string][]Point)

	width := i.Bounds().Max.X
	height := i.Bounds().Max.Y
	for y := range height {
		for x := range width {
			colorname := toHex(i.At(x, y))
			rgba, ok := color.RGBAModel.Convert(i.At(x, y)).(color.RGBA)

			if !ok {
				continue
			}

			pixels[colorname] = append(pixels[colorname], Point{X: x, Y: y, Colour: rgba})
		}
	}

	return pixels
}
