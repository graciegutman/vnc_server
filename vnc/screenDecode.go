package vnc

import (
	//"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"os/exec"
	//"io/ioutil"
)

func TakeScreenShot(f *os.File) (err error) {
	_, err = exec.LookPath("screencapture")
	if err != nil {
		log.Fatal("screencapture not installed")
	}
	_, err = exec.Command("screencapture", "-x", "-C", f.Name()).CombinedOutput()
	if err != nil {
		log.Fatal("screencapture failed")
	}
	//fmt.Println("took screenshot")
	return err
}

func ImgDecode(decodedPNG image.Image) (pixSlice []uint8, err error) {
	pixSlice = []uint8{}
	rect := decodedPNG.Bounds()
	rect_area := rect.Dx() * rect.Dy()

	count := 0
	for ; count < rect_area; count++ {
		x, y := findXY(count, rect)
		r, g, b, padding := decodePixel(x, y, decodedPNG)
		pixSlice = appendPixelValues(r, g, b, padding, pixSlice)
	}
	if err != nil {
		log.Fatal("could not construct pixSlice")
	}
	return pixSlice, err
}

func GetImageWidthHeight(decodedPNG image.Image) (uint16, uint16) {
	rect := decodedPNG.Bounds()
	rect_width := uint16(rect.Dx())
	rect_height := uint16(rect.Dy())
	return rect_width, rect_height
}

func DecodeFileToPNG(f *os.File) (decodedPNG image.Image, err error) {
	reader, err := os.Open(f.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	decodedPNG, err = png.Decode(reader)
	return decodedPNG, err
}

func ResizeImage(f *os.File) (err error) {
	_, err = exec.LookPath("convert")
	if err != nil {
		log.Fatal("imagemagick not installed")
	}
	var out []byte
	out, err = exec.Command("convert", f.Name(), "-resize", "40%%", f.Name()).CombinedOutput()
	if err != nil {
		log.Fatalf("resize failed. output: %s", out)
	}
	//fmt.Println("resized")
	return
}

func decodePixel(x, y int, img image.Image) (r, b, g, padding uint8) {
	pix := img.At(x, y)
	r32, g32, b32, _ := pix.RGBA()
	r8, g8, b8, padding := uint8(r32), uint8(g32), uint8(b32), uint8(0)
	return r8, g8, b8, padding
}

func appendPixelValues(r, g, b, padding uint8, pixSlice []uint8) []uint8 {
	pixSlice = append(pixSlice, r, g, b, padding)
	return pixSlice
}

func findXY(count int, rect image.Rectangle) (x, y int) {
	x, y = count%rect.Dx()+rect.Min.X, count/rect.Dx()+rect.Min.Y
	return x, y
}
