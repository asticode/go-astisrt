package astisrt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreamID(t *testing.T) {
	_, err := ParseStreamID("invalid prefix")
	require.Error(t, err)
	_, err = ParseStreamID("#!:invalid format")
	require.Error(t, err)
	s, err := ParseStreamID("#!::r=r,invalid,u=u")
	require.NoError(t, err)
	require.Equal(t, StreamID{
		StreamIDKeyRessourceName: {Value: "r"},
		StreamIDKeyUserName:      {Value: "u"},
	}, s)
	require.Equal(t, "#!::r=r,u=u", s.String())
	s2 := "#!:{key1={key11={key111=value111},key12=value12},key2=value2,key3={key31=value31}}"
	s, err = ParseStreamID(s2)
	require.NoError(t, err)
	gs2 := StreamID{
		"key1": {Children: map[StreamIDKey]StreamIDItem{
			"key11": {Children: map[StreamIDKey]StreamIDItem{"key111": {Value: "value111"}}},
			"key12": {Value: "value12"},
		}},
		"key2": {Value: "value2"},
		"key3": {Children: map[StreamIDKey]StreamIDItem{"key31": {Value: "value31"}}},
	}
	require.Equal(t, gs2, s)
	require.Equal(t, s2, s.String())
	s = make(StreamID)
	err = s.UnmarshalText([]byte(s2))
	require.NoError(t, err)
	require.Equal(t, gs2, s)
	b, err := s.MarshalText()
	require.NoError(t, err)
	require.Equal(t, s2, string(b))
}
