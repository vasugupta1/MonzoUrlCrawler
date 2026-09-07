package fetcher

import (
	"fmt",
	"context"
	"net/http"
	"time"
)

type fetcher interface {
	Fetch(ctx context.Context, url string) ([]string, error)
}

type HttpFetcher struct {
	client *http.Client
	parser *HtmlParser
}

func NewHttpFetcher(timeout time.Duration, htmlParser *HtmlParser) *HttpFetcher {
	return &HttpFetcher{
		client: &http.Client{
			Timeout: timeout,
		},
		htmlParser: htmlParser
	}
}

func (hf *HttpFetcher) Fetch(ctx context.Context, url string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := hf.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errof("Unsucessful status code : %s", resp.StatusCode)
	}


	return hf.parser(resp.Body, url)
}
