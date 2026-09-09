package crawler

import (
	"context"
	"log"
	"net/url"

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
		workerCount: 100, // default value
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

	visisted := make([]string, 0, len(seen))
	for u := range seen {
		visisted = append(visisted, u)
	}

	return visisted, nil
}

// func (c *Crawler) Crawl(ctx context.Context, base *url.URL) ([]string, error) {

// 	body, err := c.fetcher.Fetch(ctx, base)
// 	if err != nil {
// 		//log
// 		return nil, err
// 	}
// 	defer body.Close()

// 	urls, err := c.parser.Parse(body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	processedUrls := make([]*url.URL, 0, len(urls))
// 	for _, u := range urls {
// 		processedUrl, err := c.processor.ProcessUrl(u, base)
// 		if err != nil {
// 			c.logger.Printf("Failed to process url: %s", u)
// 			continue
// 		}
// 		processedUrls = append(processedUrls, processedUrl)

// 	}

// 	results := make(chan []*url.URL)
// 	seen := make(map[string]struct{})

// 	// set base as seen
// 	seen[base.String()] = struct{}{}

// 	activeFetches := 0
// 	for _, url := range processedUrls {
// 		key := url.String()
// 		if _, ok := seen[key]; !ok {
// 			seen[key] = struct{}{}
// 			activeFetches++
// 			go c.crawl(ctx, url, results)
// 		}
// 	}

// 	for activeFetches > 0 {
// 		select {
// 		case <-ctx.Done():
// 			activeFetches = 0
// 		case discoveredLinks := <-results:
// 			activeFetches--
// 			for _, url := range discoveredLinks {
// 				key := url.String()
// 				if _, ok := seen[key]; !ok {
// 					seen[key] = struct{}{}
// 					activeFetches++
// 					go c.crawl(ctx, url, results)
// 				}
// 			}
// 		}

// 	}

// 	var visited []string
// 	for url := range seen {
// 		visited = append(visited, url)
// 	}
// 	return visited, nil
// }

// // / this will need to be updated
// func (c *Crawler) crawl(ctx context.Context, targetUrl *url.URL, result chan<- []*url.URL) {

// 	urls, err := c.rateLimiter.Process(ctx, func() ([]*url.URL, error) {
// 		body, err := c.fetcher.Fetch(ctx, targetUrl)
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer body.Close()
// 		urls, err := c.parser.Parse(body)
// 		if err != nil {
// 			return nil, err
// 		}

// 		return urls, nil
// 	})

// 	if err != nil {
// 		c.logger.Printf("Failed to crawl %s", targetUrl)
// 		return
// 	}

// 	processedUrls := make([]*url.URL, 0, len(urls))
// 	for _, u := range urls {
// 		processedUrl, err := c.processor.ProcessUrl(u, targetUrl)
// 		if err != nil {
// 			c.logger.Printf("Failed to process url: %s", u)
// 			continue
// 		}
// 		processedUrls = append(processedUrls, processedUrl)

// 	}

// 	select {
// 	case result <- processedUrls:
// 	case <-ctx.Done():
// 		return
// 	}
// }
