package VNC

import (
	//"image"
	"log"
	"net"
	//"time"
)

var (
	Port           string = "5555"
	ChanBufferSize        = 16
	FPS                   = 50
)

type VNCServer struct {
	isConnected    bool // only one client allowed to connect
	client         *Client
	errorChan      chan struct{}
	disconnectChan chan struct{}
	// frame buffer
	screenShotService *scrShotService
	keyboardService   *keybdService
	pointerService    *mouseService
}

func CreateServer() *VNCServer {
	var server VNCServer
	server.isConnected = false
	server.errorChan = make(chan struct{}, ChanBufferSize)
	server.disconnectChan = make(chan struct{}, ChanBufferSize)

	server.pointerService = &mouseService{}
	server.pointerService.Init()
	server.keyboardService = &keybdService{}
	server.keyboardService.Init()
	server.screenShotService = &scrShotService{}
	server.screenShotService.Init(&server)
	return &server
}

func (s *VNCServer) Run() {
	ln, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Println(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("VNCServer.Run: ", err)
			// panic
		}

		if s.isConnected {
			log.Println("Server already connected, sorry")
			// send something appropriate to the client maybe ?
		} else {
			log.Println("server connected")
			s.isConnected = true
			s.client = newClient(s, conn)
			go s.client.handle()
		}
	}
}

func (s *VNCServer) isSupportedEncoding(encoding int32) bool {
	// TODO this method should go to server class
	return true
}

func makeScreenshotRGBA() []byte {
	return []byte("sucker!")
}

/*func (s *VNCServer) HandleScreenshots() {
	//testChan := make(chan struct{}, 1)
	working := false

	for {
		select {
		case <-s.requestUpdate:
			if !working {
				go func() {
					// for now assuming that update requests will ask for full screens
					// need to think what to do with requests of smaller rectangles and how to optimize that
					working = true
					res := makeScreenshotRGBA()
					s.client.frameUpdatesChan <- res
					time.Sleep(time.Second / time.Duration(FPS))
					working = false
				}()
			}
		case <-s.disconnectChan:
		case <-s.errorChan:
			return
		}

	}
}*/
