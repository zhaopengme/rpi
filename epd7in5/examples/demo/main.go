package main

import (
	"fmt"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"time"

	"github.com/gandaldf/rpi/epd7in5/epd"
)

func main() {
	fmt.Printf("start\n")
	e := epd.CreateEpd()
	defer e.Close()
	defer e.Clear()
	e.Init()
	e.Clear()

	fmt.Printf("Display\n")

	dest := draw()
	draw2dimg.SaveToPngFile("test2.png", dest)

	data := getBuffer(dest)
	fmt.Println(data[:1])

	e.DisplayBlack(data)
	e.Sleep()

	fmt.Printf("sleeping\n")
	time.Sleep(5 * time.Second)
}

func draw() *image.RGBA {
	width := 800
	height := 480
	dest := image.NewRGBA(image.Rect(0, 0, width, height)) // horizontal
	gc := draw2dimg.NewGraphicContext(dest)

	black := color.RGBA{0x00, 0x00, 0x00, 0xff}
	white := color.RGBA{0xff, 0xff, 0xff, 0xff}

	gc.SetFillColor(white)
	drawRect(gc, 0, 0, float64(width), float64(height))
	gc.Fill()

	gc.SetFillColor(black)
	drawRect(gc, 0, 8, 212, 1)
	gc.Fill()
	drawRect(gc, 0, 95, 212, 1)
	gc.Fill()

	return dest
}

func drawRect(gc *draw2dimg.GraphicContext, x, y, w, h float64) {
	gc.BeginPath()
	gc.MoveTo(x, y)
	gc.LineTo(x+w, y)
	gc.LineTo(x+w, y+h)
	gc.LineTo(x, y+h)
	gc.Close()
}

func getBuffer(image *image.RGBA) []byte {
	width := 800
	height := 480

	size := (width * height) / 8
	data := make([]byte, size)
	for i := range data {
		data[i] = 255
	}

	imageWidth := image.Rect.Dx()
	imageHeight := image.Rect.Dy()

	if imageWidth == width && imageHeight == height {
		fmt.Println("Vertical")
		for y := 0; y < imageHeight; y++ {
			for x := 0; x < imageWidth; x++ {
				if isBlack(image, x, y) {
					shift := uint32(x % 8)
					data[(x+y*width)/8] &= ^(0x80 >> shift)
				}
			}
		}
	} else if imageWidth == height && imageHeight == width {
		fmt.Println("Horizontal")
		for y := 0; y < imageHeight; y++ {
			for x := 0; x < imageWidth; x++ {
				newX := y
				newY := height - x - 1
				if isBlack(image, x, y) {
					shift := uint32(y % 8)
					data[(newX+newY*width)/8] &= ^(0x80 >> shift)
				}
			}
		}
	} else {
		fmt.Println("Invalid image size")
	}
	return data
}

func isBlack(image *image.RGBA, x, y int) bool {
	r, g, b, a := getRGBA(image, x, y)
	offset := 10
	return r < 255-offset && g < 255-offset && b < 255-offset && a > offset
}

func getRGBA(image *image.RGBA, x, y int) (int, int, int, int) {
	r, g, b, a := image.At(x, y).RGBA()
	r = r / 257
	g = g / 257
	b = b / 257
	a = a / 257

	return int(r), int(g), int(b), int(a)
}
