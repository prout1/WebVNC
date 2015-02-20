package main

import (
	"WebVNC/VNC"
	"fmt"
	"net"
)

func main() {
	vncServer := VNC.CreateServer("5555")
	go vncServer.Run()

	// create client connection
	conn, err := net.Dial("tcp", "127.0.0.1:5555")
	if err != nil {
		fmt.Println(err)
		return
	}

	proxy := VNC.CreateProxy(conn)
	proxy.Run()
}
