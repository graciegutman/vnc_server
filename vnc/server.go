package vnc

import (
    "fmt"
    "log"
    "net"
)

type ControlMessage int

const (
Start ControlMessage = iota
Stop 
)

type WorkerGroup struct {
    workerControls []chan ControlMessage
}

func Super() {
    // make TCP listener
    listener, err := CreateListener()
    checkError(err)

    // initialize the channel that each client thread
    // will use to pull screenshots from the image server
    imageOut := make(chan *FBUpdateWithImage, 50)
    workerGroup := newWorkerGroup(4, imageOut)

    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        // begin buffering images
        workerGroup.Broadcast(Start)
        // launch client thread and return to listening
        go handleClient(imageOut, conn)
    }
}

func handleClient(c chan *FBUpdateWithImage, conn net.Conn) {
    //handshake (version exchange)
    _, err := ExchangeVersions(conn)
    checkError(err)
    //fmt.Println(versionFlag)

    // send security level
    err = SendSecurity(conn)
    checkError(err)

    // init phase
    clientInitFlag, err := ReceiveClientInit(conn)
    if clientInitFlag != 1 {
        log.Fatal("Cannot give exclusive access to server")
    }

    // send serverinit msg
    pixelFormat := NewPixelFormat()
    serverInitMsg := NewServerInitMsg(pixelFormat)
    err = SendServerInit(serverInitMsg, conn)
    checkError(err)


    // main loop starts
    for {
        //read from client
        msg, msgNum, err := GetMsg(conn)
        if err != nil {
            conn.Close()
        }
        //do something with the message
        MsgDispatch(conn, msgNum, c, msg)
    }
    fmt.Printf("CLIENT exiting w/ conn %v\n", conn)
}


func newWorkerGroup(numWorkers int, imageOut chan *FBUpdateWithImage) *WorkerGroup {
    wg := &WorkerGroup{
        workerControls : make([]chan ControlMessage, numWorkers),
    }
    var control chan ControlMessage
    for i := 0; i < numWorkers; i++ {
        control = make(chan ControlMessage)
        wg.workerControls[i] = control
        go imageServer(imageOut, control)
    }
    return wg
}

func (wg *WorkerGroup) Broadcast(msg ControlMessage) {
    for _, control := range wg.workerControls {
        go func() {
            control <- msg
        }()
    }
}

func imageServer(imageOut chan *FBUpdateWithImage, control chan ControlMessage) {
    var running bool = false
    for {
        select {
            case msg := <- control:
                switch msg {
                    case Start:
                    running = true
                    case Stop:
                    running = false
                }
            default:
                if running {
                    imageOut <- NewFBUpdateWithImage()
                }
        }
    }
}
