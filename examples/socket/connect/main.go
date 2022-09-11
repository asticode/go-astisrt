package main

import (
	"fmt"
	"log"
	"net"
	"time"

	astisrt "github.com/asticode/go-astisrt/pkg"
)

func main() {
	// Handle logs
	astisrt.SetLogLevel(astisrt.LogLevelError)
	astisrt.SetLogHandler(func(ll astisrt.LogLevel, file, area, msg string, line int) { log.Println(msg) })

	// Startup srt
	if err := astisrt.Startup(); err != nil {
		log.Fatal(fmt.Errorf("main: starting up srt failed: %w", err))
	}

	// Make sure to clean up srt
	defer func() {
		if err := astisrt.CleanUp(); err != nil {
			log.Fatal(fmt.Errorf("main: cleaning up srt failed: %w", err))
		}
	}()

	// Create socket
	s, err := astisrt.NewSocket()
	if err != nil {
		log.Fatal(fmt.Errorf("main: creating socket failed: %w", err))
	}
	defer s.Close()

	// Set connect callback
	doneConnect := make(chan error)
	if err = s.SetConnectCallback(func(s *astisrt.Socket, addr *net.UDPAddr, token int, err error) {
		doneConnect <- err
	}); err != nil {
		log.Fatal(fmt.Errorf("main: setting connect callback failed: %w", err))
	}

	// Set passphrase
	if err = s.Options().SetPassphrase("passphrase"); err != nil {
		log.Fatal(fmt.Errorf("main: setting passphrase failed: %w", err))
	}

	// Set stream id
	if err = s.Options().SetStreamid("test"); err != nil {
		log.Fatal(fmt.Errorf("main: setting passphrase failed: %w", err))
	}

	// Connect
	log.Println("main: connecting")
	if _, err = s.Connect("127.0.0.1", 4000); err != nil {
		log.Fatal(fmt.Errorf("main: connecting failed: %w", err))
	}

	// Send message
	msg := "this is a test message"
	log.Printf("main: sending `%s`\n", msg)
	if _, err = s.SendMessage([]byte(msg)); err != nil {
		log.Fatal(fmt.Errorf("main: sending message failed: %w", err))
	}

	// Give time to the message to be received
	time.Sleep(500 * time.Millisecond)

	// Close socket
	s.Close()

	// Wait for disconnect
	<-doneConnect
}
