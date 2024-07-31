package bitmapoptimizer

import (
	"image"
	"image/color"
	"math"
	"slices"
	"strconv"

	"github.com/anthonynsimon/bild/effect"
)

type EdgeDetectStrategy struct {
	Buckets int
}

var _ Strategy = &EdgeDetectStrategy{}

type Coordinate struct{ X, Y int }
type EdgeDetected struct {
	ED    uint8
	Point Point
}

func (e *EdgeDetectStrategy) Process(img image.Image) map[string][]Point {

	max := img.Bounds().Max
	merged := make([]EdgeDetected, max.X*max.Y)
	edgeDetect := effect.EdgeDetection(img, 1.0)
	i := 0
	for y := range max.Y {
		for x := range max.X {
			c := img.At(x, y)

			edC := color.GrayModel.Convert(edgeDetect.At(x, y)).(color.Gray)
			rgba := color.RGBAModel.Convert(c).(color.RGBA)
			p := Point{X: x, Y: y, Colour: rgba}
			merged[i] = EdgeDetected{ED: edC.Y, Point: p}
			i++
		}
	}

	slices.SortFunc(merged, func(i, j EdgeDetected) int {
		return int(i.ED) - int(j.ED)
	})
	buckets := exponentialBuckets(merged, e.Buckets)
	final := make(map[string][]Point)

	for b := range buckets {
		points := []Point{}
		for p := range buckets[b] {
			points = append(points, buckets[b][p].Point)
		}
		final[strconv.Itoa(b)] = points
	}

	return final
}

func exponentialBuckets(list []EdgeDetected, numBuckets int) [][]EdgeDetected {
	if len(list) == 0 || numBuckets <= 0 {
		return [][]EdgeDetected{}
	}

	// Calculate the base for exponential growth
	base := math.Pow(float64(len(list)), 1/float64(numBuckets))

	// Calculate bucket sizes
	sizes := make([]int, numBuckets)
	totalSize := 0
	for i := 0; i < numBuckets; i++ {
		sizes[i] = int(math.Round(math.Pow(base, float64(i))))
		totalSize += sizes[i]
	}

	// Adjust bucket sizes to fit within list length
	excess := totalSize - len(list)
	if excess > 0 {
		for i := numBuckets - 1; i >= 0 && excess > 0; i-- {
			if sizes[i] > excess {
				sizes[i] -= excess
				excess = 0
			} else {
				excess -= sizes[i]
				sizes[i] = 0
			}
		}
	}

	// Fill buckets
	var buckets [][]EdgeDetected
	start := 0
	for i := 0; i < numBuckets; i++ {
		end := min(start+sizes[i], len(list))
		if end > start {
			buckets = append(buckets, list[start:end])
		}
		start = end
	}

	return buckets
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
