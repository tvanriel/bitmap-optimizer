package bitmapoptimizer

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"io"
	"strings"
)

type Strategy interface {
  Process(i image.Image) map[string][]Point
}
type Point struct {
  X, Y int
  Colour color.RGBA
}

func (p Point) GoString() string {
  return fmt.Sprintf(
    `bitmapoptimizer.Point{X: %d, Y: %d, Color: color.RGBA{R: %d, G: %d, B: %d, A: %d}}`,
    p.X, 
    p.Y, 
    p.Colour.R, 
    p.Colour.G, 
    p.Colour.B, 
    p.Colour.A,
  )
}

func Optimize(img image.Image, strategy Strategy, out io.Writer, packagename string) error {
  pixels := strategy.Process(img)

	i := 0
  header := strings.Join([]string{
    "// AUTOGENERATED FILE - DO NOT EDIT!\n",
    "package ", packagename, "\n",
    "\n",
    "import (\n",
    "\t\"github.com/tvanriel/bitmapoptimizer\"\n",
    "\t\"image/color\"\n",
    ")\n",
  }, "")
  _, err := out.Write([]byte(header))
	if err != nil {
		return fmt.Errorf("optimize: %w", err)
	}

	var mapbuf bytes.Buffer

	mapbuf.Write([]byte("\n\nvar Series = map[string]*[]bitmapoptimizer.Point{\n"))

	for name := range pixels {
		i++
		err := printcolours(i, name, pixels[name], out)
		if err != nil {
			return fmt.Errorf("optimize: %w", err)
		}
		_, err = fmt.Fprintf(&mapbuf, "\t\"#%s\": &Colour%d,\n", name, i)
		if err != nil {
			return fmt.Errorf("optimize: %w", err)
		}
	}

	mapbuf.Write([]byte("}\n"))
	_, err = io.Copy(out, &mapbuf)
	if err != nil {
		return fmt.Errorf("optimize: %w", err)
	}
	return nil
}

func toHex(c color.Color) string {
	r, g, b, _ := c.RGBA()

	return hex.EncodeToString([]byte{byte(r), byte(g), byte(b)})
}

func printcolours(i int, name string, points []Point, out io.Writer) error {
	s := make([]string, len(points))
	for p := range points {
		s[p] = points[p].GoString()
	}

  _, err := fmt.Fprintf(out, "// series: %s\n var Colour%d = []bitmapoptimizer.Point{%v}\n", name, i, strings.Join(s, ", "))
	if err != nil {
    return fmt.Errorf("printcolours: %w", err)
  } 
  return nil
}
