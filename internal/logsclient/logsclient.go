package logsclient

import (
	"fmt"
	"io"
	"net/http"

	"github.com/xeaser/citriage/config"
	"github.com/xeaser/citriage/pkg/httpclient"
)

type LogsClient struct {
	baseURL   string
	apiClient httpclient.ApiClient
}

func NewLogsAPIClient(cfg *config.Config, httpClient *http.Client) *LogsClient {
	return &LogsClient{
		baseURL:   fmt.Sprintf("http://%s:%d", cfg.Server.Host, cfg.Server.Port),
		apiClient: httpClient,
	}
}

func (l *LogsClient) FetchLog(buildID string) ([]byte, error) {
	requestURL := fmt.Sprintf("%s/logs/%s", l.baseURL, buildID)

	resp, err := l.apiClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("could not send request to proxy server: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("proxy server returned an error (%s): %s", resp.Status, string(body))
	}

	return body, nil
}
