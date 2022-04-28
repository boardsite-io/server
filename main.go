package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/boardsite-io/server/internal/config"
	"github.com/boardsite-io/server/internal/server"
	"github.com/boardsite-io/server/pkg/log"
)

func main() {
	ctx := context.Background()
	cfgPath := flag.String("config", "./config.yaml", "path to the config file")
	flag.Parse()
	runServer(ctx, *cfgPath)
}

func runServer(ctx context.Context, cfgPath string) {
	cfg, err := config.New(cfgPath)
	if err != nil {
		log.Global().Fatalf("parse config file: %v", err)
	}

	s := server.New(cfg)
	run, shutdown := s.Serve(ctx)
	defer shutdown()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	go func() {
		select {
		case <-stop:
			log.Global().Warn("Shutting down...")
			if err := shutdown(); err != nil {
				log.Global().Fatal(err)
			}
			os.Exit(0)
		}
	}()

	if err := run(); !errors.Is(err, http.ErrServerClosed) {
		log.Global().Fatal(err)
	}
}
