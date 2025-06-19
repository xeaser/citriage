package logsapi

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xeaser/citriage/config"
	"github.com/xeaser/citriage/internal/cache"
)

func RegisterHandlers(router *mux.Router, cfg *config.Config, logger *slog.Logger, downloader *cache.Downloader) {
	api := api{
		cfg:        cfg,
		logger:     logger,
		downloader: downloader,
	}

	router.HandleFunc("/logs/{build_id}", api.getBuildLogs).Methods(http.MethodGet, http.MethodHead)
}
