package msg

import (
    "io"
    "log"
    "fmt"
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
    SetEncodingsLen int = 3
    FBUpdateRequestLen int = 9
    KeyEventLen int = 7 
    PointerEventLen int = 5
)

func GetMsg(reader io.Reader) (msg []byte, msgKind MsgKind) {
    msgKind = readMsgKind(reader)
    msgLength := getMsgLength(msgKind)
    msg = readMsg(msgLength, reader)
    return
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

func getMsgLength(msgKind MsgKind) int {
    switch msgKind {
    case SetPixelFormat:
        return SetPixelFormatLen
    case SetEncodings:
        return SetEncodingsLen
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
    return buf
}

