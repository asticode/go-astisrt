package astisrt

// #cgo darwin LDFLAGS: -L/opt/homebrew/lib
// #cgo darwin CFLAGS: -I/opt/homebrew/include
// #cgo LDFLAGS: -lsrt
// #include <srt/srt.h>
import "C"

// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_startup
func Startup() error { return cStartup() }

// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_cleanup
func CleanUp() error { return cCleanup() }

// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_getversion
func Version() uint32 { return uint32(C.srt_getversion()) }

func boolToCInt(b bool) C.int {
	if b {
		return 1
	}
	return 0
}
