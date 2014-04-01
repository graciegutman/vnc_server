/* IP
 */

package main

// takes an argument of a valid IP address as a parameter.
import (
	"fmt"
	"net"
	"os"
)

func main() {
	//os.Args works the same way python's os args works. First arg
	//is the file path itself
	if len(os.Args) != 2 {
		//print to std error
		fmt.Fprintf(os.Stderr, "Usage: %s ip-addr\n", os.Args[0])
		//exit process with an error
		os.Exit(1)
	}
	name := os.Args[1]
	//net.ParseIP apparently has some wonky bugs. Parses name
	//to valid IP address if possible
	addr := net.ParseIP(name)
	if addr == nil {
		fmt.Println("Invalid address")
	} else {
		fmt.Println("The address is ", addr.String())
	}
	//Exit without an error
	os.Exit(0)
}
