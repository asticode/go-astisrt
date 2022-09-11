package main

import (
	"fmt"
	"log"
	"net"

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

	// Set listen callback
	if err = s.SetListenCallback(func(s *astisrt.Socket, version int, addr *net.UDPAddr, streamID string) bool {
		// Check stream id
		log.Printf("main: checking stream id %s\n", streamID)
		if streamID != "test" {
			// Set reject reason
			log.Println("main: invalid stream id")
			if err := s.SetRejectReason(1404); err != nil {
				log.Println(fmt.Errorf("main: setting reject reason failed: %w", err))
			}
			return false
		}

		// Update passphrase
		if err := s.Options().SetPassphrase("passphrase"); err != nil {
			log.Println(fmt.Errorf("main: setting passphrase failed: %w", err))
			return false
		}

		// Log
		log.Println("main: connection accepted")
		return true
	}); err != nil {
		log.Fatal(fmt.Errorf("main: setting listen callback failed: %w", err))
	}

	// Bind
	if err = s.Bind("127.0.0.1", 4000); err != nil {
		log.Fatal(fmt.Errorf("main: binding failed: %w", err))
	}

	// Listen
	log.Println("main: listening")
	if err = s.Listen(1); err != nil {
		log.Fatal(fmt.Errorf("main: listening failed: %w", err))
	}

	// Accept
	as, _, err := s.Accept()
	if err != nil {
		log.Fatal(fmt.Errorf("main: accepting failed: %w", err))
	}

	// Receive message
	b := make([]byte, 1500)
	n, err := as.ReceiveMessage(b)
	if err != nil {
		log.Fatal(fmt.Errorf("main: receiving message failed: %w", err))
		return
	}

	// Log
	log.Printf("main: received `%s`\n", b[:n])
}
