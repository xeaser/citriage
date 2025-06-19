package cache

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/xeaser/citriage/pkg/httpclient"
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remote server returned non-200 status: %s", resp.Status)

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
