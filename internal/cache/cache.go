package cache

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/xeaser/citriage/pkg/httpclient"
	"github.com/xeaser/citriage/pkg/proxyerrors"
)

const (
	jenkinsBuildsLogsUrl = "https://ci.jenkins.io/job/Core/job/jenkins/job/master/%d/consoleText"
)

type Downloader struct {
	cacheDir      string
	httpClient    httpclient.ApiClient
	downloadLocks sync.Map
}

func NewDownloader(httpClient httpclient.ApiClient, cacheDir string) (*Downloader, error) {
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}
	return &Downloader{
		cacheDir:      cacheDir,
		httpClient:    httpClient,
		downloadLocks: sync.Map{},
	}, nil
}

// GetLogFile retrieves the log file for the specified Jenkins build ID.
// If the log file is already cached locally, it returns the path to the cached file.
// Otherwise, it downloads the log file from the Jenkins server, caches it locally,
// and returns the path to the cached file. The method ensures that concurrent
// requests for the same build ID are synchronized, so the log file is only
// downloaded once. If another goroutine is already downloading the file,
// this method waits until the download is complete and then returns the cached file path.
//
// Parameters:
//   - buildID: The Jenkins build ID for which to retrieve the log file.
//
// Returns:
//   - string: The path to the cached log file.
//   - error: An error if the log file could not be retrieved or cached.
func (d *Downloader) GetLogFile(buildID int) (string, error) {
	filePath := filepath.Join(d.cacheDir, fmt.Sprintf("%d.log", buildID))

	if _, err := os.Stat(filePath); err == nil {
		return filePath, nil
	}

	_, loaded := d.downloadLocks.LoadOrStore(buildID, true)
	if loaded {
		for {
			if _, err := os.Stat(filePath); err == nil {
				return filePath, nil
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
	defer d.downloadLocks.Delete(buildID)

	if err := d.downloadAndCache(buildID, filePath); err != nil {
		return "", err
	}

	return filePath, nil
}

func (d *Downloader) downloadAndCache(buildID int, destPath string) error {
	url := fmt.Sprintf(jenkinsBuildsLogsUrl, buildID)

	resp, err := d.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return fmt.Errorf("%w: received status %s", proxyerrors.ErrUpstreamClientError, resp.Status)
	}
	if resp.StatusCode >= 500 && resp.StatusCode < 600 {
		return fmt.Errorf("%w: received status %s", proxyerrors.ErrUpstreamServerError, resp.Status)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upstream server returned unexpected status: %s", resp.Status)
	}

	tmpFile, err := os.CreateTemp(d.cacheDir, "download-*.tmp")
	if err != nil {
		return fmt.Errorf("could not create temp file: %w", err)
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		os.Remove(tmpFile.Name())
		return fmt.Errorf("failed to write to temp file: %w", err)
	}

	if err := os.Rename(tmpFile.Name(), destPath); err != nil {
		return fmt.Errorf("failed to move temp file to destination: %w", err)
	}

	return nil
}
