package astisrt

// #cgo LDFLAGS: -lsrt
// #include <srt/srt.h>
import "C"
import (
	"context"
	"net"
)

type Connection struct {
	addr *net.UDPAddr
	ctx  context.Context
	s    *Socket
}

func newConnection(addr *net.UDPAddr, s *Socket) *Connection {
	return &Connection{
		addr: addr,
		ctx:  context.Background(),
		s:    s,
	}
}

func (c *Connection) Addr() *net.UDPAddr {
	return c.addr
}

func (c *Connection) Context() context.Context {
	return c.ctx
}

func (c *Connection) WithContext(ctx context.Context) *Connection {
	c.ctx = ctx
	return c
}

func (c *Connection) Close() error {
	return c.s.Close()
}

func (c *Connection) Read(b []byte) (int, error) {
	return c.s.ReceiveMessage(b)
}

func (c *Connection) Write(b []byte) (int, error) {
	return c.s.SendMessage(b)
}

func (c *Connection) RejectReason() RejectReason {
	return c.s.RejectReason()
}

// HTTP status must be [0:1000)
// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_setrejectreason
func (c *Connection) SetPredefinedRejectReason(httpStatus int) error {
	return c.s.SetRejectReason(RejectReason(C.SRT_REJC_PREDEFINED + C.int(httpStatus)))
}

// Status must be [0:1000)
// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_setrejectreason
func (c *Connection) SetUserDefinedRejectReason(status int) error {
	return c.s.SetRejectReason(RejectReason(C.SRT_REJC_USERDEFINED + C.int(status)))
}

func (c *Connection) Stats(clear, instantaneous bool) (Stats, error) {
	return c.s.Stats(clear, instantaneous)
}
