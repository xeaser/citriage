package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/xeaser/citriage/config"
	"github.com/xeaser/citriage/internal/server"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("Failed to load configuration, %v", err)
	}

	ctx := context.Background()
	// TODO: Interface for logger with config
	logger.InfoContext(ctx, "Starting server...")

	s := server.New(&cfg, logger)
	s.ListenAndServe(ctx)

	// TODO: stop signals
	select {}
}
