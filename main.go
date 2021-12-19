package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/heat1q/boardsite/api"
)

func main() {
	ctx := context.Background()
	runServer(ctx)
}

func runServer(ctx context.Context) {
	s, err := api.NewServer()
	if err != nil {
		log.Fatal(err.Error())
	}
	run, shutdown := s.Serve(ctx)
	defer shutdown()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	go func() {
		select {
		case <-stop:
			log.Println("Shutting down...")
			if err := shutdown(); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}
	}()

	if err := run(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
