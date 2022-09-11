package astisrt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPortConversion(t *testing.T) {
	e1, e2 := uint16(4000), uint16(40975)
	require.Equal(t, e2, sockAddrPortFromRegularPort(e1))
	require.Equal(t, e1, regularPortFromSockAddrPort(e2))
}
