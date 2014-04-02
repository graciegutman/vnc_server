package vnc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	//"log"
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
		fmt.Println("can't determine length")
	}
	return 0, nil
}

func readMsg(msgLength int, reader io.Reader) ([]byte, error) {
	buf := make([]byte, msgLength)
	_, err := reader.Read(buf)
	return buf, err
}

func ParseClickEvent(clickMsg []byte) MouseData {
	var resRatio float64 = 1.25
	var clickMask uint8
	var x uint16
	var y uint16

	b1 := bytes.NewReader(clickMsg[:1])
	b2 := bytes.NewReader(clickMsg[1:3])
	b3 := bytes.NewReader(clickMsg[3:])

	binary.Read(b1, binary.BigEndian, &clickMask)
	binary.Read(b2, binary.BigEndian, &x)
	binary.Read(b3, binary.BigEndian, &y)

	// accounts for stupid retina display scaling shenaniganry
	var fx, fy float64 = (float64(x) * resRatio), (float64(y) * resRatio)

	return MouseData{
		clickMask: clickMask,
		x:         fx,
		y:         fy,
	}
}

func getSetEncodingsLen(reader io.Reader) (msgLength int, err error) {
	//fmt.Println("in Parse Set Encodings")
	// read 3 bytes from the message
	var num32BitInts uint16
	buf := make([]byte, 3)
	_, err = reader.Read(buf)
	if err != nil {
		return 0, err
	}
	//fmt.Println(buf)
	// bytes 2 and 3 make up a uint16 value that tells
	// us how many 32-bit integers follow
	b := bytes.NewReader(buf[1:])
	binary.Read(b, binary.BigEndian, &num32BitInts)
	// 4 bytes per 32 bit integer.

	msgLength = int(num32BitInts) * 4
	return msgLength, err
}
