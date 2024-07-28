# Bitmap-optimizer

Bitmap optimizer is a bespoke codegen to embed a bitmap image into a binary for quick drawing individual colours in parallel.

It generates Go-source files that have a uint32-slice per colour.

Its intended usecase is for continuous drawing bitmap images to online-graffitiboards such as reddits "Place". 

The lists are unique overall, and therefore safe to draw in parallel. 

### Usage: 
```
Usage of optimize:
  -input-file string
    	Image file to optimize (use - for stdin) (default "-")
  -output-file string
    	Go file to output (use - for stdout) (default "-")
  -package string
    	package name for Go file (default "main")
```

### Example:
```go
package main

//go:generate optimizer -input-file=file.bmp -output-file=colours.go -package main

func main() {
    // Point unpacks the uint32 in X,Y coordinates where X and Y [0, 0xFF)
    x,y := Point(Colours["#000000"][0])
    draw(x,y, "#000000")
}

func draw(_, _ int, _ string) {

}

```
