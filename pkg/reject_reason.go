package astisrt

// #cgo LDFLAGS: -lsrt
// #include <srt/srt.h>
import "C"

type RejectReason int

// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_rejectreason_str
func (r RejectReason) String() string {
	return C.GoString(C.srt_rejectreason_str(C.int(r)))
}

// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_getrejectreason
func (s *Socket) RejectReason() RejectReason {
	return s.rr
}

// Reject reason must be [1000:3000)
// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_setrejectreason
func (s *Socket) SetRejectReason(r RejectReason) error {
	return cSetRejectReason(s.u, C.int(r))
}
