package fetcher

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/vasugupta1/MonzoUrlCrawler/internal/parser"
)

type Fetcher interface {
	Fetch(ctx context.Context, url *url.URL) ([]*url.URL, error)
}

type HttpFetcher struct {
	client *http.Client
	parser *parser.HtmlParser
}

func NewHttpFetcher(timeout time.Duration, htmlParser *parser.HtmlParser) *HttpFetcher {
	return &HttpFetcher{
		client: &http.Client{
			Timeout: timeout,
		},
		parser: htmlParser,
	}
}

func (hf *HttpFetcher) Fetch(ctx context.Context, url *url.URL) ([]*url.URL, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := hf.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unsucessful status code : %d", resp.StatusCode)
	}

	return hf.parser.Parse(resp.Body)
}
