
package main

import (
	"fmt"
	"net"
    "encoding/binary"
	"os"
    "bytes"
   // "io/ioutil"
)
const SERVER_NAME string = ("Gracie")

type PixelFormat struct {
    bits_per_pixel uint8
    depth uint8
    big_endian_flag uint8
    true_colour_flag uint8
    red_max uint16
    green_max uint16
    blue_max uint16
    red_shift uint8
    green_shift uint8
    blue_shift uint8
    padding [3]byte 
}

type ServerInit struct {
    fb_width uint16
    fb_height uint16
    server_pixel_format PixelFormat
    name_length uint32
    name_string [len(SERVER_NAME)]byte
}

func main() {

	//the port we'll be listening on
	service := ":5900"
	//we get a pointer to a TCPAddr struct
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	//gives us a Conn interface from a TCPAddress and a net
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		//Waits for and returns next connection to listener
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
        version_resp := string(buf[0:])
        fmt.Println("response was ", version_resp)

        //Sends Security version 
        var security uint32 = 1
        binary.Write(conn, binary.BigEndian, security)

        //Reads Security response
        //Beginning of initialization phase (client flag)
        buf3 := make([]byte, 1)
        resp, err := conn.Read(buf3[0:])
        fmt.Println("response was ", resp)
		conn.Close()

        //Server sends a bunch of stuff about formats it will use
        pixelData := new(ServerInit)
        test_buff := new(bytes.Buffer)
        binary.Write(test_buff, binary.BigEndian, pixelData)
        fmt.Printf("%b", test_buff.Bytes())
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s \n", err.Error())
		os.Exit(1)
	}
}

