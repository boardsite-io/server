package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/heat1q/boardsite/pkg/app"
)

func main() {
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	ctx := context.Background()
	run, shutdown := app.Serve(ctx, 8000)
	defer shutdown()

	go func() {
		select {
		case <-gracefulStop:
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