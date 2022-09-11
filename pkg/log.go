package astisrt

// #cgo LDFLAGS: -lsrt
// #include <srt/srt.h>
// #include "log.h"
import "C"
import "sync"

type LogHandler func(l LogLevel, file, area, msg string, line int)

var (
	logHandler      LogHandler
	logHandlerMutex = &sync.Mutex{}
)

type LogLevel int

const (
	LogLevelCritical = LogLevel(C.LOG_CRIT)
	LogLevelDebug    = LogLevel(C.LOG_DEBUG)
	LogLevelError    = LogLevel(C.LOG_ERR)
	LogLevelNotice   = LogLevel(C.LOG_NOTICE)
	LogLevelWarning  = LogLevel(C.LOG_WARNING)
)

// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_setloglevel
func SetLogLevel(l LogLevel) {
	C.srt_setloglevel(C.int(l))
}

//export go2cHandleLog
func go2cHandleLog(ll C.int, file *C.char, line C.int, area, msg *C.char) {
	logHandlerMutex.Lock()
	h := logHandler
	logHandlerMutex.Unlock()
	if h == nil {
		return
	}
	h(LogLevel(ll), C.GoString(file), C.GoString(area), C.GoString(msg), int(line))
}

// https://github.com/Haivision/srt/blob/master/docs/API/API-functions.md#srt_setloghandler
func SetLogHandler(h LogHandler) {
	logHandlerMutex.Lock()
	logHandler = h
	logHandlerMutex.Unlock()
	C.astisrt_setloghandler(boolToCInt(h != nil))
}
