/* DaytimeServer
 */

/* SERVER BASICS
A server registers itself on a port and listens. It blocks on an "accept"
operation (program stops and waits on the accept operation). When a client
connects, the accept call returns with a connection object.

A daytime server waits for a connection, writs the current time to the client
and then closes the connection and resumes waiting.

relevant calls are:

func ListenTCP(net string, laddr *TCPAddr) (1, *TCPListener, err os.Error)
func (l *TCPListener) Accept() (c Conn, err os.Error)

The first takes a network net, and a local TCPAddress to listen on
The latter is a method called on a TCPListener that "accepts" a request
and returns a Conn interface.*/

package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	//the port we'll be listening on
	service := ":1200"
	//service is a string literal of an address with a specified port
	//why is there nothing before the port specification?
	//because we want to listen on all network interfaces.
	//we get a pointer to a TCPAddr struct
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	//gives us a Conn interface from a TCPAddress and a net
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	//infinite loop
	for {
		//Waits for and returns next connection to listener
		conn, err := listener.Accept()
		//If there is an error, return to the beginning of the loop
		if err != nil {
			continue
		}

		daytime := time.Now().String()
		//explicitly convert string to a slice of bytes
		//we don't give a shit about the return value
		conn.Write([]byte(daytime))
		//done with this client. Return to loop beginning
		conn.Close()
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
