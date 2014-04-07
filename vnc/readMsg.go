package vnc

/* This is a file representing functions necessary for
reading off the connection and serializing. Not all messages
are fully supported yet. */

import (
	"bytes"
	"encoding/binary"
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
	Unknown					 MsgKind = 404
)

const (
	SetPixelFormatLen  int = 19
	FBUpdateRequestLen int = 9
	KeyEventLen        int = 7
	PointerEventLen    int = 5
)

type MouseData struct {
	clickMask uint8
	x         float64
	y         float64
}

// Get msg reads the first byte of each message to determine the type, 
// looks up the length of that type of message, and returns the rest
// of the message with the type.
func GetMsg(reader io.Reader) (msg []byte, msgKind MsgKind, err error) {
	msgKind, err = readMsgKind(reader)
	if err != nil {
		return nil, Unknown, err
	}
	msgLength, err := getMsgLength(msgKind, reader)
	if err != nil {
		return nil, Unknown, err
	}
	msg, err = readMsg(msgLength, reader)
	return msg, msgKind, err
}

func readMsgKind(reader io.Reader) (MsgKind, error) {
	buf := make([]byte, 1)
	_, err := reader.Read(buf)
	if err != nil {
		return Unknown, err
	}
	//casting int into MsgKind
	msgKind := MsgKind(buf[0])
	return msgKind, err
}

func getMsgLength(msgKind MsgKind, reader io.Reader) (int, error){
	switch msgKind {
	case SetPixelFormat:
		return SetPixelFormatLen, nil
	case SetEncodings:
		return getSetEncodingsLen(reader)
	case FramebufferUpdateRequest:
		return FBUpdateRequestLen, nil
	case KeyEvent:
		return KeyEventLen, nil
	case PointerEvent:
		return PointerEventLen, nil
	default:
		log.Printf("can't determine length")
	}
	return 0, nil
}

func readMsg(msgLength int, reader io.Reader) ([]byte, error) {
	buf := make([]byte, msgLength)
	_, err := reader.Read(buf)
	return buf, err
}

func ParseClickEvent(clickMsg []byte) MouseData {
	// resRatio for stupid retina display scaling shenaniganry
	var resRatio float64 = 1.25
	// type of button/cursor event
	var clickMask uint8
	// coordinates of cursor on screen
	var x uint16
	var y uint16
	
	b1 := bytes.NewReader(clickMsg[:1])
	b2 := bytes.NewReader(clickMsg[1:3])
	b3 := bytes.NewReader(clickMsg[3:])

	binary.Read(b1, binary.BigEndian, &clickMask)
	binary.Read(b2, binary.BigEndian, &x)
	binary.Read(b3, binary.BigEndian, &y)

	var fx, fy float64 = (float64(x) * resRatio), (float64(y) * resRatio)

	return MouseData{
		clickMask: clickMask,
		x:         fx,
		y:         fy,
	}
}

func getSetEncodingsLen(reader io.Reader) (msgLength int, err error) {
	// read 3 bytes from the message
	var num32BitInts uint16
	buf := make([]byte, 3)
	_, err = reader.Read(buf)
	if err != nil {
		return 0, err
	}
	// bytes 2 and 3 make up a uint16 value that tells
	// us how many 32-bit integers will follow in the message
	b := bytes.NewReader(buf[1:])
	binary.Read(b, binary.BigEndian, &num32BitInts)
	// 4 bytes per 32 bit integer.
	msgLength = int(num32BitInts) * 4
	return msgLength, err
}
