package vnc

import (
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

func Super() {
	// make TCP listener
	listener, err := CreateListener()
	checkError(err)

	// initialize the image server thread(s) (workers)
	// and create a record of all worker channels
	workerGroup := newWorkerGroup(2)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		//AFTER THIS POINT, CONNECTION IS ESTABLISHED WITH CLIENT

		imageChan := make(chan *FBUpdateWithImage)
		alertChan := make(chan chan *FBUpdateWithImage)

		// send a client-unique image channel to the workers
		workerGroup.BroadcastImageChans(imageChan)

		// signal the workers to start taking screenshots and putting
		// them on client image channel(s)
		workerGroup.BroadcastControlMsg(Start)

		// launch thread that will be responsible for cleaning up after
		// client disconnects
		go cleanUpCrew(workerGroup, alertChan)

		// launch client thread and return to listening
		go handleClient(imageChan, conn, alertChan)
	}
}

func handleClient(c chan *FBUpdateWithImage, conn net.Conn, alertChan chan chan *FBUpdateWithImage) {

	defer func() {
		// Alerts cleanUpCrew that client is closing
		alertChan <- c
		conn.Close()
	}()

	// HANDSHAKE WITH CLIENT BEGIN

	// send and receive version info
	_, err := exchangeVersions(conn)
	if err != nil {
		return
	}

	// send security level supported by server
	err = sendSecurity(conn)
	if err != nil {
		return
	}

	// HANDSHAKE WITH CLIENT END

	// INIT WITH CLIENT BEGIN

	clientInitFlag, err := receiveClientInit(conn)
	// If client requests other connections can't join,
	// disconnect; feature isn't currently supported
	if clientInitFlag != 1 || err != nil {
		log.Printf("Cannot give exclusive access to server \n")
		return
	}

	// Exchange info about what image formats are going to be used
	pixelFormat := NewPixelFormat()
	serverInitMsg := NewServerInitMsg(pixelFormat)
	err = SendServerInit(serverInitMsg, conn)
	if err != nil {
		return
	}

	// INIT WITH CLIENT END

	// Allows client to receive errors from go routines client launches
	errChan := make(chan error, 10)

	// MAIN LOOP

	for {
		// Get message from client
		msg, msgNum, err := GetMsg(conn)
		if err != nil {
			log.Printf("Error received reading message %v\n", err)
			return
		}
		// React to message
		MsgDispatch(conn, msgNum, c, msg, errChan)
		select {
		case err := <-errChan:
			log.Printf("Error received from err chan %v\n", err)
			return
		default:
			continue

		}
	}

}

func imageServer(errChan chan chan *FBUpdateWithImage, imageChanChan chan chan *FBUpdateWithImage, control chan ControlMessage) {
	// Initialize a record of worker to client image channels
	ic := newServerClientImageChans()
	var running bool = false
	for {
		select {
		case msg := <-control:
			switch msg {
			case Start:
				running = true
			case Stop:
				running = false
			}
		// If super sends an image channel, append it to list of active
		// client image channels
		case imageChan := <-imageChanChan:
			ic.imageChannels = append(ic.imageChannels, imageChan)

		// If cleanUpCrew sends an image channel, remove client channel
		// from list of active client image channels
		case closedChan := <-errChan:
			log.Printf("close channel %v\n", closedChan)
			removeChan(closedChan, ic)
		default:
			if running {
				image := NewFBUpdateWithImage()
				ic.BroadcastImage(image)
			}
		}
	}
}

func newServerClientImageChans() *ServerClientImageChans {
	ic := &ServerClientImageChans{
		imageChannels: make([]chan *FBUpdateWithImage, 0),
	}
	return ic
}

// initializes a worker group with numWorkers threads
func newWorkerGroup(numWorkers int) *WorkerGroup {
	wg := &WorkerGroup{
		workerControls:       make([]chan ControlMessage, numWorkers),
		workerImageChanChans: make([]chan chan *FBUpdateWithImage, numWorkers),
		workerErrChans:       make([]chan chan *FBUpdateWithImage, numWorkers),
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

// Sends Start and Stop message to workers
func (wg *WorkerGroup) BroadcastControlMsg(msg ControlMessage) {
	for _, ctrl := range wg.workerControls {
		go func(control chan ControlMessage) {
			control <- msg
		}(ctrl)
	}
}

// Sends new image channel to workers
func (wg *WorkerGroup) BroadcastImageChans(imageChan chan *FBUpdateWithImage) {
	for _, imgChanChan := range wg.workerImageChanChans {
		go func(imageChanChan chan chan *FBUpdateWithImage) {
			imageChanChan <- imageChan
		}(imgChanChan)
	}
}

// Sends a remove client channel message to workers
func (wg *WorkerGroup) BroadcastChanToClose(imageChan chan *FBUpdateWithImage) {
	for _, workerErrChan := range wg.workerErrChans {
		go func(ErrChan chan chan *FBUpdateWithImage) {
			ErrChan <- imageChan
		}(workerErrChan)
	}
}

// Sends image to all clients
func (ic *ServerClientImageChans) BroadcastImage(image *FBUpdateWithImage) {
	if len(ic.imageChannels) > 0 {
		for _, imageChan := range ic.imageChannels {
			log.Printf("Sending image %p over image pipe %v\n", image, imageChan)
			imageChan <- image
		}
	}
}

// Cleans up no-longer-used client channels on the image server when a client disconnects
func cleanUpCrew(wg *WorkerGroup, alertChan chan chan *FBUpdateWithImage) {
	for {
		select {
		case deadClient := <-alertChan:
			wg.BroadcastChanToClose(deadClient)
		default:
			continue
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
