package astisrt

import (
	"errors"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	err := newError(_Ctype_int(ErrEconnsetup.(Error).srtErrno), 3)
	require.Equal(t, "astisrt: Connection setup failure", err.Error())
	require.True(t, errors.Is(err, ErrEconnsetup))
	unwrap := errors.Unwrap(err)
	require.NotNil(t, unwrap)
	require.True(t, errors.Is(unwrap, syscall.Errno(3)))
}
