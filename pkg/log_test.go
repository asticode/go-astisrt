package astisrt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLog(t *testing.T) {
	var msgs []string
	SetLogHandler(func(l LogLevel, file, area, msg string, line int) { msgs = append(msgs, msg) })
	defer SetLogHandler(nil)

	SetLogLevel(LogLevelDebug)
	s1, err := NewSocket()
	require.NoError(t, err)
	defer s1.Close()
	require.True(t, len(msgs) > 0)

	msgs = []string{}
	SetLogLevel(LogLevelWarning)
	s2, err := NewSocket()
	require.NoError(t, err)
	defer s2.Close()
	require.True(t, len(msgs) == 0)
}
