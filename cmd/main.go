package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/vasugupta1/MonzoUrlCrawler/internal/crawler"
	"github.com/vasugupta1/MonzoUrlCrawler/internal/fetcher"
	"github.com/vasugupta1/MonzoUrlCrawler/internal/parser"
	"github.com/vasugupta1/MonzoUrlCrawler/internal/urlprocessor"
)

//1. Total links found is 42011
//3: Improve Concurrency
//4: Add logging in crawler, fetcher and html parser

func main() {
	startURL, _ := url.Parse("https://crawlme.monzo.com/")
	htmlParser := parser.NewHtmlParser()
	fetchTimeout := 20 * time.Second
	hf := fetcher.NewHttpFetcher(fetchTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	logger := log.New(os.Stdout, "[CRAWLER] ", log.LstdFlags)

	processor := urlprocessor.NewUrlProcessor(urlprocessor.WithLogger(logger),
		urlprocessor.WithSupportedScheme("http"),
		urlprocessor.WithSupportedScheme("https"))

	crawler := crawler.NewCrawler(hf,
		htmlParser,
		processor,
		crawler.WithLogger(logger))

	foundUrls, err := crawler.Crawl(ctx, startURL)

	if err != nil {
		log.Fatalf("Fetch failed: %v", err)
	}

	fmt.Printf("Found total of %d links\n", len(foundUrls))
}
