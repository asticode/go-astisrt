package astisrt

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestServerConnectionOptions(t *testing.T) {
	// Create server
	latency := int32(42)
	srv, err := NewServer(ServerOptions{
		ConnectionOptions: []ConnectionOption{WithLatency(latency)},
	})
	require.NoError(t, err)
	defer srv.Close()

	// Assert options
	v, err := srv.s.Options().Latency()
	require.NoError(t, err)
	require.Equal(t, latency, v)
}

func TestServerOnBeforeAccept(t *testing.T) {
	// Quiet logs
	SetLogHandler(func(l LogLevel, file, area, msg string, line int) {})
	defer SetLogHandler(nil)

	// Create server
	srv, err := NewServer(ServerOptions{
		Host: "127.0.0.1",
		OnBeforeAccept: func(c *Connection, version int, streamID string) bool {
			if streamID == "invalid" {
				c.SetPredefinedRejectReason(http.StatusNotFound) //nolint: errcheck
				return false
			}
			return true
		},
		Port: 4000,
	})
	require.NoError(t, err)
	chanListenAndServe := make(chan error)
	defer func() {
		srv.Close()
		<-chanListenAndServe
	}()

	// Listen and serve
	go func() { chanListenAndServe <- srv.ListenAndServe(1) }()

	// Dial with invalid stream id
	c1, err := Dial(DialOptions{
		ConnectionOptions: []ConnectionOption{WithStreamid("invalid")},
		Host:              "127.0.0.1",
		Port:              4000,
	})
	require.ErrorIs(t, err, ErrEconnrej)
	require.Equal(t, RejectReason(1404), c1.RejectReason())
	require.Equal(t, "Application-defined rejection reason", c1.RejectReason().String())

	// Dial with valid stream id
	c2, err := Dial(DialOptions{
		ConnectionOptions: []ConnectionOption{WithStreamid("valid")},
		Host:              "127.0.0.1",
		Port:              4000,
	})
	require.NoError(t, err)
	c2.Close()
}

type ctxKey string

const (
	ctxKeyTest ctxKey = "test"
)

func TestServerHandler(t *testing.T) {
	// Quiet logs
	SetLogHandler(func(l LogLevel, file, area, msg string, line int) {})
	defer SetLogHandler(nil)

	// Create server
	ctxValue := "test"
	chanCtxValue := make(chan string)
	srv, err := NewServer(ServerOptions{
		Handler: ServerHandlerFunc(func(c *Connection) {
			var s string
			if v := c.Context().Value(ctxKeyTest); v != nil {
				s = v.(string)
			}
			chanCtxValue <- s
		}),
		Host: "127.0.0.1",
		OnBeforeAccept: func(c *Connection, version int, streamID string) bool {
			*c = *(c.WithContext(context.WithValue(context.Background(), ctxKeyTest, ctxValue)))
			return true
		},
		Port: 4000,
	})
	require.NoError(t, err)
	chanListenAndServe := make(chan error)
	defer func() {
		srv.Close()
		<-chanListenAndServe
	}()

	// Listen and serve
	go func() { chanListenAndServe <- srv.ListenAndServe(1) }()

	// Dial
	c, err := Dial(DialOptions{
		Host: "127.0.0.1",
		Port: 4000,
	})
	require.NoError(t, err)
	defer c.Close()

	// Assert
	require.Equal(t, ctxValue, <-chanCtxValue)
}

type readResult struct {
	b   []byte
	err error
}

func TestServerReadWriteStats(t *testing.T) {
	// Quiet logs
	SetLogHandler(func(l LogLevel, file, area, msg string, line int) {})
	defer SetLogHandler(nil)

	// Create server
	chanRead := make(chan readResult)
	srv, err := NewServer(ServerOptions{
		Handler: ServerHandlerFunc(func(c *Connection) {
			b := make([]byte, 1500)
			n, err := c.Read(b)
			if err != nil {
				chanRead <- readResult{err: err}
			} else {
				chanRead <- readResult{b: b[:n]}
			}
			<-chanRead
		}),
		Host: "127.0.0.1",
		Port: 4000,
	})
	require.NoError(t, err)
	chanListenAndServe := make(chan error)
	defer func() {
		srv.Close()
		<-chanListenAndServe
	}()

	// Listen and serve
	go func() { chanListenAndServe <- srv.ListenAndServe(1) }()

	// Dial
	c, err := Dial(DialOptions{
		Host: "127.0.0.1",
		Port: 4000,
	})
	require.NoError(t, err)
	defer c.Close()

	// Write
	msg := "test message"
	n, err := c.Write([]byte(msg))
	require.NoError(t, err)
	require.Equal(t, 12, n)

	// Wait for read
	rr := <-chanRead

	// Assert
	require.NoError(t, rr.err)
	require.Equal(t, msg, string(rr.b))

	// Get stats
	s, err := c.Stats(true, false)
	require.NoError(t, err)
	require.True(t, s.ByteSent() > 0)
	require.True(t, s.ByteSentTotal() > 0)
	s, err = c.Stats(true, false)
	require.NoError(t, err)
	require.False(t, s.ByteSent() > 0)

	// Release read
	chanRead <- readResult{}
}

func TestServerShutdownSuccess(t *testing.T) {
	// Quiet logs
	SetLogHandler(func(l LogLevel, file, area, msg string, line int) {})
	defer SetLogHandler(nil)

	// Create server
	var handled bool
	chanHandle := make(chan bool)
	srv, err := NewServer(ServerOptions{
		Handler: ServerHandlerFunc(func(c *Connection) { handled = <-chanHandle }),
		Host:    "127.0.0.1",
		Port:    4000,
	})
	require.NoError(t, err)
	chanListenAndServe := make(chan error)
	defer func() {
		srv.Close()
		<-chanListenAndServe
	}()

	// Listen and serve
	go func() { chanListenAndServe <- srv.ListenAndServe(1) }()

	// Dial
	c, err := Dial(DialOptions{
		Host: "127.0.0.1",
		Port: 4000,
	})
	require.NoError(t, err)
	defer c.Close()

	// Shutdown
	go func() {
		time.Sleep(100 * time.Millisecond)
		chanHandle <- true
	}()
	err = srv.Shutdown(context.Background())
	require.NoError(t, err)
	require.True(t, handled)
}

func TestServerShutdownDeadlineExceeded(t *testing.T) {
	// Quiet logs
	SetLogHandler(func(l LogLevel, file, area, msg string, line int) {})
	defer SetLogHandler(nil)

	// Create server
	chanHandle := make(chan bool)
	srv, err := NewServer(ServerOptions{
		Handler: ServerHandlerFunc(func(c *Connection) { <-chanHandle }),
		Host:    "127.0.0.1",
		Port:    4000,
	})
	require.NoError(t, err)
	chanListenAndServe := make(chan error)
	defer func() {
		srv.Close()
		select {
		case <-chanListenAndServe:
		default:
		}
	}()

	// Listen and serve
	go func() {
		chanListenAndServe <- srv.ListenAndServe(1)
		close(chanListenAndServe)
	}()

	// Dial
	c, err := Dial(DialOptions{
		Host: "127.0.0.1",
		Port: 4000,
	})
	require.NoError(t, err)
	defer c.Close()

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	err = srv.Shutdown(ctx)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	err = <-chanListenAndServe
	require.ErrorIs(t, err, ErrServerClosed)

	// Second ListenAndServe should fail
	require.ErrorIs(t, srv.ListenAndServe(1), ErrServerClosed)

	// Cancel handle
	chanHandle <- true
}
