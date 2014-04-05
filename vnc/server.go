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
    // For controlling worker threads (start, stop producing images)
    workerControls []chan ControlMessage
    // For notifying worker threads of new clients to send images to
    workerImageChanChans []chan chan *FBUpdateWithImage
    // For notifying worker threads of closed clients whose image
    // Channels need to be closed
    workerErrChans []chan chan *FBUpdateWithImage
}

type ServerClientImageChans struct {
    // Per-client channels for a worker to send images on
    imageChannels []chan *FBUpdateWithImage
}

func cleanUpCrew(wg *WorkerGroup, alertChan chan chan *FBUpdateWithImage) {
    for {
        select {
            case deadClient := <- alertChan:
                wg.BroadcastChanToClose(deadClient)
            default:
                continue
        }
    }
}

func Super() {
    // make TCP listener
    listener, err := CreateListener()
    checkError(err)

    // initialize the channel that each client thread
    // will use to pull screenshots from the image server

    workerGroup := newWorkerGroup(1)

    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        imageChan := make(chan *FBUpdateWithImage)
        alertChan := make(chan chan *FBUpdateWithImage)
        // begin buffering images
        workerGroup.BroadcastImageChans(imageChan)
        workerGroup.BroadcastControlMsg(Start)

        // launch thread that will be responsible for cleaning up after
        // client disconnects
        go cleanUpCrew(workerGroup, alertChan)

        // launch client thread and return to listening
        go handleClient(imageChan, conn, alertChan)
    }
}

// when you create imageChans, workerGroup, etc, create deadClients channel
// workerGroup needs to listen on deadClients, and invoke wg.BroadcastChanToClose on each incoming message

func handleClient(c chan *FBUpdateWithImage, conn net.Conn, alertChan chan chan *FBUpdateWithImage) {
    //handshake (version exchange)
    defer func() {
        // Send the "i'm dying" message on the channel
        alertChan <- c
        //send message to thread that handles garbage collection
        conn.Close()
    }()
 
    _, err := exchangeVersions(conn)
    if err != nil {
        return
    }

    // send security level
    err = sendSecurity(conn)
    if err != nil {
        return
    }

    // init phase
    clientInitFlag, err := receiveClientInit(conn)
    // if client requests other connections can't join
    // inform client that this isn't possible (at the moment)
    if clientInitFlag != 1 || err != nil { 
        fmt.Printf("Cannot give exclusive access to server")
        return
    }

    // send serverinit msg
    pixelFormat := NewPixelFormat()
    serverInitMsg := NewServerInitMsg(pixelFormat)
    err = SendServerInit(serverInitMsg, conn)
    if err != nil {
        return
    }

    errChan := make(chan error, 10)

    // main loop starts
    for {
        //read from client
        msg, msgNum, err := GetMsg(conn)
        if err != nil {
            fmt.Println("err ", err, " getting message from conn")
            fmt.Println("closing connection 1")
            return
        }
        //do something with the message
        MsgDispatch(conn, msgNum, c, msg, errChan)
        select {
            case err := <-errChan:
                fmt.Println("Error Received from err chan", err)
                fmt.Println("closing connection 2")
                return
            default:
                continue

            }
    }

}

func newServerClientImageChans() *ServerClientImageChans {
    ic := &ServerClientImageChans{
        imageChannels : make([]chan *FBUpdateWithImage, 0),
    }
    return ic
}

func newWorkerGroup(numWorkers int) *WorkerGroup {
    wg := &WorkerGroup{
        workerControls : make([]chan ControlMessage, numWorkers),
        workerImageChanChans : make([]chan chan *FBUpdateWithImage, numWorkers),
        workerErrChans : make([]chan chan *FBUpdateWithImage, numWorkers),
    }
    var control chan ControlMessage
    
    for i := 0; i < numWorkers; i++ {
        control = make(chan ControlMessage)
        wg.workerControls[i] = control
        imageChanChan := make(chan chan *FBUpdateWithImage)
        wg.workerImageChanChans[i] = imageChanChan
        errChan := make(chan chan *FBUpdateWithImage)
        wg.workerErrChans[i] = errChan
        go imageServer(errChan, imageChanChan, control)
    }

    return wg
}


func (wg *WorkerGroup) BroadcastControlMsg(msg ControlMessage) {
    for _, ctrl := range wg.workerControls {
        go func(control chan ControlMessage) {
            control <- msg
        }(ctrl)
    }
}

func (wg *WorkerGroup) BroadcastImageChans(imageChan chan *FBUpdateWithImage) {
    for _, imgChanChan := range wg.workerImageChanChans {
        go func(imageChanChan chan chan *FBUpdateWithImage) {
            imageChanChan <- imageChan
        }(imgChanChan)
    }
}

func (wg *WorkerGroup) BroadcastChanToClose(imageChan chan *FBUpdateWithImage) {
    for _, workerErrChan := range wg.workerErrChans {
        go func(ErrChan chan chan *FBUpdateWithImage) {
            ErrChan <- imageChan
        }(workerErrChan)
    }
}

func (ic *ServerClientImageChans) BroadcastImage(image *FBUpdateWithImage) {
    if len(ic.imageChannels) > 0 {
        for _, imageChan := range ic.imageChannels {
            log.Printf("Sending image %p over image pipe %v\n", image, imageChan)
            imageChan <- image
        }
    }
}

// I should definitely stick the image channels in a map. Coming soon.
func removeChan(closedChan chan *FBUpdateWithImage, ic *ServerClientImageChans) {
    for i, imageChan := range ic.imageChannels {
        if imageChan == closedChan {
            ic.imageChannels = append(ic.imageChannels[:i], ic.imageChannels[i+1:]...)
        }
    }
}

func imageServer(errChan chan chan *FBUpdateWithImage, imageChanChan chan chan *FBUpdateWithImage, control chan ControlMessage) {
    ic := newServerClientImageChans()
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
            case imageChan := <- imageChanChan:
                ic.imageChannels = append(ic.imageChannels, imageChan)
            case closedChan := <- errChan:
                fmt.Print("close channel %v", closedChan)
                removeChan(closedChan, ic)
            default:
                if running {
                    image := NewFBUpdateWithImage()
                    ic.BroadcastImage(image)
                }
        }
    }
}
