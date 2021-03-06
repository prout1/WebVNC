package VNC

import (
	//"image"
	"log"
	"net"
	//"time"
)

var (
	ChanBufferSize = 1024
	FPS            = 50
)

type VNCServer struct {
	isConnected    bool // only one client allowed to connect
	port           string
	client         *Client
	errorChan      chan error
	disconnectChan chan struct{}
	// message handlers

	// encoding handlers

	// frame buffer
	screenShotService *scrShotService
	keyboardService   *keybdService
	pointerService    *mouseService
}

func CreateServer(port string) *VNCServer {
	var server VNCServer
	server.isConnected = false
	server.port = port
	server.errorChan = make(chan error, ChanBufferSize)
	server.disconnectChan = make(chan struct{}, ChanBufferSize)

	server.pointerService = &mouseService{}
	server.pointerService.Init()

	server.keyboardService = &keybdService{}
	server.keyboardService.Init()

	server.screenShotService = &scrShotService{}
	server.screenShotService.Init(&server)
	return &server
}

func (s *VNCServer) startServices() {
	go s.pointerService.Run()
	go s.keyboardService.Run()
	go s.screenShotService.Run()
}

func (s *VNCServer) Run() {
	s.startServices()

	ln, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		log.Println(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("VNCServer.Run: ", err)
			// panic
			return
		}

		if s.isConnected {
			log.Println("Server already connected, sorry")
		} else {
			s.isConnected = true
			s.client = newClient(s, conn)
			go s.client.handle()
		}
	}
}

func (s *VNCServer) Stop() {
	s.keyboardService.Stop()
	s.pointerService.Stop()
	s.screenShotService.Stop()
}

func (s *VNCServer) sendError(err error) {
	s.errorChan <- err
}

func (s *VNCServer) isSupportedEncoding(encoding int32) bool {
	return true
}
