package astisrt

// #cgo LDFLAGS: -lsrt
// #include <srt/srt.h>
import "C"
import (
	"net"
	"sync"
	"syscall"
	"unsafe"
)

type ConnectCallback func(s *Socket, addr *net.UDPAddr, token int, err error)
type ListenCallback func(s *Socket, version int, addr *net.UDPAddr, streamID string) bool

var (
	acceptedSockets       = make(map[C.SRTSOCKET]*Socket)
	acceptedSocketsMutex  = &sync.Mutex{}
	connectCallbacks      = make(map[C.SRTSOCKET]ConnectCallback)
	connectCallbacksMutex = &sync.Mutex{}
	listenCallbacks       = make(map[C.SRTSOCKET]ListenCallback)
	listenCallbacksMutex  = &sync.Mutex{}
)

// For groups only
// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_connect_callback
func (s *Socket) SetConnectCallback(c ConnectCallback) (err error) {
	if err = cConnectCallback(s.u); err != nil {
		return
	}
	connectCallbacksMutex.Lock()
	connectCallbacks[s.u] = c
	connectCallbacksMutex.Unlock()
	return
}

//export go2cConnectCallback
func go2cConnectCallback(opaque unsafe.Pointer, u C.SRTSOCKET, errorcode C.int, peeraddr *C.struct_sockaddr, token C.int) {
	// Get callback
	connectCallbacksMutex.Lock()
	cb, ok := connectCallbacks[*(*C.SRTSOCKET)(opaque)]
	connectCallbacksMutex.Unlock()

	// No callback
	if !ok || cb == nil {
		return
	}

	// Create socket
	s := newSocketFromC(u)

	// Create addr
	addr, _ := udpAddrFromSockaddr((*syscall.RawSockaddrAny)(unsafe.Pointer(peeraddr)))

	// Callback
	cb(s, addr, int(token), newError(errorcode, 0))
}

//export go2cListenCallback
func go2cListenCallback(opaque unsafe.Pointer, u C.SRTSOCKET, version C.int, peeraddr *C.struct_sockaddr, streamid *C.char) int {
	// Get callback
	listenCallbacksMutex.Lock()
	cb, ok := listenCallbacks[*(*C.SRTSOCKET)(opaque)]
	listenCallbacksMutex.Unlock()

	// No callback
	if !ok || cb == nil {
		return int(C.SRT_SUCCESS)
	}

	// Create socket
	s := newSocketFromC(u)

	// Create addr
	addr, _ := udpAddrFromSockaddr((*syscall.RawSockaddrAny)(unsafe.Pointer(peeraddr)))

	// Callback
	if ok = cb(s, int(version), addr, C.GoString(streamid)); !ok {
		return int(C.SRT_ERROR)
	}

	// Store socket so that Accept() gets the proper Go object
	acceptedSocketsMutex.Lock()
	acceptedSockets[u] = s
	acceptedSocketsMutex.Unlock()
	return int(C.SRT_SUCCESS)
}

// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_listen_callback
func (s *Socket) SetListenCallback(c ListenCallback) (err error) {
	if err = cListenCallback(s.u); err != nil {
		return
	}
	listenCallbacksMutex.Lock()
	listenCallbacks[s.u] = c
	listenCallbacksMutex.Unlock()
	return
}
