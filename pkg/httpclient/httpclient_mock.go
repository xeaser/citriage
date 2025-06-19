package httpclient

import (
	"errors"
	"net/http"
)

type MockClient struct {
	GetFunc func(url string) (*http.Response, error)
}

func (m *MockClient) Get(url string) (*http.Response, error) {
	if m.GetFunc != nil {
		return m.GetFunc(url)
	}
	return nil, errors.New("GetFunc not implemented")
}
