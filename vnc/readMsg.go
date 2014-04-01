package vnc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

type MsgKind int

const (
	SetPixelFormat           MsgKind = 0
	SetEncodings             MsgKind = 2
	FramebufferUpdateRequest MsgKind = 3
	KeyEvent                 MsgKind = 4
	PointerEvent             MsgKind = 5
	ClientCutText            MsgKind = 6
)

const (
	SetPixelFormatLen  int = 19
	FBUpdateRequestLen int = 9
	KeyEventLen        int = 7
	PointerEventLen    int = 5
)

type MouseData struct {
	clickMask uint8
	x         uint16
	y         uint16
}

func GetMsg(reader io.Reader) (msg []byte, msgKind MsgKind) {
	msgKind = readMsgKind(reader)
	fmt.Println("Message type:", msgKind)
	msgLength := getMsgLength(msgKind, reader)
	msg = readMsg(msgLength, reader)
	return
}

func readMsgKind(reader io.Reader) MsgKind {
	buf := make([]byte, 1)
	_, err := reader.Read(buf)
	fmt.Println(buf)
	if err != nil {
		log.Fatal(err)
	}
	msgKind := MsgKind(buf[0])
	return msgKind
}

func getMsgLength(msgKind MsgKind, reader io.Reader) int {
	switch msgKind {
	case SetPixelFormat:
		return SetPixelFormatLen
	case SetEncodings:
		return getSetEncodingsLen(reader)
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

func ParseClickEvent(clickMsg []byte) MouseData {
	var clickMask uint8
	var x uint16
	var y uint16

	b1 := bytes.NewReader(clickMsg[:1])
	b2 := bytes.NewReader(clickMsg[1:3])
	b3 := bytes.NewReader(clickMsg[3:])

	binary.Read(b1, binary.BigEndian, &clickMask)
	binary.Read(b2, binary.BigEndian, &x)
	binary.Read(b3, binary.BigEndian, &y)

	return MouseData{
		clickMask: clickMask,
		x:         x,
		y:         y,
	}
}

func getSetEncodingsLen(reader io.Reader) (msgLength int) {
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
