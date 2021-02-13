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
	"github.com/heat1q/boardsite/redis"
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

	ctx := context.Background()
	run, shutdown := app.Serve(ctx, port)
	defer shutdown()

	if err := redis.InitPool(); err != nil {
		log.Fatal(err)
	}
	defer redis.ClosePool()
	log.Println("Redis connection pool initialized.")

	gracefulStopHook(shutdown)

	if err := run(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func gracefulStopHook(shutdown func() error) {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)
	signal.Notify(stop, syscall.SIGKILL)

	go func() {
		select {
		case <-stop:
			log.Println("Shutting down...")
			redis.ClosePool()
			if err := shutdown(); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}
	}()
}
