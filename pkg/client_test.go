package astisrt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type disconnectResult struct {
	c   *Connection
	err error
}

func TestDial(t *testing.T) {
	// Quiet logs
	SetLogHandler(func(l LogLevel, file, area, msg string, line int) {})
	defer SetLogHandler(nil)

	// Create server
	srv, err := NewServer(ServerOptions{
		Handler: ServerHandlerFunc(func(c *Connection) { c.Close() }),
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
	chanOnDisconnect := make(chan disconnectResult)
	c, err := Dial(DialOptions{
		ConnectionOptions: []ConnectionOption{WithStreamid("streamid")},
		Host:              "127.0.0.1",
		OnDisconnect:      func(c *Connection, err error) { chanOnDisconnect <- disconnectResult{c: c, err: err} },
		Port:              4000,
	})
	require.NoError(t, err)
	defer c.Close()

	// Assert dial
	require.Equal(t, "127.0.0.1:4000", c.Addr().String())
	s, err := c.Options().Streamid()
	require.NoError(t, err)
	require.Equal(t, "streamid", s)
	dr := <-chanOnDisconnect
	require.Equal(t, c, dr.c)
	require.ErrorIs(t, dr.err, ErrEconnlost)
}
