package buildsapi

import (
	"log/slog"

	"github.com/gorilla/mux"
	"github.com/xeaser/citriage/config"
)

func RegisterHandlers(router *mux.Router, cfg *config.Config, logger *slog.Logger) {
	api := api{
		cfg:    cfg,
		logger: logger,
	}

	// router.HandleFunc("/logs/{builds_id}", api.returnContentLength).Methods(http.MethodHead)
	router.HandleFunc("/logs/{builds_id}", api.returnContentLength).Methods("GET")
}
