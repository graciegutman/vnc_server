vnc_server
==========

A very slow, eminently multithreaded VNC server in Go

##Files Contained in This Repository

+ main.go: build with CC=clang go build main.go
+ vnc/createmsg.go: message serialization
+ vnc/readMsg.go: message deserialization
+ vnc/screenDecode.go: screenshot grabbing and processing
+ vnc/mouse.go: cursor handling
+ vnc/server.go: server control flow and utilities

#####A Disclaimer About Requirements and Dependencies

This VNC server is a student project; I didn't build it with portability in mind. If you really intend to run this server, you need to have Imagemagick, Cocoa, and Go installed. The server has not been tested on any platform other than a Mac OSX 10.9. 

##Overview

When it came time to decide what to hack on for four weeks at Hackbright, I
knew I was less interested in the product side of development, and more
interested in dispelling as much magic and learning as many foreign
programming concepts as possible. I ended up choosing to write a VNC server in
Go, which would expose me to networks, lower level OS concepts, concurrency,
and a new language (which turned out to be two new languages).

At the highest level, a VNC server works like this: the VNC server and VNC client talk to each other using the Remote Framebuffer (RFB) protocol. After a brief handshake and initialization, the client asks the server for its screen data, and the server responds by sending the client screen data. The client sends the server info about its cursor and keyboard events, and the server reads the info and treats those click and keyboard events as its own. The combined effect of these interactions is remote control of the server computer.

My VNC server makes heavy use of go routines--Go's lightweight threads--to support multiple clients, respond quickly to messages, and speed up screen grabbing. It currently supports some cursor events but not keyboard events (although implementation of those should be straightforward).

##The Choice To Use Go
I chose to use Golang for a number of reasons, not the least of which was my desire to learn a language that didn't abstract as much away from me as JavaScript and Python had. I believed Go was an appropriate choice for this project because of the following: 

+ The RFB protocol calls for very specific data types in its messages, and it made sense to use a language that was outright about what types of data it was using (for example, uint32 vs. int16)
+ Go has an interesting and accessible way of dealing with multi-threading. Go routines are lightweight threads, which can pass data to each other through structures called channels (as opposed to accessing common global data, which makes race conditions more likely).
+ The lower level-ness of the language allowed for finer-grained control of network connections and OS processes.

##High-Level Architecture 

My actual server design is composed of four main segments, each of which occupies its own thread or multiple threads. 

###Super:
A supervisor thread that oversees the client server/connection and initializes and launches the client and image server threads. It can access the image server threads by way of several channels stored in a struct called WorkerGroup. Having separate super and client threads allows for multiple clients to connect to the server at once.

###Image Server:
The image server is responsible for taking screenshots and putting them on channels that can be read by client threads. The image server can be initialized as an arbitrary number of threads to increase the frequency by which screenshots are taken.

###Client:
The client thread handles the handshake, initialization, and main communication loop with the VNC client. In the main loop, the client thread reads incoming client messages and determines the proper response message. Then, each response is constructed and sent in its own thread. Because of this, the write functionality doesn’t run the risk of blocking the read functionality until everything is written to the network. When the client returns, it sends a “dying” message to a clean up thread which removes the client’s unique channel from the image server.

###Clean up Crew:
Is responsible for notifying image server to remove the dying client’s channel.

## Implementation Details
### Server and Control Flow


In progress

## Final Thoughts

In progress