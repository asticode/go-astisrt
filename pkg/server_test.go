package astisrt

import (
	"context"
	"net/http"
	"testing"

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
