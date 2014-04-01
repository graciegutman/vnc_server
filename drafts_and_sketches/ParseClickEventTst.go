package main

import (
        "vnc/vnc"
        "fmt"
        )

func main() {
    msg := []byte{1, 2, 23, 0, 217}
    click := vnc.ParseClickEvent(msg)
    fmt.Println(click)
}