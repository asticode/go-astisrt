package astisrt

// #cgo LDFLAGS: -lsrt
// #include <srt/srt.h>
import "C"
import (
	"errors"
	"fmt"
	"net"
	"syscall"
	"unsafe"
)

type Socket struct {
	rr RejectReason
	u  C.SRTSOCKET
}

// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_create_socket
func NewSocket() (s *Socket, err error) {
	// Create c socket
	var u C.SRTSOCKET
	if u, err = cCreateSocket(); err != nil {
		return
	}

	// Create socket
	s = newSocketFromC(u)
	return
}

func newSocketFromC(u C.SRTSOCKET) *Socket {
	return &Socket{u: u}
}

// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_close
func (s *Socket) Close() error {
	return cClose(s.u)
}

// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_bind
func (s *Socket) Bind(host string, port uint16) (err error) {
	// Create addr
	var addr *C.struct_sockaddr
	var size C.int
	if addr, size, err = sockAddrFromHostAndPort(host, port); err != nil {
		err = fmt.Errorf("astisrt: creating addr failed: %w", err)
		return
	}

	// Bind
	if err = cBind(s.u, addr, size); err != nil {
		err = fmt.Errorf("astisrt: binding failed: %w", err)
		return
	}
	return
}

// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_listen
func (s *Socket) Listen(backlog int) error {
	return cListen(s.u, C.int(backlog))
}

// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_connect
func (s *Socket) Connect(host string, port uint16) (addr *net.UDPAddr, err error) {
	// Create addr
	var caddr *C.struct_sockaddr
	var size C.int
	if caddr, size, err = sockAddrFromHostAndPort(host, port); err != nil {
		err = fmt.Errorf("astisrt: creating addr failed: %w", err)
		return
	}

	// Connect
	var rr C.int
	if rr, err = cConnect(s.u, caddr, size); err != nil {
		if errors.Is(err, ErrEconnrej) {
			s.rr = RejectReason(rr)
		}
		err = fmt.Errorf("astisrt: connecting failed: %w", err)
		return
	}

	// Create udp addr
	if addr, err = udpAddrFromSockaddr((*syscall.RawSockaddrAny)(unsafe.Pointer(caddr))); err != nil {
		err = fmt.Errorf("astisrt: creating udp addr failed: %w", err)
		return
	}
	return
}

// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_accept
func (s *Socket) Accept() (a *Socket, addr *net.UDPAddr, err error) {
	// Accept
	var u C.SRTSOCKET
	var raw syscall.RawSockaddrAny
	size := C.int(syscall.SizeofSockaddrAny)
	if u, err = cAccept(s.u, (*C.struct_sockaddr)(unsafe.Pointer(&raw)), &size); err != nil {
		err = fmt.Errorf("astisrt: accepting failed: %w", err)
		return
	}

	// Get socket
	acceptedSocketsMutex.Lock()
	var ok bool
	if a, ok = acceptedSockets[u]; ok {
		delete(acceptedSockets, u)
	} else {
		a = newSocketFromC(u)
	}
	acceptedSocketsMutex.Unlock()

	// Create udp addr
	if addr, err = udpAddrFromSockaddr(&raw); err != nil {
		err = fmt.Errorf("astisrt: creating udp addr failed: %w", err)
		return
	}
	return
}

// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_sendmsg2
func (s *Socket) SendMessage(b []byte) (int, error) {
	n, err := cSendMsg2(s.u, (*C.char)(unsafe.Pointer(&b[0])), C.int(len(b)), nil)
	return int(n), err
}

// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_recvmsg2
func (s *Socket) ReceiveMessage(b []byte) (int, error) {
	n, err := cRecMsg2(s.u, (*C.char)(unsafe.Pointer(&b[0])), C.int(len(b)), nil)
	return int(n), err
}
