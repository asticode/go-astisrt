package astisrt

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
)

var (
	ErrServerClosed = errors.New("astisrt: server closed")
)

type ServerHandlerFunc func(c *Connection)

func (f ServerHandlerFunc) ServeSRT(c *Connection) { f(c) }

type ServerHandler interface {
	ServeSRT(c *Connection)
}

var NopServerHandler = nopServerHandler{}

type nopServerHandler struct{}

func (h nopServerHandler) ServeSRT(c *Connection) {}

type Server struct {
	accepted     map[*Socket]*Connection
	closed       bool
	h            ServerHandler
	handling     map[*Connection]bool
	host         string
	m            *sync.Mutex
	port         uint16
	s            *Socket
	shutdownChan chan bool
}

type ServerOnBeforeAccept func(c *Connection, version int, streamID string) bool

type ServerOptions struct {
	// Connection options
	ConnectionOptions []ConnectionOption

	// Handler to invoke on accepted connections
	// Defaults to NopServerHandler
	Handler ServerHandler

	// Host to bind with
	Host string

	// Callback executed before accepting a new connection.
	// Return true if connection should be accepted, false otherwise.
	OnBeforeAccept ServerOnBeforeAccept

	// Port to bind with
	Port uint16
}

func (srv *Server) listenCallback(cb ServerOnBeforeAccept) ListenCallback {
	return func(s *Socket, version int, addr *net.UDPAddr, streamID string) bool {
		// Create connection
		c := newConnection(addr, s)

		// Callback
		if ok := cb(c, version, streamID); !ok {
			return false
		}

		// Store connection
		srv.m.Lock()
		srv.accepted[s] = c
		srv.m.Unlock()
		return true
	}
}

func NewServer(o ServerOptions) (srv *Server, err error) {
	// Create server
	srv = &Server{
		accepted: make(map[*Socket]*Connection),
		h:        NopServerHandler,
		handling: make(map[*Connection]bool),
		host:     o.Host,
		m:        &sync.Mutex{},
		port:     o.Port,
	}

	// Update handler
	if o.Handler != nil {
		srv.h = o.Handler
	}

	// Create socket
	if srv.s, err = NewSocket(); err != nil {
		err = fmt.Errorf("astisrt: creating socket failed: %w", err)
		return
	}

	// Apply connection options
	if err = applyConnectionOptions(srv.s, o.ConnectionOptions); err != nil {
		err = fmt.Errorf("astisrt: applying connection options failed: %w", err)
		return
	}

	// Set callbacks
	if o.OnBeforeAccept != nil {
		if err = srv.s.SetListenCallback(srv.listenCallback(o.OnBeforeAccept)); err != nil {
			err = fmt.Errorf("astisrt: setting listen callback failed: %w", err)
			return
		}
	}
	return
}

func (srv *Server) close(shutdown bool) (err error) {
	// Lock
	srv.m.Lock()
	defer srv.m.Unlock()

	// Server is already closed
	if srv.closed {
		return nil
	}

	// Server is closed
	srv.closed = true

	// Create shutdown chan
	if shutdown && len(srv.handling) > 0 {
		srv.shutdownChan = make(chan bool)
	}

	// Close connections being handled
	for c := range srv.handling {
		if err = c.Close(); err != nil {
			err = fmt.Errorf("astisrt: closing connection being handled failed: %w", err)
			return
		}
	}

	// Close server socket
	if err = srv.s.Close(); err != nil {
		err = fmt.Errorf("astisrt: closing server socket failed: %w", err)
		return
	}
	return
}

func (srv *Server) Close() error {
	return srv.close(false)
}

func (srv *Server) ListenAndServe(backlog int) (err error) {
	// Server is closed
	srv.m.Lock()
	if srv.closed {
		srv.m.Unlock()
		return ErrServerClosed
	}
	srv.m.Unlock()

	// Bind
	if err = srv.s.Bind(srv.host, srv.port); err != nil {
		err = fmt.Errorf("astisrt: binding failed: %w", err)
		return
	}

	// Listen
	if err = srv.s.Listen(backlog); err != nil {
		err = fmt.Errorf("astisrt: listening failed: %w", err)
		return
	}

	// Loop
	for {
		// Accept
		var s *Socket
		var addr *net.UDPAddr
		if s, addr, err = srv.s.Accept(); err != nil {
			// Server is closed
			srv.m.Lock()
			if srv.closed {
				srv.m.Unlock()
				return ErrServerClosed
			}
			srv.m.Unlock()

			// TODO Retry on temporary error (e.g. http.Server)

			// Fatal error
			err = fmt.Errorf("astisrt: accepting failed: %w", err)
			return
		}

		// Get connection
		srv.m.Lock()
		c, ok := srv.accepted[s]
		if !ok {
			c = newConnection(addr, s)
		} else {
			delete(srv.accepted, s)
		}
		srv.m.Unlock()

		// Handle
		go srv.handleAcceptedConnection(c)
	}
}

func (srv *Server) handleAcceptedConnection(c *Connection) {
	// Store connection
	srv.m.Lock()
	if srv.closed {
		srv.m.Unlock()
		return
	}
	srv.handling[c] = true
	srv.m.Unlock()

	// Make sure connection is properly closed when handling is done
	defer func() {
		// Lock
		srv.m.Lock()
		defer srv.m.Unlock()

		// Close connection
		c.Close()

		// Remove connection
		delete(srv.handling, c)

		// Close shutdown chan
		if srv.shutdownChan != nil && len(srv.handling) == 0 {
			close(srv.shutdownChan)
		}
	}()

	// Serve
	srv.h.ServeSRT(c)
}

func (srv *Server) Shutdown(ctx context.Context) (err error) {
	// Close
	if err = srv.close(true); err != nil {
		err = fmt.Errorf("astisrt: closing failed: %w", err)
		return
	}

	// Get shutdown chan
	srv.m.Lock()
	c := srv.shutdownChan
	if c == nil {
		srv.m.Unlock()
		return
	}
	srv.m.Unlock()

	// Wait for either context to be done or shutdown chan to be closed
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c:
		return
	}
}
