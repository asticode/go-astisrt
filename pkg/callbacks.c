#include "_cgo_export.h"
#include <srt/srt.h>

void astisrt_connect_callback_fn(void* opaque, SRTSOCKET ns, int errorcode, const struct sockaddr* peeraddr, int token) {
    return go2cConnectCallback(opaque, ns, errorcode, (struct sockaddr*) peeraddr, token);
}

int astisrt_listen_callback_fn(void* opaque, SRTSOCKET ns, int hs_version, const struct sockaddr* peeraddr, const char* streamid) {
    return go2cListenCallback(opaque, ns, hs_version, (struct sockaddr*) peeraddr, (char*) streamid);
}