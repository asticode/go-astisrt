package astisrt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO Why is this test taking 5s ?
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
	chanOnDisconnect := make(chan *Connection)
	c, err := Dial(DialOptions{
		ConnectionOptions: []ConnectionOption{WithStreamid("streamid")},
		Host:              "127.0.0.1",
		OnDisconnect:      func(c *Connection, err error) { chanOnDisconnect <- c },
		Port:              4000,
	})
	require.NoError(t, err)
	defer c.Close()

	// Assert dial
	require.Equal(t, "127.0.0.1:4000", c.Addr().String())
	s, err := c.Options().Streamid()
	require.NoError(t, err)
	require.Equal(t, "streamid", s)
	require.Equal(t, c, <-chanOnDisconnect)
}
