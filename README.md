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

####Super:
A supervisor thread that oversees the client server/connection and initializes and launches the client and image server threads. It can access the image server threads by way of several channels stored in a struct called WorkerGroup. Having separate super and client threads allows for multiple clients to connect to the server at once.

####Image Server:
The image server is responsible for taking screenshots and putting them on channels that can be read by client threads and sent on the network. The image server can be a single thread taking, processing, and writing screenshots, or it can be initialized as multiple threads doing these things concurrently. 

Currently, writing the very large image data to the network is a speed bottleneck, so splitting the server into threads that can take and process more than 2 or 3 images per second is a waste of overhead, because the network can't handle more than 2 or 3 images per second. Once I implement better image encoding, I'll experiment with splitting the server into more threads.

####Client:
The client thread handles the handshake, initialization, and main communication loop with the VNC client. In the main loop, the client thread reads incoming client messages and determines the proper response message. Then, each response is constructed and sent in its own thread. Because of this, the write functionality doesn’t run the risk of blocking the read functionality until everything is written to the network. When the client returns, it sends a “dying” message to a clean up thread which removes the client’s unique channel from the image server.

####Clean up Crew:
Is responsible for notifying image server to remove the dying client’s channel.

##Implementation Details
####Server and Control Flow:
#####Before The Client Connects
When the VNC server is launched, Super instantiates a TCP listener on port 5900. It then calls newWorkerGroup() to both launch the image server, and construct a workerGroup--a struct that has access to the image server's channels. 

Each image server thread is initialized with three channels:

1. The workerControls channel accepts control messages, Start and Stop, which start and stop the production and processing of images. 

2. workerImageChanChans is a channel that accepts image channels. Each client gets a unique image channel upon connection. The reason for this is that reading from any given channel is destructive. If there were more than one client reading from one image channel, each client would only see ~ (number of screenshots written to channel) / (number of clients connected). Instead, when a client connects, Super creates a unique image channel. The image channel is given to the client thread as a parameter, and is broadcast to all image server threads through their image channel channels. The image server threads then add the image channel to its list of active client image channels. 

3. workerErrChans workerErrChan is a channel that accepts imageChannels. It allows for the destruction of active client image channels when clients disconnect. Channels block when a sender writes to it and a receiver isn't reading from it. Conversely, channels block when a receiver tries to read from a channel and a sender isn't writing to it. Not deleting image channels would mean the image server would continue to try to write images to the channel, but there would be no client reading from the channel, creating a block. 

Once the image channel is initialized, Super waits in a loop for a connection. Once it has one, several events occur.

1. A unique image channel is made for the client and broadcasted to the image server threads. 
2. A "Start" message is sent to the image server to start taking and processing screenshots.

3. An alert channel is created so the client can alert the cleanup thread that it is closing.

4. cleanUpCrew--the thread in charge of cleanup--is launched as a go routine with access to workerGroup and the alert channel

5. handleClient--the thread in charge of client processes--is launched with access to its image channel, and its alert channel. 

In progress

## Final Thoughts

In progress