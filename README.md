# WebVNC

WebVNC is an application which allows users to reach a remote computer (acting as VNC server) using only a web browser.
This application has several components, which can be used combined, as well as seperately. 

- VNC server, implementing the VNC standart and designed to be easily extensible.
- VNC cilent - easily extensible and VNC compatible as well. This layer is only responsible for assembling and sending messages to 
the VNC server, it's input can come from anywhere - service on the client machine, networks, etc.
- WebSocket proxy - optional layer of the application, made for convenience. It's basically a wrapper of the client, taking JSON 
messages from the browser via websocket and using the client api to send the data to the server in a VNC standart compatible way.
- Browser client - javascript client layer, used to capture user input and send according JSON messages to the websocket proxy.

##Install: 
```
go get "https://github.com/prout1/WebVNC.git"
```

##Usage:

```go
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

```
