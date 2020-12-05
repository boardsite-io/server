package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/heat1q/boardsite/app"
)

func main() {
	portStr := flag.String("port", os.Getenv("B_PORT"), "listening on port")
	redisHost := flag.String("redis-host", os.Getenv("B_REDIS_HOST"), "redis hostname")
	redisPort := flag.String("redis-port", os.Getenv("B_REDIS_PORT"), "redis port")
	flag.Parse()

	// input args may overwrite os env
	os.Setenv("B_PORT", *portStr)
	os.Setenv("B_REDIS_HOST", *redisHost)
	os.Setenv("B_REDIS_PORT", *redisPort)

	port, err := strconv.Atoi(*portStr)
	if err != nil {
		log.Fatal(err)
	}

	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	ctx := context.Background()
	run, shutdown := app.Serve(ctx, port)
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
