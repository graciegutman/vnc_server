vnc_server
==========

A very slow, eminently multithreaded VNC server in Go

Files Contained in This Repository
__________________________________

main.go: build with CC=clang go build main.go

vnc/createmsg.go: message serialization

vnc/readMsg.go: message deserialization

vnc/screenDecode.go: screenshot grabbing and processing

vnc/mouse.go: cursor handling

vnc/server.go: server control flow and utilities

Overview
_________

When it came time to decide what to hack on for four weeks at Hackbright, I knew I was less interested in the product side of development, and more interested in dispelling as much magic and learning as many foreign programming concepts as possible. I ended up choosing to write a VNC server in Go, which would expose me to networks, lower level OS concepts, concurrency, and a new language. 
