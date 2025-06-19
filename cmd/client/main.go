package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/xeaser/citriage/config"
	"github.com/xeaser/citriage/internal/dedupe"
	"github.com/xeaser/citriage/internal/logsclient"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	buildIDFlag := flag.String("build-id", "", "The build ID of the log to fetch")
	flag.Parse()
	buildID := *buildIDFlag
	if buildID == "" {
		logger.Error("Error: --build-id is required")
		os.Exit(1)
	}

	cfg, err := config.Get()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to load configuration, %v", err))
		os.Exit(1)
	}

	httpClient := &http.Client{
		Timeout: time.Duration(cfg.HttpClient.TimeoutSecs) * time.Second,
	}

	logApiClient := logsclient.NewLogsAPIClient(&cfg, httpClient)
	logger.Info(fmt.Sprintf("Fetching log for build '%s'", buildID))

	logData, err := logApiClient.FetchLog(buildID)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to fetch log: %v", err))
		os.Exit(1)
	}
	logger.Info("Log fetched successfully. Processing for deduplication...")

	processedLog := dedupe.Process(logData)

	logger.Info("\n--- Deduplicated Log Output ---")
	fmt.Print(processedLog)
}
