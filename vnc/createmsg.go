package vnc

import(
    "net"
    "fmt"
    "os"
    "encoding/binary"
    "log"
    "io/ioutil"
)

const(
    FramebufferUpdate uint8 = 0
    SetColourMapEntries uint8 = 1
    Bell uint8 = 2
    ServerCutText uint8 = 3
)

const(
    VersionNumber string = "RFB 003.003\n"
    FBWidth uint16 = 1024
    FBHeight uint16 = 640
    ServerNameLen uint32 = 3
    SecurityType uint32 = 1
)

var ServerName [3]byte = [3]byte{1, 2, 3}

type ServerInit struct {
    fbWidth            uint16
    fbHeight           uint16
    serverPixelFormat PixelFormat
    nameLength         uint32
    nameString         [3]byte
}

type PixelFormat struct {
    bitsPerPixel   uint8
    depth            uint8
    bigEndianFlag  uint8
    trueColourFlag uint8
    redMax          uint16
    greenMax        uint16
    blueMax         uint16
    redShift        uint8
    greenShift      uint8
    blueShift       uint8
    padding          [3]byte
}

type FrameBufferUpdate struct {
    messageType         uint8
    padding             [1]byte
    numberOfRectangles   uint16
    x                    uint16
    y                    uint16
    width                uint16
    height               uint16
    encodingType         int32
}

type FrameBufferWithImage struct {
    frameBufferMsg FrameBufferUpdate
    pixelFormat []uint8
}

func CreateListener() (listener *net.TCPListener, err error) {
    //the port we'll be listening on
    service := ":5900"
    tcpAddr, err := net.ResolveTCPAddr("tcp", service)
    checkError(err)
    listener, err = net.ListenTCP("tcp", tcpAddr)
    checkError(err)
    return listener, err
}

func sendVersion(conn net.Conn) (err error) {
    _, err = conn.Write([]byte(VersionNumber))
    return err
}

func receiveVersion(conn net.Conn) (version string, err error) {
    buf := make([]byte, 12)
    _, err = conn.Read(buf)
    version = string(buf)
    return version, err
}

func ExchangeVersions(conn net.Conn) (versionFlag bool, err error) {
    err = sendVersion(conn)
    checkError(err)

    version, err := receiveVersion(conn)
    checkError(err)

    if version != VersionNumber {
        fmt.Fprintf(os.Stderr, "Fatal error: versions don't match")
        os.Exit(1)
    }
    versionFlag = true
    return versionFlag, err
}

func SendSecurity(conn net.Conn) (err error) {
    err = binary.Write(conn, binary.BigEndian, SecurityType)
    return err
}

func ReceiveClientInit(conn net.Conn) (clientInitFlag int, err error) {
    buf := make([]byte, 1)
    resp, err := conn.Read(buf)
    clientInitFlag = int(resp)
    return clientInitFlag, err
}

func NewServerInitMsg(pixelFormat PixelFormat) ServerInit {
    return ServerInit {
                    fbWidth: FBWidth, 
                    fbHeight: FBHeight, 
                    serverPixelFormat: pixelFormat,
                    nameLength: ServerNameLen,
                    nameString: ServerName,
    }
}

func NewPixelFormat() PixelFormat {
    return PixelFormat {
                    bitsPerPixel: 32,
                    depth: 24,
                    bigEndianFlag: 0,
                    trueColourFlag: 1,
                    redMax: 255,
                    greenMax: 255,
                    blueMax: 255,
                    redShift: 0,
                    greenShift: 8,
                    blueShift: 16,
    }
}

func SendServerInit(serverInitMsg ServerInit, conn net.Conn) (err error) {
    binary.Write(conn, binary.BigEndian, serverInitMsg)
    return err
}

func NewFrameBufferRaw(width, height uint16) FrameBufferUpdate {
    return FrameBufferUpdate {
                            messageType: 0,
                            numberOfRectangles: 1,
                            x: 0,
                            y: 0,
                            width: width,
                            height: height,
                            encodingType: 0,
    }
}

//func MakeAndSendFrameBufferRaw(conn net.Conn) (err error) {
    //fb, ps := NewFrameBufferWithImageRaw()
    //err = SendFrameBufferRaw(conn, fb, ps)
    //checkError(err)
    //return 
//}

func NewFrameBufferWithImageRaw(c chan *FrameBufferWithImage) {
    for {
        f, err := ioutil.TempFile("", "screenshot")
        if err != nil {
            log.Fatal("could not create temp file")
        }
        
        err = TakeScreenShot(f)
        if err != nil{
            log.Fatal("screenshot failed")
        }
        
        err = ResizeImage(f)
        checkError(err)

        image, _ := DecodeFileToPNG(f)
        width, height := GetImageWidthHeight(image)
        newFrameBuffer := NewFrameBufferRaw(width, height)
        pixSlice, err := ImgDecodeRaw(image)
        checkError(err)
        
        os.Remove(f.Name())

        fb := &FrameBufferWithImage{newFrameBuffer, pixSlice}
        c <- fb
    }
}

func SendFrameBufferRaw(conn net.Conn, c chan *FrameBufferWithImage) {
    fb := <-c
    binary.Write(conn, binary.BigEndian, fb.frameBufferMsg)
    binary.Write(conn, binary.LittleEndian, fb.pixelFormat)
    return
}

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s \n", err.Error())
        os.Exit(1)
    }
}
