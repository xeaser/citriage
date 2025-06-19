package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/xeaser/citriage/config"
	"github.com/xeaser/citriage/internal/cache"
	"github.com/xeaser/citriage/internal/server"
)

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	logger.InfoContext(ctx, "Iniatilsing server dependecies...")

	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("Failed to load configuration, %v", err)
	}

	apiClient := &http.Client{
		Timeout: time.Duration(cfg.HttpClient.TimeoutSecs) * time.Second,
	}

	downloader, err := cache.NewDownloader(apiClient, cfg.Cache.Dir)
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to create logs downloader, err: %v", err))
		log.Fatal("Server initialization failed")
	}

	s := server.New(&cfg, logger, downloader)
	s.ListenAndServe(ctx)

	// TODO: stop signals
	select {}
}
