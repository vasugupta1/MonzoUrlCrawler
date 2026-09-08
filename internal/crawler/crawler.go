package crawler

import (
	"context"
	"net/url"

	"github.com/vasugupta1/MonzoUrlCrawler/internal/fetcher"
)

type Crawler struct {
	Fetcher fetcher.Fetcher
}

func NewCrawler(fetcher fetcher.Fetcher) *Crawler {
	return &Crawler{
		Fetcher: fetcher,
	}
}

func (c *Crawler) Crawl(ctx context.Context, base *url.URL) ([]string, error) {

	urls, err := c.Fetcher.Fetch(ctx, base)
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
			break
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

func (c *Crawler) crawl(ctx context.Context, url *url.URL, result chan<- []*url.URL) {
	select {
	case <-ctx.Done():
		return
	default:
		urls, err := c.Fetcher.Fetch(ctx, url)
		if err != nil {
			result <- nil
			return
		}
		result <- urls
	}
}
