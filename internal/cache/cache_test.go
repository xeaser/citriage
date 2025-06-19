package cache

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/xeaser/citriage/pkg/httpclient"
)

func TestDownloader_GetLogFile(t *testing.T) {
	t.Run("downloads file successfully when not in cache", func(t *testing.T) {
		mockClient := &httpclient.MockClient{
			GetFunc: func(url string) (*http.Response, error) {
				body := io.NopCloser(bytes.NewReader([]byte("log content")))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       body,
				}, nil
			},
		}

		tempDir := t.TempDir()
		downloader, err := NewDownloader(mockClient, tempDir)
		if err != nil {
			t.Fatalf("Failed to create downloader: %v", err)
		}

		filePath, err := downloader.GetLogFile(123)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}

		expectedPath := filepath.Join(tempDir, "123.log")
		if filePath != expectedPath {
			t.Errorf("Expected file path %s, but got %s", expectedPath, filePath)
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read created file: %v", err)
		}
		if string(content) != "log content" {
			t.Errorf("Expected file content 'log content', but got '%s'", string(content))
		}
	})

	t.Run("returns existing file when already in cache", func(t *testing.T) {
		getCalled := false
		mockClient := &httpclient.MockClient{
			GetFunc: func(url string) (*http.Response, error) {
				getCalled = true
				return nil, errors.New("http client should not be called")
			},
		}

		tempDir := t.TempDir()
		cachedFilePath := filepath.Join(tempDir, "10.log")
		if err := os.WriteFile(cachedFilePath, []byte("existing content"), 0644); err != nil {
			t.Fatalf("Failed to create pre-cached file: %v", err)
		}

		downloader, err := NewDownloader(mockClient, tempDir)
		if err != nil {
			t.Fatalf("Failed to create downloader: %v", err)
		}

		filePath, err := downloader.GetLogFile(10)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		if getCalled {
			t.Error("Expected HTTP client's Get method NOT to be called, but it was.")
		}
		if filePath != cachedFilePath {
			t.Errorf("Expected path %s, got %s", cachedFilePath, filePath)
		}
	})
}
