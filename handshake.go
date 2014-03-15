
package main

import (
	"fmt"
	"net"
    "encoding/binary"
	"os"
)

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
        var buf [12]byte
        _, err = conn.Read(buf[0:])
        version_resp := string(buf[0:])
        fmt.Println("response was ", version_resp)

        //Sends Security version 
        var security uint32 = 1
        binary.Write(conn, binary.BigEndian, security)

        //Reads Security response
        var buf3 [1]byte
        resp, err := conn.Read(buf3[0:])
        fmt.Println("response was ", resp)
		conn.Close()
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s \n", err.Error())
		os.Exit(1)
	}
}

