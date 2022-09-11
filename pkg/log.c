#include "_cgo_export.h"
#include <srt/srt.h>

void astisrt_handlelog(void* opaque, int level, const char* file, int line, const char* area, const char* message) {
    go2cHandleLog(level, (char*) file, line, (char*) area, (char*) message);
}

void astisrt_setloghandler(int b) {
    if (b == 1) {
        srt_setloghandler(NULL, astisrt_handlelog);
    } else {
        srt_setloghandler(NULL, NULL);
    }
}