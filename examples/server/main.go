package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	astisrt "github.com/asticode/go-astisrt/pkg"
)

type ctxKey string

const (
	ctxKeyStreamID ctxKey = "stream.id"
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

	// Create server
	s, err := astisrt.NewServer(astisrt.ServerOptions{
		// Provide options that will be passed to accepted connections
		ConnectionOptions: []astisrt.ConnectionOption{
			astisrt.WithLatency(300),
			astisrt.WithTranstype(astisrt.TranstypeLive),
		},

		// Specify how an incoming connection should be handled before being accepted
		// When false is returned, the connection is rejected.
		OnBeforeAccept: func(c *astisrt.Connection, version int, streamID string) bool {
			// Check stream id
			log.Printf("main: checking stream id %s\n", streamID)
			if streamID != "test" {
				// Set reject reason
				log.Println("main: invalid stream id")
				if err := c.SetPredefinedRejectReason(http.StatusNotFound); err != nil {
					log.Println(fmt.Errorf("main: setting predefined reject reason failed: %w", err))
				}
				return false
			}

			// Update passphrase
			if err := c.Options().SetPassphrase("passphrase"); err != nil {
				log.Println(fmt.Errorf("main: setting passphrase failed: %w", err))
				return false
			}

			// Add stream id to context
			log.Println("main: connection accepted")
			*c = *c.WithContext(context.WithValue(c.Context(), ctxKeyStreamID, streamID))
			return true
		},

		// Similar to http.Handler behavior, specify how a connection
		// will be handled once accepted
		Handler: astisrt.ServerHandlerFunc(func(c *astisrt.Connection) {
			// Get stream id from context
			if v := c.Context().Value(ctxKeyStreamID); v != nil {
				log.Printf("main: handling connection with stream id %s\n", v.(string))
			}

			// Loop
			for {
				// Read
				b := make([]byte, 1500)
				n, err := c.Read(b)
				if err != nil {
					log.Println(fmt.Errorf("main: reading failed: %w", err))
					return
				}

				// Log
				log.Printf("main: read `%s`\n", b[:n])

				// Get stats
				s, err := c.Stats(false, false)
				if err != nil {
					log.Println(fmt.Errorf("main: getting stats failed: %w", err))
					continue
				}

				// Log
				log.Printf("main: %d total bytes received\n", s.ByteRecvTotal())
			}
		}),

		// Addr the server should be listening to
		Host: "127.0.0.1",
		Port: 4000,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("main: creating server failed: %w", err))
	}
	defer s.Close()

	// Listen and serve in a goroutine
	doneListenAndServe := make(chan error)
	go func() {
		log.Println("main: listening")
		doneListenAndServe <- s.ListenAndServe(1)
	}()

	// Wait for SIGTERM
	<-doneSignal

	// Create shutdown context with a timeout to make sure it's cancelled if it takes too much time
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown
	log.Println("main: shutting down")
	if err = s.Shutdown(ctx); err != nil {
		log.Println(fmt.Errorf("main: shutting down failed: %w", err))
	}

	// Wait for listen and serve to be done
	<-doneListenAndServe
}
