package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

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

	// Capture SIGTERM
	doneSignal := make(chan os.Signal, 1)
	signal.Notify(doneSignal, os.Interrupt)

	// Dial
	log.Println("main: dialing")
	doneWrite := make(chan error)
	c, err := astisrt.Dial(astisrt.DialOptions{
		// Provide options to the connection
		ConnectionOptions: []astisrt.ConnectionOption{
			astisrt.WithLatency(300),
			astisrt.WithPassphrase("passphrase"),
			astisrt.WithStreamid("test"),
		},

		// Callback when the connection is disconnected
		OnDisconnect: func(c *astisrt.Connection, err error) { doneWrite <- err },

		// Addr that should be dialed
		Host: "127.0.0.1",
		Port: 4000,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("main: dialing failed: %w", err))
	}
	defer c.Close()

	// Write in a goroutine
	go func() {
		defer func() { close(doneWrite) }()

		// Loop
		log.Println("Write some text and press Enter to send it to the server. Write \"exit\" to quit.")
		r := bufio.NewReader(os.Stdin)
		for {
			// Read from stdin
			t, err := r.ReadString('\n')
			if err != nil {
				log.Println(fmt.Errorf("main: reading from stdin failed: %w", err))
				return
			}
			t = strings.TrimSpace(t)

			// Exit
			if t == "exit" {
				return
			}

			// Write to the server
			if _, err = c.Write([]byte(t)); err != nil {
				log.Println(fmt.Errorf("main: writing to server failed: %w", err))
				return
			}
		}
	}()

	// Wait for either SIGTERM or write end
	select {
	case <-doneSignal:
		c.Close()
	case err := <-doneWrite:
		log.Println(err)
	}

	// Make sure write is done
	select {
	case <-doneWrite:
	default:
	}
}
