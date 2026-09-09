package crawler

import (
	"context"
	"log"
	"net/url"

	"github.com/vasugupta1/MonzoUrlCrawler/internal/concurrency"
	"github.com/vasugupta1/MonzoUrlCrawler/internal/fetcher"
	"github.com/vasugupta1/MonzoUrlCrawler/internal/parser"
	"github.com/vasugupta1/MonzoUrlCrawler/internal/urlprocessor"
)

type Option func(*Crawler)

type Crawler struct {
	fetcher     fetcher.Fetcher
	workerCount int
	logger      *log.Logger
	processor   urlprocessor.Processor
	parser      *parser.HtmlParser
}

func WithWorkercount(count int) Option {
	return func(c *Crawler) {
		c.workerCount = count
	}
}

func WithLogger(logger *log.Logger) Option {
	return func(c *Crawler) {
		c.logger = logger
	}
}

func NewCrawler(fetcher fetcher.Fetcher, parser *parser.HtmlParser, processor urlprocessor.Processor, opts ...Option) *Crawler {
	c := &Crawler{
		fetcher:     fetcher,
		parser:      parser,
		workerCount: 1000, // default value
		processor:   processor,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Crawler) Crawl(ctx context.Context, base *url.URL) ([]string, error) {
	//set base as seen so we don't crawl this again incase a page has the index button href link
	seen := make(map[string]struct{})
	seen[base.String()] = struct{}{}

	workerPool := concurrency.NewWorkerPool(c.workerCount)
	workerPool.Submit(base)
	workerPool.Start(ctx, func(ctx context.Context, url *url.URL) ([]*url.URL, error) {
		body, err := c.fetcher.Fetch(ctx, url)
		if err != nil {
			return nil, err
		}
		defer body.Close()
		urls, err := c.parser.Parse(body)
		if err != nil {
			return nil, err
		}
		return urls, nil
	})

	for result := range workerPool.Results() {
		if result.Err != nil {
			workerPool.Done()
			continue
		}
		for _, discoverUrl := range result.DiscoverdUrls {
			processedUrl, err := c.processor.ProcessUrl(discoverUrl, result.SourceUrl)
			if err != nil {
				//log and continue
				continue
			}

			key := processedUrl.String()
			if _, exists := seen[key]; !exists {
				seen[key] = struct{}{}
				workerPool.Submit(processedUrl)
			}
		}
		workerPool.Done()
	}

	visisted := make([]string, 0, len(seen))
	for u := range seen {
		visisted = append(visisted, u)
	}

	return visisted, nil
}
