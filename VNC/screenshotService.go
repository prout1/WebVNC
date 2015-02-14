package VNC

import (
	"image"
	//"log"
	"reflect"
	//"sync"
	"time"
	"unsafe"
)

/*
#cgo LDFLAGS:-lgdi32
#include <stdio.h>
#include <stdlib.h>

#include <windows.h>
#include <wingdi.h>

void ScreenCap(void * ScreenData)
{
	HDC hScreen = GetDC(GetDesktopWindow());
	int ScreenX = GetDeviceCaps(hScreen, HORZRES);
	int ScreenY = GetDeviceCaps(hScreen, VERTRES);
	int BitsPerPixel = GetDeviceCaps(hScreen, BITSPIXEL);
	BITMAPINFOHEADER bmi = { 0 };
	//BYTE* ScreenData = (BYTE*)malloc((BitsPerPixel/8) * ScreenX * ScreenY);
	HDC hdcMem = CreateCompatibleDC(hScreen);
	HBITMAP hBitmap = CreateCompatibleBitmap(hScreen, ScreenX, ScreenY);
	HGDIOBJ hOld = SelectObject(hdcMem, hBitmap);

	bmi.biSize = sizeof(BITMAPINFOHEADER);
	bmi.biPlanes = 1;
	bmi.biBitCount = BitsPerPixel;
	bmi.biWidth = ScreenX;
	bmi.biHeight = -ScreenY;
	bmi.biCompression = BI_RGB;
	bmi.biSizeImage = (BitsPerPixel * ScreenX * ScreenY) / 8;

	BitBlt(hdcMem, 0, 0, ScreenX, ScreenY, hScreen, 0, 0, SRCCOPY);
	SelectObject(hdcMem, hOld);

	GetDIBits(hdcMem, hBitmap, 0, ScreenY, ScreenData, (BITMAPINFO*)&bmi, DIB_RGB_COLORS);

	{
		// flipping red and blue byte values, because Go needs them in different order than this
		int i, lim = ScreenX*ScreenY;
		int * buf = (int*)ScreenData;
		int y;

		for(i=0; i<lim; ++i) {
			y = buf[i];
			buf[i] = ((y & 0xff00ff00) | ((y & 0x000000ff) << 16) | ((y & 0x00ff0000) >> 16));
		}
	}
	ReleaseDC(GetDesktopWindow(), hScreen);
	DeleteDC(hdcMem);
	DeleteObject(hBitmap);
	//return ScreenData;
}
*/
import "C"

type scrShotService struct {
	currentShot *image.RGBA
	server      *VNCServer
	timeout     time.Duration
	requests    chan *FrameBufferRequest
}

func (s *scrShotService) Init(serv *VNCServer) {
	s.server = serv
	s.currentShot = nil
	s.requests = make(chan *FrameBufferRequest, ChanBufferSize)
	s.timeout = 2000 * time.Millisecond
}

func (s *scrShotService) Run() {
	working := false
	for {
		select {
		case <-s.requests:
			if !working {
				go func() {
					// for now assuming that update requests will ask for full screens
					// need to think what to do with requests of smaller rectangles and how to optimize that
					working = true
					//log.Println("screen update")
					s.server.client.frameUpdatesChan <- s.captureScreen().Pix
					working = false
				}()
			}
			//TODO : add disconnect and error chan here as well
		case <-s.server.disconnectChan:
		case <-s.server.errorChan:
			return
		}

	}
}

func (s *scrShotService) Stop(serv *VNCServer) {

}

func (s *scrShotService) Request(req *FrameBufferRequest) {
	s.requests <- req
}

func getScreenDimensions() (width, height, bytesPerPixel int32) {
	width = int32(C.GetSystemMetrics(C.SM_CXSCREEN))
	height = int32(C.GetSystemMetrics(C.SM_CYSCREEN))
	bytesPerPixel = 4
	return
}

func CaptureScreenInternal() *image.RGBA {
	width, height, bytesPerPixel := getScreenDimensions()
	pix := make([]byte, width*height*bytesPerPixel, width*height*bytesPerPixel)
	h := (*reflect.SliceHeader)(unsafe.Pointer(&pix))
	pixelData := unsafe.Pointer(h.Data)

	//TODO error handling
	C.ScreenCap(pixelData)

	// construct the image
	res := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	res.Pix = pix
	return res
}

func (s *scrShotService) captureScreen() *image.RGBA {
	// doing screen capturing at minimum s.timeout
	//var wg sync.WaitGroup
	var res *image.RGBA = nil
	//log.Println(s.timeout)
	//go func() {
	//	wg.Add(1)
	res = CaptureScreenInternal()
	//	wg.Done()
	//}()
	//go func() {
	//	wg.Add(1)
	//	time.Sleep(s.timeout)
	//	wg.Done()
	//}()

	//wg.Wait()
	return res
}
