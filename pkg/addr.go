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

func sockAddrFromHostAndPort(host string, port uint16) (addr *C.struct_sockaddr, size C.int, err error) {
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

	// Create addr
	if v := ip.To4(); v != nil {
		// Create raw
		raw := syscall.RawSockaddrInet4{
			Family: unix.AF_INET,
			Port:   sockAddrPortFromRegularPort(port),
		}

		// Copy addr
		copy(raw.Addr[:], v)

		// Create addr
		addr = (*C.struct_sockaddr)(unsafe.Pointer(&raw))
		size = syscall.SizeofSockaddrInet4
	} else if v := ip.To16(); v != nil {
		// Create raw
		raw := syscall.RawSockaddrInet6{
			Family: unix.AF_INET6,
			Port:   sockAddrPortFromRegularPort(port),
		}

		// Copy addr
		copy(raw.Addr[:], v)

		// Create addr
		addr = (*C.struct_sockaddr)(unsafe.Pointer(&raw))
		size = syscall.SizeofSockaddrInet6
	} else {
		err = errors.New("astisrt: ip is neither IPv4 nor IPv6")
		return
	}
	return
}

func udpAddrFromSockaddr(rawAny *syscall.RawSockaddrAny) (addr *net.UDPAddr, err error) {
	switch rawAny.Addr.Family {
	case unix.AF_INET6:
		raw6 := (*syscall.RawSockaddrInet6)(unsafe.Pointer(rawAny))
		addr = &net.UDPAddr{
			IP:   raw6.Addr[:],
			Port: int(regularPortFromSockAddrPort(raw6.Port)),
		}
	case unix.AF_INET:
		raw4 := (*syscall.RawSockaddrInet4)(unsafe.Pointer(rawAny))
		addr = &net.UDPAddr{
			IP:   raw4.Addr[:],
			Port: int(regularPortFromSockAddrPort(raw4.Port)),
		}
	default:
		err = errors.New("astisrt: addr is neither IPv4 nor IPv6")
		return
	}
	return
}
