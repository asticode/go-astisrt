FROM debian:latest

RUN apt-get update

##
# git
##

RUN apt-get install -y git nano

##
# srt
##

ARG SRT_VERSION
                
RUN apt-get install -y \
    build-essential \
    tclsh \
    cmake \
    pkg-config \
    libssl-dev

RUN mkdir -p /opt/srt \
    && git clone https://github.com/Haivision/srt /opt/srt/src \
    && cd /opt/srt/src \
    && git checkout ${SRT_VERSION}

RUN cd /opt/srt/src \
    && ./configure --prefix=.. \
    && make \
    && make install

##
# go
##

ARG GO_VERSION

RUN apt-get install -y wget
RUN <<_EOF_ sh
arch=$(dpkg --print-architecture)
goArch="amd64"
case \${arch} in
    arm64) goArch="arm64" ;;
esac
wget -O /tmp/go.tar.gz https://dl.google.com/go/go${GO_VERSION}.linux-\${goArch}.tar.gz
tar -C /opt -xzf /tmp/go.tar.gz
_EOF_
ENV PATH="$PATH:/opt/go/bin"