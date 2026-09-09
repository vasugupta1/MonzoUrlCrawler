package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Fetcher interface {
	Fetch(ctx context.Context, url *url.URL) (io.ReadCloser, error)
}

type HttpFetcher struct {
	client *http.Client
}

func NewHttpFetcher(timeout time.Duration) *HttpFetcher {
	return &HttpFetcher{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (hf *HttpFetcher) Fetch(ctx context.Context, url *url.URL) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := hf.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("Unsucessful status code : %d", resp.StatusCode)
	}

	return resp.Body, nil
}
