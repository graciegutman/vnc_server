package vnc 

import (
    "bytes"
    "io"
    "log"
    "fmt"
    "encoding/binary"
)

type MsgKind int

const (
    SetPixelFormat MsgKind = 0
    SetEncodings MsgKind = 2
    FramebufferUpdateRequest MsgKind = 3
    KeyEvent MsgKind = 4
    PointerEvent MsgKind = 5
    ClientCutText MsgKind = 6
)

const (
    SetPixelFormatLen int = 19
    FBUpdateRequestLen int = 9
    KeyEventLen int = 7 
    PointerEventLen int = 5
)

func GetMsg(reader io.Reader) (msg []byte, msgKind MsgKind) {
    msgKind = readMsgKind(reader)
    fmt.Println("Message type:", msgKind)
    msgLength := getMsgLength(msgKind, reader)
    msg = readMsg(msgLength, reader)
    return
}

func parseSetEncodings(reader io.Reader) (msgLength int){
    fmt.Println("in Parse Set Encodings")
    // read 3 bytes from the message
    var num32BitInts uint16
    buf := make([]byte, 3)
    _, err := reader.Read(buf)
    if err != nil {
        log.Fatal("couldn't read Set Encodings Msg")
    }
    fmt.Println(buf)
    // bytes 2 and 3 make up a uint16 value that tells
    // us how many 32-bit integers follow
    b := bytes.NewReader(buf[1:])
    err = binary.Read(b, binary.BigEndian, &num32BitInts)
    // 4 bytes per 32 bit integer.

    msgLength = int(num32BitInts) * 4
    return msgLength
}

func readMsgKind(reader io.Reader) MsgKind {
    buf := make([]byte, 1)
    _, err := reader.Read(buf)
    if err != nil {
        log.Fatal("could not read first byte in buffer")
    }
    msgKind := MsgKind(buf[0])
    return msgKind
}

func getMsgLength(msgKind MsgKind, reader io.Reader) int {
    switch msgKind {
    case SetPixelFormat:
        return SetPixelFormatLen
    case SetEncodings:
        return parseSetEncodings(reader)
    case FramebufferUpdateRequest:
        return FBUpdateRequestLen
    case KeyEvent:
        return KeyEventLen
    case PointerEvent:
        return PointerEventLen
    default:
        fmt.Println("can't determine length")
   }
   return 0
}

func readMsg(msgLength int, reader io.Reader) []byte {
    buf := make([]byte, msgLength)
    _, err := reader.Read(buf)
    if err != nil {
        log.Fatal("could not read rest of message in buffer")
    }
    fmt.Println(buf)
    return buf
}

