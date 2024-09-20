package astisrt

import (
	"fmt"
	"net"
)

type DialOnDisconnect func(c *Connection, err error)

func (f DialOnDisconnect) connectCallback(c *Connection) ConnectCallback {
	return func(s *Socket, addr *net.UDPAddr, token int, err error) { f(c, err) }
}

type DialOptions struct {
	// Connection options
	ConnectionOptions []ConnectionOption

	// Host to connect to
	Host string

	// Callback executed when connection disconnects.
	OnDisconnect DialOnDisconnect

	// Port to connect to
	Port uint16
}

func Dial(o DialOptions) (c *Connection, err error) {
	// Create socket
	var s *Socket
	if s, err = NewSocket(); err != nil {
		err = fmt.Errorf("astisrt: creating socket failed: %w", err)
		return
	}

	// Make sure socket is closed in case of error
	defer func() {
		if err != nil {
			s.Close()
		}
	}()

	// Create connection
	c = newConnection(nil, s)

	// Apply connection options
	if err = applyConnectionOptions(s, o.ConnectionOptions); err != nil {
		err = fmt.Errorf("astisrt: applying connection options failed: %w", err)
		return
	}

	// Set callbacks
	if o.OnDisconnect != nil {
		if err = s.SetConnectCallback(o.OnDisconnect.connectCallback(c)); err != nil {
			err = fmt.Errorf("astisrt: setting connect callback failed: %w", err)
			return
		}
	}

	// Connect
	if c.addr, err = s.Connect(o.Host, o.Port); err != nil {
		err = fmt.Errorf("astisrt: connecting failed: %w", err)
		return
	}
	return
}
