package astisrt

import "C"
import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

func orderPortBytes(port uint16, o binary.ByteOrder) uint16 {
	b := (*[2]byte)(unsafe.Pointer(&port))
	return o.Uint16((*b)[:])
}

func sockAddrPortFromRegularPort(port uint16) uint16 {
	return orderPortBytes(port, binary.BigEndian)
}

func regularPortFromSockAddrPort(port uint16) uint16 {
	return orderPortBytes(port, binary.BigEndian)
}

type sockaddr struct {
	inet4 *syscall.RawSockaddrInet4
	inet6 *syscall.RawSockaddrInet6
}

func (sa *sockaddr) addr() *C.struct_sockaddr {
	if sa.inet4 != nil {
		return (*C.struct_sockaddr)(unsafe.Pointer(sa.inet4))
	} else if sa.inet6 != nil {
		return (*C.struct_sockaddr)(unsafe.Pointer(sa.inet6))
	}
	return nil
}

func (sa *sockaddr) size() C.int {
	if sa.inet4 != nil {
		return syscall.SizeofSockaddrInet4
	} else if sa.inet6 != nil {
		return syscall.SizeofSockaddrInet6
	}
	return 0
}

func (sa *sockaddr) toUDP() *net.UDPAddr {
	if sa.inet4 != nil {
		return &net.UDPAddr{
			IP:   sa.inet4.Addr[:],
			Port: int(regularPortFromSockAddrPort(sa.inet4.Port)),
		}
	} else if sa.inet6 != nil {
		return &net.UDPAddr{
			IP:   sa.inet6.Addr[:],
			Port: int(regularPortFromSockAddrPort(sa.inet6.Port)),
		}
	}
	return nil
}

func newSockaddrFromHostAndPort(host string, port uint16) (addr *sockaddr, err error) {
	// Parse ip
	ip := net.ParseIP(host)
	if ip == nil {
		// Lookup ip
		var ips []net.IP
		if ips, err = net.LookupIP(host); err != nil {
			err = fmt.Errorf("astisrt: looking up ip failed: %w", err)
			return
		}

		// No ip
		if len(ips) < 1 {
			err = errors.New("astisrt: ip lookup didn't return enough ips")
			return
		}

		// Set ip
		ip = ips[0]
	}

	// Check ip
	if v := ip.To4(); v != nil {
		// Create sockaddr
		addr = &sockaddr{inet4: &syscall.RawSockaddrInet4{
			Family: unix.AF_INET,
			Port:   sockAddrPortFromRegularPort(port),
		}}

		// Copy addr
		copy(addr.inet4.Addr[:], v)
	} else if v := ip.To16(); v != nil {
		// Create sockaddr
		addr = &sockaddr{inet6: &syscall.RawSockaddrInet6{
			Family: unix.AF_INET6,
			Port:   sockAddrPortFromRegularPort(port),
		}}

		// Copy addr
		copy(addr.inet6.Addr[:], v)
	} else {
		err = errors.New("astisrt: ip is neither IPv4 nor IPv6")
		return
	}
	return
}

func newSockaddrFromSockaddrAny(in *syscall.RawSockaddrAny) (out *sockaddr, err error) {
	switch in.Addr.Family {
	case unix.AF_INET6:
		out = &sockaddr{inet6: (*syscall.RawSockaddrInet6)(unsafe.Pointer(in))}
	case unix.AF_INET:
		out = &sockaddr{inet4: (*syscall.RawSockaddrInet4)(unsafe.Pointer(in))}
	default:
		err = errors.New("astisrt: addr is neither IPv4 nor IPv6")
		return
	}
	return
}
