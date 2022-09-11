#include <srt/srt.h>

void astisrt_connect_callback_fn(void* opaque, SRTSOCKET ns, int errorcode, const struct sockaddr* peeraddr, int token);
int astisrt_listen_callback_fn(void* opaque, SRTSOCKET ns, int hs_version, const struct sockaddr* peeraddr, const char* streamid);