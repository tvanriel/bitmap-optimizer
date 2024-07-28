package main

import (
	"flag"
	"image"
	"io"
	"log/slog"
	"os"

	optimizer "github.com/tvanriel/bitmap-optimizer"

	_ "golang.org/x/image/bmp"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var (
	filename    string
	outname     string
	packagename string
)

func main() {
	flag.StringVar(&filename, "input-file", "-", "Image file to optimize (use - for stdin)")
	flag.StringVar(&outname, "output-file", "-", "Go file to output (use - for stdout)")
	flag.StringVar(&packagename, "package", "main", "package name for Go file")
	flag.Parse()

	var input io.Reader
	var output io.Writer
	if filename == "-" || filename == "" {
		input = os.Stdin
	} else {
		f, err := os.Open(filename)
		if err != nil {
			slog.Error("open input file", "err", err)
			return
		}
		defer f.Close()
		input = f
	}
	if outname == "-" || outname == "" {
		output = os.Stdout
	} else {
		f, err := os.OpenFile(outname, os.O_CREATE|os.O_RDWR, 0o755)
		if err != nil {
			slog.Error("open output file", "err", err)
		}
		defer f.Close()
		output = f
	}

	img, _, err := image.Decode(input)
	if err != nil {
		slog.Error("decode image", "err", err)
		return
	}
  err = optimizer.Optimize(img, output, packagename)
  if err != nil {
    slog.Error("optimize", "err", err)
  }
}
