package crawler

import (
	"context"
	"net/url"

	"github.com/vasugupta1/MonzoUrlCrawler/internal/concurrency"
	"github.com/vasugupta1/MonzoUrlCrawler/internal/fetcher"
)

type Crawler struct {
	fetcher     fetcher.Fetcher
	rateLimiter *concurrency.RateLimiter[[]*url.URL]
}

func NewCrawler(fetcher fetcher.Fetcher, limit int) *Crawler {
	return &Crawler{
		fetcher:     fetcher,
		rateLimiter: concurrency.NewRateLimiter[[]*url.URL](limit),
	}
}

func (c *Crawler) Crawl(ctx context.Context, base *url.URL) ([]string, error) {

	urls, err := c.fetcher.Fetch(ctx, base)
	if err != nil {
		//log
		return nil, err
	}

	results := make(chan []*url.URL)
	seen := make(map[string]struct{})

	// set base as seen
	seen[base.String()] = struct{}{}

	activeFetches := 0
	for _, url := range urls {
		key := url.String()
		if _, ok := seen[key]; !ok {
			seen[key] = struct{}{}
			activeFetches++
			go c.crawl(ctx, url, results)
		}
	}

	for activeFetches > 0 {
		select {
		case <-ctx.Done():
			activeFetches = 0
		case discoveredLinks := <-results:
			activeFetches--
			for _, url := range discoveredLinks {
				key := url.String()
				if _, ok := seen[key]; !ok {
					seen[key] = struct{}{}
					activeFetches++
					go c.crawl(ctx, url, results)
				}
			}
		}

	}

	var visited []string
	for url := range seen {
		visited = append(visited, url)
	}
	return visited, nil
}

func (c *Crawler) crawl(ctx context.Context, targetUrl *url.URL, result chan<- []*url.URL) {

	urls, err := c.rateLimiter.Process(ctx, func() ([]*url.URL, error) {
		return c.fetcher.Fetch(ctx, targetUrl)
	})

	var sendResult []*url.URL
	if err == nil {
		sendResult = urls
	}

	select {
	case result <- sendResult:
	case <-ctx.Done():
		return
	}
}
