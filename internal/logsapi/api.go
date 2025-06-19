package logsapi

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xeaser/citriage/config"
	"github.com/xeaser/citriage/internal/cache"
)

const (
	urlBuildIdVar = "build_id"

	urlOffsetKey = "offset"
	urlLimitKey  = "limit"
)

type api struct {
	logger     *slog.Logger
	cfg        *config.Config
	downloader *cache.Downloader
}

func (a api) getBuildLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	ctx := r.Context()
	vars := mux.Vars(r)
	buildID, err := strconv.Atoi(vars[urlBuildIdVar])
	if err != nil {
		errMsg := "Invalid build id"
		a.logger.ErrorContext(ctx, errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
	}

	a.logger.InfoContext(ctx, fmt.Sprintf("Received request for build: %d", buildID))

	offsetParam := r.URL.Query().Get(urlOffsetKey)
	limitParam := r.URL.Query().Get(urlLimitKey)

	filePath, err := a.downloader.GetLogFile(buildID)
	if err != nil {
		a.logger.ErrorContext(ctx, fmt.Sprintf("Error getting log file for build %d: %v", buildID, err))
		http.Error(w, fmt.Sprintf("Failed to retrieve log file due to error: %s", err.Error()), http.StatusNoContent)
		return
	}

	logData, err := os.ReadFile(filePath)
	if err != nil {
		a.logger.ErrorContext(ctx, fmt.Sprintf("Error opening cached file for build %d: %v", buildID, err))
		http.Error(w, "Could not open log file", http.StatusInternalServerError)
		return
	}
	fileSize := len(logData)

	w.Header().Set("Content-Length", strconv.Itoa(fileSize))
	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}

	dataSlice := handleFilePagination(offsetParam, limitParam, fileSize, logData)

	w.Header().Set("Content-Length", strconv.Itoa(len(dataSlice)))
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(dataSlice)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Error writing response, buildId: %d, err: %v", buildID, err))
	}

	a.logger.InfoContext(ctx, fmt.Sprintf("Request processed for build: %d", buildID))
}

func handleFilePagination(offsetParam string, limitParam string, fileSize int, logData []byte) []byte {
	offset, _ := strconv.Atoi(offsetParam)
	if offset < 0 || int(offset) > fileSize {
		offset = 0
	}

	limit, _ := strconv.Atoi(limitParam)
	if limit <= 0 {
		limit = fileSize
	}

	start := offset
	end := min(start+limit, fileSize)
	if start > end {
		start = end
	}

	dataSlice := logData[start:end]
	return dataSlice
}
