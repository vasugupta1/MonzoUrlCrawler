package crawler

import (
	"context"
	"log"
	"net/url"

	"github.com/vasugupta1/MonzoUrlCrawler/internal/concurrency"
	"github.com/vasugupta1/MonzoUrlCrawler/internal/fetcher"
)

type Option func(*Crawler)

type Crawler struct {
	fetcher     fetcher.Fetcher
	rateLimiter *concurrency.RateLimiter[[]*url.URL]
	logger      *log.Logger
}

func NewCrawler(fetcher fetcher.Fetcher, opts ...Option) *Crawler {
	c := &Crawler{
		fetcher: fetcher,
		//set default in case caller of NewCrawler forgets to add it
		rateLimiter: concurrency.NewRateLimiter[[]*url.URL](100),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithRateLimit(limit int) Option {
	return func(c *Crawler) {
		c.rateLimiter = concurrency.NewRateLimiter[[]*url.URL](limit)
	}
}

func WithLogger(logger *log.Logger) Option {
	return func(c *Crawler) {
		c.logger = logger
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

	if err != nil {
		c.logger.Printf("Failed to card %s", targetUrl)
		return
	}

	select {
	case result <- urls:
	case <-ctx.Done():
		return
	}
}
