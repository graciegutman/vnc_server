package main

import (
    "fmt"
//    "os"
    "os/exec"
    "log"
)

func main() {
    
    err := screencap()    
    if err !=nil {
        fmt.Println("screencap failed")
    }
}

func screencap()(err error){
    _, err = exec.LookPath("screencapture")
    if err != nil {
        log.Fatal("screencapture not installed")
    }

    _, err = exec.Command("screencapture", "-x", "frame.png").CombinedOutput()
    return err
}
