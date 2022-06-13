package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.SetFlags(0)

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

// run starts a http.Server for the passed in address
// with all requests handled by echoServer.
func run() error {
	if len(os.Args) != 2 {
		return errors.New("please provide an address:port to listen on as the first argument")
	}
	address := os.Args[1]

	baseCtx := context.Background()
	serverCtx, cancel := context.WithCancel(baseCtx)

	server := NewServer(log.Printf)

	errc := make(chan error, 1)
	go func() {
		errc <- server.Serve(address, serverCtx)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	// cancel the server context to stop all services
	cancel()

	ctx, cancel := context.WithTimeout(baseCtx, time.Second*10)
	defer cancel()

	return server.Shutdown(ctx)
}
