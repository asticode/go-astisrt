# SRT server, client and socket in GO

First off, a big thanks to the [official GO bindings](https://github.com/Haivision/srtgo) that was an amazing source of inspiration for this project.

However a few things bothered me in it therefore I decided to write my own bindings with a few goals in mind:

- [x] split the API between both high level entities, similar to GO's `net/http` entities, with very guided methods (`ListenAndServe`, `Shutdown`, `Dial`, `Read`, `Write` etc.) and low level entities, closer to C, with I-should-know-what-I-am-doing-before-using-them methods (`Socket`, etc.)
- [x] provide named and typed option setters and getters
- [x] make sure all errors are handled properly since they are thread-stored and ban the use of runtime.LockOSThread()
- [x] make sure there's a context specific to each connection in high level methods
- [x] make sure pointers are the same between the `ListenCallback` and `Accept()`, and between the `ConnectCallback` and `Connect()` 
- [x] only use blocking mode in high level entities

`astisrt` has been tested on `v1.5.0`.

## Examples

Examples are located in the [examples](examples) directory

WARNING: the code below doesn't handle errors for readibility purposes. However you SHOULD!

### Server

[Go to full example](examples/server/main.go)

```go
// Capture SIGTERM
doneSignal := make(chan os.Signal, 1)
signal.Notify(doneSignal, os.Interrupt)

// Create server
s, _ := astisrt.NewServer(astisrt.ServerOptions{
    // Provide options that will be passed to accepted connections
    ConnectionOptions: []astisrt.ConnectionOption{
        astisrt.WithLatency(300),
        astisrt.WithTranstype(astisrt.TranstypeLive),
    },

    // Specify how an incoming connection should be handled before being accepted
    // When false is returned, the connection is rejected.
    OnBeforeAccept: func(c *astisrt.Connection, version int, streamID string) bool {
        // Check stream id
        if streamID != "test" {
            // Set reject reason
            c.SetPredefinedRejectReason(http.StatusNotFound)
            return false
        }

        // Update passphrase
        c.Options().SetPassphrase("passphrase")

        // Add stream id to context
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
            n, _ := c.Read(b)

            // Log
            log.Printf("main: read `%s`\n", b[:n])

            // Get stats
            s, _ := c.Stats(false, false)

            // Log
            log.Printf("main: %d total bytes received\n", s.ByteRecvTotal())
        }
    }),

    // Addr the server should be listening to
    Host: "127.0.0.1",
    Port: 4000,
})
defer s.Close()

// Listen and serve in a goroutine
doneListenAndServe := make(chan error)
go func() { doneListenAndServe <- s.ListenAndServe(1) }()

// Wait for SIGTERM
<-doneSignal

// Create shutdown context with a timeout to make sure it's cancelled if it takes too much time
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Shutdown
s.Shutdown(ctx)

// Wait for listen and serve to be done
<-doneListenAndServe
```

### Client

[Go to full example](examples/client/main.go)

```go
// Capture SIGTERM
doneSignal := make(chan os.Signal, 1)
signal.Notify(doneSignal, os.Interrupt)

// Dial
doneWrite := make(chan err)
c, _ := astisrt.Dial(astisrt.DialOptions{
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
defer c.Close()

// Write in a goroutine
go func() {
    defer func() { close(doneWrite) }()

    // Loop
    r := bufio.NewReader(os.Stdin)
    for {
        // Read from stdin
        t, _ := r.ReadString('\n')

        // Write to the server
        c.Write([]byte(t))
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
```

### Socket

#### Listen

[Go to full example](examples/socket/listen/main.go)

```go
// Create socket
s, _ := astisrt.NewSocket()
defer s.Close()

// Set listen callback
s.SetListenCallback(func(s *astisrt.Socket, version int, addr *net.UDPAddr, streamID string) bool {
    // Check stream id
    if streamID != "test" {
        // Set reject reason
        s.SetRejectReason(1404)
        return false
    }

    // Update passphrase
    s.Options().SetPassphrase("passphrase")
    return true
})

// Bind
s.Bind("127.0.0.1", 4000)

// Listen
s.Listen(1)

// Accept
as, _, _ := s.Accept()

// Receive message
b := make([]byte, 1500)
n, _ := as.ReceiveMessage(b)

// Log
log.Printf("main: received `%s`\n", b[:n])
```

#### Connect

[Go to full example](examples/socket/listen/main.go)

```go
// Create socket
s, _ := astisrt.NewSocket()
defer s.Close()

// Set connect callback
doneConnect := make(chan error)
s.SetConnectCallback(func(s *astisrt.Socket, addr *net.UDPAddr, token int, err error) {
    doneConnect <- err
}

// Set passphrase
s.Options().SetPassphrase("passphrase")

// Set stream id
s.Options().SetStreamid("test")

// Connect
s.Connect("127.0.0.1", 4000)

// Send message
s.SendMessage([]byte("this is a test message"))

// Give time to the message to be received
time.Sleep(500 * time.Millisecond)

// Close socket
s.Close()

// Wait for disconnect
<-doneConnect
```

## Install `srtlib` from source

You can find the instructions to install `srtlib` [here](https://github.com/Haivision/srt/tree/master/docs/build).

However if you don't feel like doing it manually you can use the following command:

```sh
$ make install-srt
```

`srtlib` will be built from source in a directory named `tmp` and located in you working directory.

For your GO code to pick up `srtlib` dependency automatically, you'll need to add the following environment variables:

(don't forget to replace `{{ path to your working directory }}` with the absolute path to your working directory)

```sh
export CGO_LDFLAGS="-L{{ path to your working directory }}/tmp/v1.5.0/lib/",
export CGO_CXXFLAGS="-I{{ path to your working directory }}/tmp/v1.5.0/include/",
export PKG_CONFIG_PATH="{{ path to your working directory }}/tmp/v1.5.0/lib/pkgconfig",
```