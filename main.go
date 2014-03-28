package main

import (
    "vnc/vnc"
    //"encoding/binary"
    "fmt"
    "os"
    "log"
    //"bytes"
    // "io/ioutil"
)

func main() {
    //establish TCP connection
    listener, err := vnc.CreateListener()
    checkError(err)

    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        //send versions
        versionFlag, err := vnc.ExchangeVersions(conn)
        checkError(err)
        fmt.Println(versionFlag)

        //send security level
        err = vnc.SendSecurity(conn)
        checkError(err)

        //init phase
        clientInitFlag, err := vnc.ReceiveClientInit(conn)
        if clientInitFlag != 1 {
            log.Fatal("Cannot give exclusive access to server")
        }

        //send serverinit msg
        pixelFormat := vnc.NewPixelFormat()
        serverInitMsg := vnc.NewServerInitMsg(pixelFormat)
        err = vnc.SendServerInit(serverInitMsg, conn)
        checkError(err)

        c := make(chan *vnc.FrameBufferWithImage, 100)

        go vnc.NewFrameBufferWithImageRaw(c)
        go vnc.NewFrameBufferWithImageRaw(c)
        go vnc.NewFrameBufferWithImageRaw(c)
        //main loop starts
        for {
            //read from client
            //ADD A MESSAGE HANDLER
            _, msgnum := vnc.GetMsg(conn)
            if msgnum == 3 {
            vnc.SendFrameBufferRaw(conn, c)
            }
        }
    } //conn.Close()
}

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s \n", err.Error())
        os.Exit(1)
    }
}