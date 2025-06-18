package buildsapi

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xeaser/citriage/config"
	"github.com/xeaser/citriage/internal/response"
)

type api struct {
	logger *slog.Logger
	cfg    *config.Config
}

func (a *api) returnContentLength(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	buildIdVar := vars["build_id"]
	_, err := strconv.Atoi(buildIdVar)
	if err != nil {
		errMsg := "Invalid build id"
		a.logger.ErrorContext(ctx, errMsg)
		response.RespondWithJSON(w, http.StatusBadRequest, errMsg)
	}
}
