vnc_server
==========

A very slow, eminently multithreaded VNC server in Go

#H2 Files Contained in This Repository

main.go: build with CC=clang go build main.go
vnc/createmsg.go: message serialization
vnc/readMsg.go: message deserialization
vnc/screenDecode.go: screenshot grabbing and processing
vnc/mouse.go: cursor handling
vnc/server.go: server control flow and utilities

#H4 A Disclaimer About Requirements and Dependencies

This VNC server is a student project; I didn't build it with portability in mind. If you really intend to run this server, you need to have Imagemagick, Cocoa, and Go installed. The server has not been tested on any platform other than a Mac OSX 10.9. 

#H2 Overview

When it came time to decide what to hack on for four weeks at Hackbright, I
knew I was less interested in the product side of development, and more
interested in dispelling as much magic and learning as many foreign
programming concepts as possible. I ended up choosing to write a VNC server in
Go, which would expose me to networks, lower level OS concepts, concurrency,
and a new language.

At the highest level, a VNC server works like this: the VNC server and VNC client connect and shake hands. The client asks the server for its screen data, and the server responds by sending the client screen data. The client sends the server info about its cursor and keyboard events, and the server reads the info and treats those click and keyboard events as its own. The combined effect of these interactions is remote control of a computer.

My VNC server makes heavy use of go routines--Go's lightweight threads--to support multiple clients, respond quickly to messages, and speed up screen grabbing. It currently supports some cursor events but not keyboard events (although implementation of those should be straightforward).

#H2 The Choice To Use Golang

In progress

#H2 High-Level Architecture 

In progress

#H2 Implementation Details

In progress

#H2 Final Thoughts

In progress