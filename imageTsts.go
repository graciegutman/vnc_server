package main 

import (
        "image"
        //"vnc/vnc"
        //"io/ioutil"
        //"time"
        "fmt"
        "vnc/vnc"
        //"os"
        )

// in a loop

//      get X, Y value

//      get color values at X, Y 

//      find deltas at each color value

//      if deltas > or < 0, X becomes left bound, Y becomes upper and lower bound

//      continue to loop

//      if deltas > or < 0, and this X < old RX, RX = thisX if X > old LX, LX = thisX
//      if deltas > or < 0, and this Y > old LY, move LY = thisY

var threshold int = 5000

type Point struct{
    x, y int
}
type RectBounds struct{
    leftb, upperb, rightb, lowerb Point
}

func (rb *RectBounds) setRight (x, y int) {
    rb.rightb.x = x
    rb.rightb.y = y
}

func (rb *RectBounds) setLeft (x, y int) {
    rb.leftb.x = x
    rb.leftb.y = y
}

func (rb *RectBounds) setUpper (x, y int) {
    rb.upperb.x = x
    rb.upperb.y = y
}

func (rb *RectBounds) setLower (x, y int) {
    rb.lowerb.x = x
    rb.lowerb.y = y
}

func abs(x uint32) uint32 {
    switch {
        case x < 0:
            return -x
        case x == 0:
            return 0 // return correctly abs(-0)
        }
    return x
}

func pixDifference(x, y, threshold int, img1, img2 image.Image)(bool) {
    t := uint32(threshold)

    r1, g1, b1, a1 := img1.At(x, y).RGBA()
    r2, g2, b2, a2 := img2.At(x, y).RGBA()
    if abs(r2 - r1) > t || abs(g2 - g1) > t || abs(b2 - b1) > t || abs(a2 - a1) > t {
        return true
    }
    return false
}

func findFirstPixel(rectArea int, img1, img2 image.Image, rect image.Rectangle, threshold int) *RectBounds {
    count := 0
    for ; count < rectArea; count++ {
        x, y := vnc.FindXY(count, rect)
        if pixDifference(x, y, threshold, img1, img2) {
            fmt.Println("in loop iteration, %d", count)
            rb := &RectBounds{}

            rb.setLeft(x, y)
            rb.setRight(x, y)
            rb.setUpper(x, y)
            rb.setLower(x, y)
            return rb
        }
    }
    return &RectBounds{}
}

func findRectBounds(rectArea int, img1, img2 image.Image, rb *RectBounds, rect image.Rectangle, threshold int) {
    totalChanges := 0
    count := 0
    for ; count < rectArea; count++ {
        x, y := vnc.FindXY(count, rect)
        if pixDifference(x, y, threshold, img1, img2) {
            totalChanges += 1
            if x < rb.leftb.x {
                fmt.Println("leftbound changed to", x, y)
                rb.setLeft(x, y)
            }else if x > rb.rightb.x {
                fmt.Println("rightbound changed to", x, y)
                rb.setRight(x, y)
            }

            if y > rb.lowerb.y {
                fmt.Println("lowerbound changed to", x, y)
                rb.setLower(x, y)
            }
        }
    }
    fmt.Println(totalChanges)
}

func main() {

// decode into image.images

    PNG1, _ := vnc.DecodeFileToPNGtst("screen4.png")
    PNG2, _ := vnc.DecodeFileToPNGtst("screen4.png")

// prelim rectangle bounds are X,Y value of first pixel that differs.
    PNGRect := PNG1.Bounds()
    rectArea := PNGRect.Dx() * PNGRect.Dy()

// establish rect RX, LX, UY, LY
    rb := findFirstPixel(rectArea, PNG1, PNG2, PNGRect, threshold)

    findRectBounds(rectArea, PNG1, PNG2, rb, PNGRect, threshold)
    fmt.Println(rb.leftb.x, rb.rightb.x, rb.upperb.y, rb.lowerb.y)
}


