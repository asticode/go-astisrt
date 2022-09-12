// Code generated by astisrt using internal/cmd/options. DO NOT EDIT.
package astisrt

import (
	"fmt"
	"syscall"
)

type ConnectionOptions struct {
	c *Connection
}

func (c *Connection) Options() *ConnectionOptions {
	return &ConnectionOptions{c: c}
}

type ConnectionOption struct {
	do 		 func(s *Socket) error
	name 	 string
}

func applyConnectionOptions(s *Socket, opts []ConnectionOption) (err error) {
	for _, opt := range opts {
		if err = opt.do(s); err != nil {
			err = fmt.Errorf("astisrt: applying %s option failed: %w", opt.name, err)
			return
		}
	}
	return
}

func (co *ConnectionOptions) Bindtodevice() (string, error) {
	return co.c.s.Options().Bindtodevice()
}

func (co *ConnectionOptions) SetBindtodevice(v string) error {
	return co.c.s.Options().SetBindtodevice(v)
}

func WithBindtodevice(v string) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetBindtodevice(v) },
		name: "Bindtodevice",
	}
}

func (co *ConnectionOptions) SetCongestion(v string) error {
	return co.c.s.Options().SetCongestion(v)
}

func WithCongestion(v string) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetCongestion(v) },
		name: "Congestion",
	}
}

func (co *ConnectionOptions) SetConntimeo(v int32) error {
	return co.c.s.Options().SetConntimeo(v)
}

func WithConntimeo(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetConntimeo(v) },
		name: "Conntimeo",
	}
}

func (co *ConnectionOptions) Drifttracer() (bool, error) {
	return co.c.s.Options().Drifttracer()
}

func (co *ConnectionOptions) SetDrifttracer(v bool) error {
	return co.c.s.Options().SetDrifttracer(v)
}

func WithDrifttracer(v bool) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetDrifttracer(v) },
		name: "Drifttracer",
	}
}

func (co *ConnectionOptions) SetEnforcedencryption(v bool) error {
	return co.c.s.Options().SetEnforcedencryption(v)
}

func WithEnforcedencryption(v bool) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetEnforcedencryption(v) },
		name: "Enforcedencryption",
	}
}

func (co *ConnectionOptions) Event() (int32, error) {
	return co.c.s.Options().Event()
}

func (co *ConnectionOptions) Fc() (int32, error) {
	return co.c.s.Options().Fc()
}

func (co *ConnectionOptions) SetFc(v int32) error {
	return co.c.s.Options().SetFc(v)
}

func WithFc(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetFc(v) },
		name: "Fc",
	}
}

func (co *ConnectionOptions) SetGroupconnect(v int32) error {
	return co.c.s.Options().SetGroupconnect(v)
}

func WithGroupconnect(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetGroupconnect(v) },
		name: "Groupconnect",
	}
}

func (co *ConnectionOptions) SetGroupminstabletimeo(v int32) error {
	return co.c.s.Options().SetGroupminstabletimeo(v)
}

func WithGroupminstabletimeo(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetGroupminstabletimeo(v) },
		name: "Groupminstabletimeo",
	}
}

func (co *ConnectionOptions) Grouptype() (int32, error) {
	return co.c.s.Options().Grouptype()
}

func (co *ConnectionOptions) Inputbw() (int64, error) {
	return co.c.s.Options().Inputbw()
}

func (co *ConnectionOptions) SetInputbw(v int64) error {
	return co.c.s.Options().SetInputbw(v)
}

func WithInputbw(v int64) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetInputbw(v) },
		name: "Inputbw",
	}
}

func (co *ConnectionOptions) Iptos() (int32, error) {
	return co.c.s.Options().Iptos()
}

func (co *ConnectionOptions) SetIptos(v int32) error {
	return co.c.s.Options().SetIptos(v)
}

func WithIptos(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetIptos(v) },
		name: "Iptos",
	}
}

func (co *ConnectionOptions) Ipttl() (int32, error) {
	return co.c.s.Options().Ipttl()
}

func (co *ConnectionOptions) SetIpttl(v int32) error {
	return co.c.s.Options().SetIpttl(v)
}

func WithIpttl(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetIpttl(v) },
		name: "Ipttl",
	}
}

func (co *ConnectionOptions) Ipv6only() (int32, error) {
	return co.c.s.Options().Ipv6only()
}

func (co *ConnectionOptions) SetIpv6only(v int32) error {
	return co.c.s.Options().SetIpv6only(v)
}

func WithIpv6only(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetIpv6only(v) },
		name: "Ipv6only",
	}
}

func (co *ConnectionOptions) Isn() (int32, error) {
	return co.c.s.Options().Isn()
}

func (co *ConnectionOptions) Kmpreannounce() (int32, error) {
	return co.c.s.Options().Kmpreannounce()
}

func (co *ConnectionOptions) SetKmpreannounce(v int32) error {
	return co.c.s.Options().SetKmpreannounce(v)
}

func WithKmpreannounce(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetKmpreannounce(v) },
		name: "Kmpreannounce",
	}
}

func (co *ConnectionOptions) Kmrefreshrate() (int32, error) {
	return co.c.s.Options().Kmrefreshrate()
}

func (co *ConnectionOptions) SetKmrefreshrate(v int32) error {
	return co.c.s.Options().SetKmrefreshrate(v)
}

func WithKmrefreshrate(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetKmrefreshrate(v) },
		name: "Kmrefreshrate",
	}
}

func (co *ConnectionOptions) Kmstate() (KmState, error) {
	return co.c.s.Options().Kmstate()
}

func (co *ConnectionOptions) Latency() (int32, error) {
	return co.c.s.Options().Latency()
}

func (co *ConnectionOptions) SetLatency(v int32) error {
	return co.c.s.Options().SetLatency(v)
}

func WithLatency(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetLatency(v) },
		name: "Latency",
	}
}

func (co *ConnectionOptions) Linger() (syscall.Linger, error) {
	return co.c.s.Options().Linger()
}

func (co *ConnectionOptions) SetLinger(v syscall.Linger) error {
	return co.c.s.Options().SetLinger(v)
}

func WithLinger(v syscall.Linger) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetLinger(v) },
		name: "Linger",
	}
}

func (co *ConnectionOptions) Lossmaxttl() (int32, error) {
	return co.c.s.Options().Lossmaxttl()
}

func (co *ConnectionOptions) SetLossmaxttl(v int32) error {
	return co.c.s.Options().SetLossmaxttl(v)
}

func WithLossmaxttl(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetLossmaxttl(v) },
		name: "Lossmaxttl",
	}
}

func (co *ConnectionOptions) Maxbw() (int64, error) {
	return co.c.s.Options().Maxbw()
}

func (co *ConnectionOptions) SetMaxbw(v int64) error {
	return co.c.s.Options().SetMaxbw(v)
}

func WithMaxbw(v int64) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetMaxbw(v) },
		name: "Maxbw",
	}
}

func (co *ConnectionOptions) SetMessageapi(v bool) error {
	return co.c.s.Options().SetMessageapi(v)
}

func WithMessageapi(v bool) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetMessageapi(v) },
		name: "Messageapi",
	}
}

func (co *ConnectionOptions) Mininputbw() (int64, error) {
	return co.c.s.Options().Mininputbw()
}

func (co *ConnectionOptions) SetMininputbw(v int64) error {
	return co.c.s.Options().SetMininputbw(v)
}

func WithMininputbw(v int64) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetMininputbw(v) },
		name: "Mininputbw",
	}
}

func (co *ConnectionOptions) Minversion() (int32, error) {
	return co.c.s.Options().Minversion()
}

func (co *ConnectionOptions) SetMinversion(v int32) error {
	return co.c.s.Options().SetMinversion(v)
}

func WithMinversion(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetMinversion(v) },
		name: "Minversion",
	}
}

func (co *ConnectionOptions) Mss() (int32, error) {
	return co.c.s.Options().Mss()
}

func (co *ConnectionOptions) SetMss(v int32) error {
	return co.c.s.Options().SetMss(v)
}

func WithMss(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetMss(v) },
		name: "Mss",
	}
}

func (co *ConnectionOptions) Nakreport() (bool, error) {
	return co.c.s.Options().Nakreport()
}

func (co *ConnectionOptions) SetNakreport(v bool) error {
	return co.c.s.Options().SetNakreport(v)
}

func WithNakreport(v bool) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetNakreport(v) },
		name: "Nakreport",
	}
}

func (co *ConnectionOptions) Oheadbw() (int32, error) {
	return co.c.s.Options().Oheadbw()
}

func (co *ConnectionOptions) SetOheadbw(v int32) error {
	return co.c.s.Options().SetOheadbw(v)
}

func WithOheadbw(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetOheadbw(v) },
		name: "Oheadbw",
	}
}

func (co *ConnectionOptions) Packetfilter() (string, error) {
	return co.c.s.Options().Packetfilter()
}

func (co *ConnectionOptions) SetPacketfilter(v string) error {
	return co.c.s.Options().SetPacketfilter(v)
}

func WithPacketfilter(v string) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetPacketfilter(v) },
		name: "Packetfilter",
	}
}

func (co *ConnectionOptions) SetPassphrase(v string) error {
	return co.c.s.Options().SetPassphrase(v)
}

func WithPassphrase(v string) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetPassphrase(v) },
		name: "Passphrase",
	}
}

func (co *ConnectionOptions) SetPayloadsize(v int32) error {
	return co.c.s.Options().SetPayloadsize(v)
}

func WithPayloadsize(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetPayloadsize(v) },
		name: "Payloadsize",
	}
}

func (co *ConnectionOptions) Pbkeylen() (int32, error) {
	return co.c.s.Options().Pbkeylen()
}

func (co *ConnectionOptions) SetPbkeylen(v int32) error {
	return co.c.s.Options().SetPbkeylen(v)
}

func WithPbkeylen(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetPbkeylen(v) },
		name: "Pbkeylen",
	}
}

func (co *ConnectionOptions) Peeridletimeo() (int32, error) {
	return co.c.s.Options().Peeridletimeo()
}

func (co *ConnectionOptions) SetPeeridletimeo(v int32) error {
	return co.c.s.Options().SetPeeridletimeo(v)
}

func WithPeeridletimeo(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetPeeridletimeo(v) },
		name: "Peeridletimeo",
	}
}

func (co *ConnectionOptions) Peerlatency() (int32, error) {
	return co.c.s.Options().Peerlatency()
}

func (co *ConnectionOptions) SetPeerlatency(v int32) error {
	return co.c.s.Options().SetPeerlatency(v)
}

func WithPeerlatency(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetPeerlatency(v) },
		name: "Peerlatency",
	}
}

func (co *ConnectionOptions) Peerversion() (int32, error) {
	return co.c.s.Options().Peerversion()
}

func (co *ConnectionOptions) Rcvbuf() (int32, error) {
	return co.c.s.Options().Rcvbuf()
}

func (co *ConnectionOptions) SetRcvbuf(v int32) error {
	return co.c.s.Options().SetRcvbuf(v)
}

func WithRcvbuf(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetRcvbuf(v) },
		name: "Rcvbuf",
	}
}

func (co *ConnectionOptions) Rcvdata() (int32, error) {
	return co.c.s.Options().Rcvdata()
}

func (co *ConnectionOptions) Rcvkmstate() (int32, error) {
	return co.c.s.Options().Rcvkmstate()
}

func (co *ConnectionOptions) Rcvlatency() (int32, error) {
	return co.c.s.Options().Rcvlatency()
}

func (co *ConnectionOptions) SetRcvlatency(v int32) error {
	return co.c.s.Options().SetRcvlatency(v)
}

func WithRcvlatency(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetRcvlatency(v) },
		name: "Rcvlatency",
	}
}

func (co *ConnectionOptions) Rcvtimeo() (int32, error) {
	return co.c.s.Options().Rcvtimeo()
}

func (co *ConnectionOptions) SetRcvtimeo(v int32) error {
	return co.c.s.Options().SetRcvtimeo(v)
}

func WithRcvtimeo(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetRcvtimeo(v) },
		name: "Rcvtimeo",
	}
}

func (co *ConnectionOptions) Rendezvous() (bool, error) {
	return co.c.s.Options().Rendezvous()
}

func (co *ConnectionOptions) SetRendezvous(v bool) error {
	return co.c.s.Options().SetRendezvous(v)
}

func WithRendezvous(v bool) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetRendezvous(v) },
		name: "Rendezvous",
	}
}

func (co *ConnectionOptions) Retransmitalgo() (int32, error) {
	return co.c.s.Options().Retransmitalgo()
}

func (co *ConnectionOptions) SetRetransmitalgo(v int32) error {
	return co.c.s.Options().SetRetransmitalgo(v)
}

func WithRetransmitalgo(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetRetransmitalgo(v) },
		name: "Retransmitalgo",
	}
}

func (co *ConnectionOptions) Reuseaddr() (bool, error) {
	return co.c.s.Options().Reuseaddr()
}

func (co *ConnectionOptions) SetReuseaddr(v bool) error {
	return co.c.s.Options().SetReuseaddr(v)
}

func WithReuseaddr(v bool) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetReuseaddr(v) },
		name: "Reuseaddr",
	}
}

func (co *ConnectionOptions) SetSender(v bool) error {
	return co.c.s.Options().SetSender(v)
}

func WithSender(v bool) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetSender(v) },
		name: "Sender",
	}
}

func (co *ConnectionOptions) Sndbuf() (int32, error) {
	return co.c.s.Options().Sndbuf()
}

func (co *ConnectionOptions) SetSndbuf(v int32) error {
	return co.c.s.Options().SetSndbuf(v)
}

func WithSndbuf(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetSndbuf(v) },
		name: "Sndbuf",
	}
}

func (co *ConnectionOptions) Snddata() (int32, error) {
	return co.c.s.Options().Snddata()
}

func (co *ConnectionOptions) SetSnddropdelay(v int32) error {
	return co.c.s.Options().SetSnddropdelay(v)
}

func WithSnddropdelay(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetSnddropdelay(v) },
		name: "Snddropdelay",
	}
}

func (co *ConnectionOptions) Sndkmstate() (int32, error) {
	return co.c.s.Options().Sndkmstate()
}

func (co *ConnectionOptions) Sndtimeo() (int32, error) {
	return co.c.s.Options().Sndtimeo()
}

func (co *ConnectionOptions) SetSndtimeo(v int32) error {
	return co.c.s.Options().SetSndtimeo(v)
}

func WithSndtimeo(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetSndtimeo(v) },
		name: "Sndtimeo",
	}
}

func (co *ConnectionOptions) State() (int32, error) {
	return co.c.s.Options().State()
}

func (co *ConnectionOptions) Streamid() (string, error) {
	return co.c.s.Options().Streamid()
}

func (co *ConnectionOptions) SetStreamid(v string) error {
	return co.c.s.Options().SetStreamid(v)
}

func WithStreamid(v string) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetStreamid(v) },
		name: "Streamid",
	}
}

func (co *ConnectionOptions) Tlpktdrop() (bool, error) {
	return co.c.s.Options().Tlpktdrop()
}

func (co *ConnectionOptions) SetTlpktdrop(v bool) error {
	return co.c.s.Options().SetTlpktdrop(v)
}

func WithTlpktdrop(v bool) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetTlpktdrop(v) },
		name: "Tlpktdrop",
	}
}

func (co *ConnectionOptions) SetTranstype(v Transtype) error {
	return co.c.s.Options().SetTranstype(v)
}

func WithTranstype(v Transtype) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetTranstype(v) },
		name: "Transtype",
	}
}

func (co *ConnectionOptions) SetTsbpdmode(v bool) error {
	return co.c.s.Options().SetTsbpdmode(v)
}

func WithTsbpdmode(v bool) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetTsbpdmode(v) },
		name: "Tsbpdmode",
	}
}

func (co *ConnectionOptions) UdpRcvbuf() (int32, error) {
	return co.c.s.Options().UdpRcvbuf()
}

func (co *ConnectionOptions) SetUdpRcvbuf(v int32) error {
	return co.c.s.Options().SetUdpRcvbuf(v)
}

func WithUdpRcvbuf(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetUdpRcvbuf(v) },
		name: "UdpRcvbuf",
	}
}

func (co *ConnectionOptions) UdpSndbuf() (int32, error) {
	return co.c.s.Options().UdpSndbuf()
}

func (co *ConnectionOptions) SetUdpSndbuf(v int32) error {
	return co.c.s.Options().SetUdpSndbuf(v)
}

func WithUdpSndbuf(v int32) ConnectionOption {
	return ConnectionOption{
		do: func(s *Socket) error { return s.Options().SetUdpSndbuf(v) },
		name: "UdpSndbuf",
	}
}

func (co *ConnectionOptions) Version() (int32, error) {
	return co.c.s.Options().Version()
}