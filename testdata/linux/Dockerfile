FROM golang:1.19

RUN \
  apt-get update && \
  DEBIAN_FRONTEND=noninteractive apt-get install -y \
    tclsh pkg-config cmake libssl-dev build-essential git \
  && apt-get clean

ENV SRT_VERSION="v1.5.3"

RUN \
  cd /opt && \
  git clone https://github.com/asticode/go-astisrt.git && \
  cd go-astisrt && \
  make install-srt version="${SRT_VERSION}" && \
  mv tmp/${SRT_VERSION} /opt/srt && \
  cd .. && \
  rm -rf go-astisrt

ENV LD_LIBRARY_PATH="$LD_LIBRARY_PATH:/opt/srt/lib/"
ENV CGO_CFLAGS="-I/opt/srt/include/"
ENV CGO_LDFLAGS="-L/opt/srt/lib/"
ENV PKG_CONFIG_PATH="/opt/srt/lib/pkgconfig"
ENV GOCACHE="/opt/astisrt/tmp/linux/gocache"
ENV GOPATH="/opt/astisrt/tmp/linux/gopath"

CMD \
  cd /opt/astisrt/tmp/linux/gopath/src/github.com/asticode/go-astisrt/pkg && \
  go test