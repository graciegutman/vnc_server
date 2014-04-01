package main

import (
    "vnc/vnc"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	//"bytes"
	// "io/ioutil"
)

const SERVER_NAME string = ("Gracie")

type PixelFormat struct {
	bits_per_pixel   uint8
	depth            uint8
	big_endian_flag  uint8
	true_colour_flag uint8
	red_max          uint16
	green_max        uint16
	blue_max         uint16
	red_shift        uint8
	green_shift      uint8
	blue_shift       uint8
	padding          [3]byte
}

type ServerInit struct {
	fb_width            uint16
	fb_height           uint16
	server_pixel_format PixelFormat
	name_length         uint32
	name_string         [3]byte
}

type FrameBufferUpdate struct {
	message_type         uint8
	padding              [1]byte
	number_of_rectangles uint16
	x                    uint16
	y                    uint16
	width                uint16
	height               uint16
	encoding_type        int32
}

func main() {

	//the port we'll be listening on
	service := ":5900"
	//we get a pointer to a TCPAddr struct
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	//gives us a *TCPListener interface from a TCPAddress and a net
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		//If there is an error, return to the beginning of the loop
		if err != nil {
			continue
		}
		fmt.Println("connection went through")

		//Send version number
		version := "RFB 003.003\n"
		conn.Write([]byte(version))

		//Reads Version response
		buf := make([]byte, 12)
		_, err = conn.Read(buf[0:])
		version_resp := buf
		fmt.Println(version_resp)

		//Sends Security version
		var security uint32 = 1
		binary.Write(conn, binary.BigEndian, security)

		//Reads Security response
		//Beginning of initialization phase (client flag)
		buf3 := make([]byte, 1)
		resp, err := conn.Read(buf3)
		fmt.Println(resp)

		//Server sends a bunch of stuff about formats it will use
		pf := PixelFormat{
			bits_per_pixel:   32,
			depth:            24,
			big_endian_flag:  0,
			true_colour_flag: 1,
			red_max:          255,
			green_max:        255,
			blue_max:         255,
			red_shift:        16,
			green_shift:      8,
			blue_shift:       0}

		pixelData := &ServerInit{
			fb_width:            1280,
			fb_height:           1024,
			server_pixel_format: pf,
			name_length:         3,
			name_string:         [3]byte{1, 2, 3},
		}
		//Write ServerInit message to conn
		binary.Write(conn, binary.BigEndian, pixelData)
		/*Okay: not that it hadn't been hairy and gross before, but here is
		  my abomination. This is the beginning of the phase in which I couldn't
		  just hardcode a "read this many bytes" function, because this is the
		  part where I'd have to take multiple different types of requests.
		  My first priority was "get green rectangle on the screen," so I basically
		  called read for some arbitrary amount of bytes, just to block long enough
		  that I knew I had a request for a framebuffer. Then I just wrote a
		  framebuffer in a loop. Voila, framebuffer. And then tears. Well,
		  this was a nice experiment. Time to start over with what I know now.*/

		for i := 0; i < 20; i++ {
            message, msgnum := vnc.GetMsg(conn)
            fmt.Println(message, msgnum)
			//Create pix array of all green values
			pix_array := createPixArray(int(pixelData.fb_height), int(pixelData.fb_width))

			fb_update := &FrameBufferUpdate{
				number_of_rectangles: 1,
				x:                    0,
				y:                    0,
				width:                pixelData.fb_width,
				height:               pixelData.fb_height,
				encoding_type:        0}

			err := binary.Write(conn, binary.BigEndian, fb_update)
			fmt.Printf("Fb update error: %v\n", err)
			err = binary.Write(conn, binary.LittleEndian, pix_array)
			fmt.Printf("pix array error: %v\n", err)

		}

		conn.Close()
	}
}

func createPixArray(width, height int) []uint32 {
	size := width * height
	pix_slice := make([]uint32, size)
	for i := 0; i < (size); i++ {
		pix_slice[i] = 65280
	}
	return pix_slice
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s \n", err.Error())
		os.Exit(1)
	}
}
