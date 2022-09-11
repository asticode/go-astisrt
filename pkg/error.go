package astisrt

// #cgo LDFLAGS: -lsrt
// #include <srt/srt.h>
import "C"
import (
	"sync"
	"syscall"
)

// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#diagnostics-1

type Error struct {
	srtErrno int
	sysErrno syscall.Errno
}

func newError(srtErrno, sysErrno C.int) error {
	return Error{
		srtErrno: int(srtErrno),
		sysErrno: syscall.Errno(sysErrno),
	}
}

var (
	ErrEunknown        = newError(C.SRT_EUNKNOWN, 0)
	ErrSuccess         = newError(C.SRT_SUCCESS, 0)
	ErrEconnsetup      = newError(C.SRT_ECONNSETUP, 0)
	ErrEnoserver       = newError(C.SRT_ENOSERVER, 0)
	ErrEconnrej        = newError(C.SRT_ECONNREJ, 0)
	ErrEsockfail       = newError(C.SRT_ESOCKFAIL, 0)
	ErrEsecfail        = newError(C.SRT_ESECFAIL, 0)
	ErrEsclosed        = newError(C.SRT_ESCLOSED, 0)
	ErrEconnfail       = newError(C.SRT_ECONNFAIL, 0)
	ErrEconnlost       = newError(C.SRT_ECONNLOST, 0)
	ErrEnoconn         = newError(C.SRT_ENOCONN, 0)
	ErrEresource       = newError(C.SRT_ERESOURCE, 0)
	ErrEthread         = newError(C.SRT_ETHREAD, 0)
	ErrEnobuf          = newError(C.SRT_ENOBUF, 0)
	ErrEsysobj         = newError(C.SRT_ESYSOBJ, 0)
	ErrEfile           = newError(C.SRT_EFILE, 0)
	ErrEinvrdoff       = newError(C.SRT_EINVRDOFF, 0)
	ErrErdperm         = newError(C.SRT_ERDPERM, 0)
	ErrEinvwroff       = newError(C.SRT_EINVWROFF, 0)
	ErrEwrperm         = newError(C.SRT_EWRPERM, 0)
	ErrEinvop          = newError(C.SRT_EINVOP, 0)
	ErrEboundsock      = newError(C.SRT_EBOUNDSOCK, 0)
	ErrEconnsock       = newError(C.SRT_ECONNSOCK, 0)
	ErrEinvparam       = newError(C.SRT_EINVPARAM, 0)
	ErrEinvsock        = newError(C.SRT_EINVSOCK, 0)
	ErrEunboundsock    = newError(C.SRT_EUNBOUNDSOCK, 0)
	ErrEnolisten       = newError(C.SRT_ENOLISTEN, 0)
	ErrErdvnoserv      = newError(C.SRT_ERDVNOSERV, 0)
	ErrErdvunbound     = newError(C.SRT_ERDVUNBOUND, 0)
	ErrEinvalmsgapi    = newError(C.SRT_EINVALMSGAPI, 0)
	ErrEinvalbufferapi = newError(C.SRT_EINVALBUFFERAPI, 0)
	ErrEduplisten      = newError(C.SRT_EDUPLISTEN, 0)
	ErrElargemsg       = newError(C.SRT_ELARGEMSG, 0)
	ErrEinvpollid      = newError(C.SRT_EINVPOLLID, 0)
	ErrEpollempty      = newError(C.SRT_EPOLLEMPTY, 0)
	ErrEbindconflict   = newError(C.SRT_EBINDCONFLICT, 0)
	ErrEasyncfail      = newError(C.SRT_EASYNCFAIL, 0)
	ErrEasyncsnd       = newError(C.SRT_EASYNCSND, 0)
	ErrEasyncrcv       = newError(C.SRT_EASYNCRCV, 0)
	ErrEtimeout        = newError(C.SRT_ETIMEOUT, 0)
	ErrEcongest        = newError(C.SRT_ECONGEST, 0)
	ErrEpeererr        = newError(C.SRT_EPEERERR, 0)
)

var errorStrMutex = &sync.Mutex{}

// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_strerror
func (err Error) Error() string {
	errorStrMutex.Lock()
	defer errorStrMutex.Unlock()
	return "astisrt: " + C.GoString(C.srt_strerror(C.int(err.srtErrno), C.int(err.sysErrno)))
}

func (err Error) Is(e error) bool {
	a, ok := e.(Error)
	if !ok {
		return false
	}
	return int(a.srtErrno) == int(err.srtErrno)
}

func (err Error) Unwrap() error {
	if err.sysErrno != 0 {
		return err.sysErrno
	}
	return nil
}

// TODO Add Temporary, Timeout, etc. functions
