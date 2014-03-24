package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

func main() {
	reader, err := os.Open("frame.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	decoded_png, err := png.Decode(reader)
	//loops through pixels in decoded_png

	pixArray := []uint8{}

	rect := decoded_png.Bounds()
	rect_area := rect.Dx() * rect.Dy()
	count := 0
	for ; count < rect_area; count++ {
		x, y := findXY(count, rect)
		r, g, b, padding := decodePixel(x, y, decoded_png)
		pixArray = appendPixelValues(r, g, b, padding, pixArray)
	}
	fmt.Println(pixArray)
}

func decodePixel(x, y int, img image.Image) (r, g, b, padding uint8) {
	pix := img.At(x, y)
	r32, g32, b32, _ := pix.RGBA()
	r8, g8, b8, padding := uint8(r32), uint8(g32), uint8(b32), uint8(0)
	return r8, g8, b8, padding
}

func appendPixelValues(r, g, b, padding uint8, pixArray []uint8) []uint8 {
	pixArray = append(pixArray, padding, r, g, b)
	return pixArray
}

func findXY(count int, rect image.Rectangle) (x, y int) {
	x, y = count%rect.Dx()+rect.Min.X, count%rect.Dx()+rect.Min.Y
	return x, y
}
