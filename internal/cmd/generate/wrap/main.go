package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type fnArg struct {
	Format string
	Name   string
	Type   string
}

type fn struct {
	CArgs               []fnArg
	CName               string
	RejectReasonArgName string
	GoArgs              []fnArg
	GoName              string
	Return              string
}

func (f fn) joinArgs(args []fnArg) string {
	var ss []string
	for _, a := range args {
		s := a.Name
		if a.Format != "" {
			s = fmt.Sprintf(a.Format, a.Name)
		}
		ss = append(ss, s)
	}
	return strings.Join(ss, ", ")
}

func (f fn) JoinCArgs() string {
	return f.joinArgs(f.CArgs)
}

func (f fn) JoinGoArgs() string {
	s := f.joinArgs(f.GoArgs)
	if len(f.GoArgs) > 0 {
		s += ", "
	}
	return s
}

func (f fn) JoinCArgsWithType() (s string) {
	for _, a := range f.CArgs {
		if a.Type == "" {
			continue
		}
		s += a.Type + " " + a.Name + ", "
	}
	return
}

func (f fn) JoinGoArgsWithType() (s string) {
	var ss []string
	for _, a := range f.GoArgs {
		if a.Type == "" {
			continue
		}
		ss = append(ss, a.Name+" "+a.Type)
	}
	return strings.Join(ss, ", ")
}

func (f fn) CRetCompare() string {
	if f.Return == "SRTSOCKET" {
		return "SRT_INVALID_SOCK"
	}
	return "SRT_ERROR"
}

func (f fn) GoRetCompare() string {
	if f.Return == "SRTSOCKET" {
		return "SRT_INVALID_SOCK_"
	}
	return "SRT_ERROR_"
}

func (f fn) CReturn() string {
	if f.Return != "" {
		return f.Return
	}
	return "int"
}

func (f fn) GoOutputArgs() string {
	var args []string
	if f.Return == "SRTSOCKET" {
		args = append(args, "C.SRTSOCKET")
	} else if f.Return == "int" {
		args = append(args, "C.int")
	}
	if f.RejectReasonArgName != "" {
		args = append(args, "C.int")
	}
	args = append(args, "error")
	o := strings.Join(args, ", ")
	if len(args) > 1 {
		return "(" + o + ")"
	}
	return o
}

func (f fn) GoReturnArgsError() string {
	var args []string
	if f.Return == "SRTSOCKET" || f.Return == "int" {
		args = append(args, "ret")
	}
	if f.RejectReasonArgName != "" {
		args = append(args, "rejectReason")
	}
	args = append(args, "newError(srtErrno, sysErrno)")
	return strings.Join(args, ", ")
}

func (f fn) GoReturnArgsSuccess() string {
	var args []string
	if f.Return == "SRTSOCKET" || f.Return == "int" {
		args = append(args, "ret")
	}
	if f.RejectReasonArgName != "" {
		args = append(args, "rejectReason")
	}
	args = append(args, "nil")
	return strings.Join(args, ", ")
}

var fns = []fn{
	{
		CArgs: []fnArg{
			{
				Name: "u",
				Type: "SRTSOCKET",
			},
			{
				Name: "opt",
				Type: "SRT_SOCKOPT",
			},
			{
				Name: "optval",
				Type: "void*",
			},
			{
				Name: "optlen",
				Type: "int*",
			},
		},
		CName: "getsockflag",
		GoArgs: []fnArg{
			{
				Name: "s",
				Type: "C.SRTSOCKET",
			},
			{
				Name: "o",
				Type: "C.SRT_SOCKOPT",
			},
			{
				Name: "p",
				Type: "unsafe.Pointer",
			},
			{
				Format: "(*C.int)(unsafe.Pointer(%s))",
				Name:   "size",
				Type:   "*int",
			},
		},
		GoName: "GetSockFlag",
	},
	{
		CArgs: []fnArg{
			{
				Name: "u",
				Type: "SRTSOCKET",
			},
			{
				Name: "opt",
				Type: "SRT_SOCKOPT",
			},
			{
				Name: "optval",
				Type: "const void*",
			},
			{
				Name: "optlen",
				Type: "int",
			},
		},
		CName: "setsockflag",
		GoArgs: []fnArg{
			{
				Name: "s",
				Type: "C.SRTSOCKET",
			},
			{
				Name: "o",
				Type: "C.SRT_SOCKOPT",
			},
			{
				Name: "p",
				Type: "unsafe.Pointer",
			},
			{
				Format: "(C.int)(%s)",
				Name:   "size",
				Type:   "int",
			},
		},
		GoName: "SetSockFlag",
	},
	{
		CName:  "startup",
		GoName: "Startup",
	},
	{
		CName:  "cleanup",
		GoName: "Cleanup",
	},
	{
		CName:  "create_socket",
		GoName: "CreateSocket",
		Return: "SRTSOCKET",
	},
	{
		CArgs: []fnArg{
			{
				Name: "u",
				Type: "SRTSOCKET",
			},
		},
		CName: "close",
		GoArgs: []fnArg{
			{
				Name: "s",
				Type: "C.SRTSOCKET",
			},
		},
		GoName: "Close",
	},
	{
		CArgs: []fnArg{
			{
				Name: "u",
				Type: "SRTSOCKET",
			},
			{
				Name: "name",
				Type: "const struct sockaddr*",
			},
			{
				Name: "namelen",
				Type: "int",
			},
		},
		CName: "bind",
		GoArgs: []fnArg{
			{
				Name: "s",
				Type: "C.SRTSOCKET",
			},
			{
				Name: "addr",
				Type: "*C.struct_sockaddr",
			},
			{
				Name: "size",
				Type: "C.int",
			},
		},
		GoName: "Bind",
	},
	{
		CArgs: []fnArg{
			{
				Name: "u",
				Type: "SRTSOCKET",
			},
			{
				Name: "backlog",
				Type: "int",
			},
		},
		CName: "listen",
		GoArgs: []fnArg{
			{
				Name: "s",
				Type: "C.SRTSOCKET",
			},
			{
				Name: "backlog",
				Type: "C.int",
			},
		},
		GoName: "Listen",
	},
	{
		CArgs: []fnArg{
			{
				Name: "u",
				Type: "SRTSOCKET",
			},
			{
				Name: "name",
				Type: "const struct sockaddr*",
			},
			{
				Name: "namelen",
				Type: "int",
			},
		},
		CName: "connect",
		GoArgs: []fnArg{
			{
				Name: "s",
				Type: "C.SRTSOCKET",
			},
			{
				Name: "addr",
				Type: "*C.struct_sockaddr",
			},
			{
				Name: "size",
				Type: "C.int",
			},
		},
		GoName:              "Connect",
		RejectReasonArgName: "u",
	},
	{
		CArgs: []fnArg{
			{
				Name: "u",
				Type: "SRTSOCKET",
			},
			{
				Name: "addr",
				Type: "struct sockaddr*",
			},
			{
				Name: "addrlen",
				Type: "int*",
			},
		},
		CName: "accept",
		GoArgs: []fnArg{
			{
				Name: "s",
				Type: "C.SRTSOCKET",
			},
			{
				Name: "addr",
				Type: "*C.struct_sockaddr",
			},
			{
				Name: "size",
				Type: "*C.int",
			},
		},
		GoName: "Accept",
		Return: "SRTSOCKET",
	},
	{
		CArgs: []fnArg{
			{
				Name: "lsn",
				Type: "SRTSOCKET",
			},
			{
				Name: "astisrt_listen_callback_fn",
			},
			{
				Name: "opaque",
				Type: "void*",
			},
		},
		CName: "listen_callback",
		GoArgs: []fnArg{
			{
				Format: "*%s",
				Name:   "lsn",
				Type:   "*C.SRTSOCKET",
			},
			{
				Name: "unsafe.Pointer(lsn)",
			},
		},
		GoName: "ListenCallback",
	},
	{
		CArgs: []fnArg{
			{
				Name: "clr",
				Type: "SRTSOCKET",
			},
			{
				Name: "astisrt_connect_callback_fn",
			},
			{
				Name: "opaque",
				Type: "void*",
			},
		},
		CName: "connect_callback",
		GoArgs: []fnArg{
			{
				Format: "*%s",
				Name:   "clr",
				Type:   "*C.SRTSOCKET",
			},
			{
				Name: "unsafe.Pointer(clr)",
			},
		},
		GoName: "ConnectCallback",
	},
	{
		CArgs: []fnArg{
			{
				Name: "sock",
				Type: "SRTSOCKET",
			},
			{
				Name: "value",
				Type: "int",
			},
		},
		CName: "setrejectreason",
		GoArgs: []fnArg{
			{
				Name: "sock",
				Type: "C.SRTSOCKET",
			},
			{
				Name: "value",
				Type: "C.int",
			},
		},
		GoName: "SetRejectReason",
	},
	{
		CArgs: []fnArg{
			{
				Name: "u",
				Type: "SRTSOCKET",
			},
			{
				Name: "perf",
				Type: "SRT_TRACEBSTATS*",
			},
			{
				Name: "clear",
				Type: "int",
			},
			{
				Name: "instantaneous",
				Type: "int",
			},
		},
		CName: "bistats",
		GoArgs: []fnArg{
			{
				Name: "u",
				Type: "C.SRTSOCKET",
			},
			{
				Name: "perf",
				Type: "*C.SRT_TRACEBSTATS",
			},
			{
				Name: "clear",
				Type: "C.int",
			},
			{
				Name: "instantaneous",
				Type: "C.int",
			},
		},
		GoName: "BiStats",
	},
	{
		CArgs: []fnArg{
			{
				Name: "u",
				Type: "SRTSOCKET",
			},
			{
				Name: "buf",
				Type: "char*",
			},
			{
				Name: "len",
				Type: "int",
			},
			{
				Name: "mctrl",
				Type: "SRT_MSGCTRL*",
			},
		},
		CName: "recvmsg2",
		GoArgs: []fnArg{
			{
				Name: "u",
				Type: "C.SRTSOCKET",
			},
			{
				Name: "buf",
				Type: "*C.char",
			},
			{
				Name: "len",
				Type: "C.int",
			},
			{
				Name: "mctrl",
				Type: "*C.SRT_MSGCTRL",
			},
		},
		GoName: "RecMsg2",
		Return: "int",
	},
	{
		CArgs: []fnArg{
			{
				Name: "u",
				Type: "SRTSOCKET",
			},
			{
				Name: "buf",
				Type: "char*",
			},
			{
				Name: "len",
				Type: "int",
			},
			{
				Name: "mctrl",
				Type: "SRT_MSGCTRL*",
			},
		},
		CName: "sendmsg2",
		GoArgs: []fnArg{
			{
				Name: "u",
				Type: "C.SRTSOCKET",
			},
			{
				Name: "buf",
				Type: "*C.char",
			},
			{
				Name: "len",
				Type: "C.int",
			},
			{
				Name: "mctrl",
				Type: "*C.SRT_MSGCTRL",
			},
		},
		GoName: "SendMsg2",
		Return: "int",
	},
}

type tpl struct {
	content  string
	filename string
}

var tpls = []tpl{
	{
		content: `// Code generated by astisrt using internal/cmd/generate/wrap. DO NOT EDIT.
#include <srt/srt.h>
#include "callbacks.h"
{{ range $fn := . }}
{{ $fn.CReturn }} astisrt_{{ $fn.CName }}({{ $fn.JoinCArgsWithType }}int *srtErrno, int *sysErrno{{ if $fn.RejectReasonArgName }}, int *rejectReason{{ end }}) {
	{{ $fn.CReturn }} ret = srt_{{ $fn.CName }}({{ $fn.JoinCArgs }});
	if (ret == {{ $fn.CRetCompare }}) {
		*srtErrno = srt_getlasterror(sysErrno);
		srt_clearlasterror();{{ if $fn.RejectReasonArgName }}
		if (*srtErrno == SRT_ECONNREJ) {
			*rejectReason = srt_getrejectreason({{ $fn.RejectReasonArgName }});
		}{{ end }}
	}
	return ret;
}
{{ end }}`,
		filename: "wrap.c",
	},
	{
		content: `// Code generated by astisrt using internal/cmd/generate/wrap. DO NOT EDIT.
#include <srt/srt.h>
#include "callbacks.h"
{{ range $fn := . }}
{{ $fn.CReturn }} astisrt_{{ $fn.CName }}({{ $fn.JoinCArgsWithType }}int *srtErrno, int *sysErrno{{ if $fn.RejectReasonArgName }}, int *rejectReason{{ end }});{{ end }}`,
		filename: "wrap.h",
	},
	{
		content: `// Code generated by astisrt using internal/cmd/generate/wrap. DO NOT EDIT.
package astisrt

// #cgo LDFLAGS: -lsrt
// #include "static_consts.h"
// #include <srt/srt.h>
// #include "wrap.h"
import "C"
import (
	"unsafe"
)
{{ range $fn := . }}
func c{{ $fn.GoName }}({{ $fn.JoinGoArgsWithType }}) {{ $fn.GoOutputArgs }} {
	var srtErrno C.int
	var sysErrno C.int
	{{ if $fn.RejectReasonArgName }}var rejectReason C.int
	{{ end }}ret := C.astisrt_{{ $fn.CName }}({{ $fn.JoinGoArgs }}&srtErrno, &sysErrno{{ if $fn.RejectReasonArgName }}, &rejectReason{{ end }})
	if ret == C.{{ $fn.GoRetCompare }} {
		return {{ $fn.GoReturnArgsError }}
	}
	return {{ $fn.GoReturnArgsSuccess }}
}
{{ end }}`,
		filename: "wrap.go",
	},
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(fmt.Errorf("main: getting working directory failed: %w", err))
	}

	for _, tpl := range tpls {
		log.Printf("generating %s\n", tpl.filename)

		f, err := os.Create(filepath.Join(dir, "pkg", tpl.filename))
		if err != nil {
			log.Fatal(fmt.Errorf("main: creating file failed: %w", err))
		}
		defer f.Close()

		if err = template.Must(template.New("tmpl").Parse(tpl.content)).Execute(f, fns); err != nil {
			log.Fatal(fmt.Errorf("main: executing template failed: %w", err))
		}
	}
}
