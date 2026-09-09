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

func main() {
	startURL, _ := url.Parse("https://crawlme.monzo.com/")
	fetchTimeout := 20 * time.Second
	hf := fetcher.NewHttpFetcher(fetchTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	logger := log.New(os.Stdout, "[CRAWLER] ", log.LstdFlags)

	processor := urlprocessor.NewUrlProcessor("http", "https")
	crawler := crawler.NewCrawler(hf,
		parser.NewHtmlParser(),
		processor,
		crawler.WithLogger(logger),
		crawler.WithWorkercount(1000))

	foundUrls, err := crawler.Crawl(ctx, startURL)

	if err != nil {
		log.Fatalf("Fetch failed: %v", err)
	}

	fmt.Printf("Found total of %d links\n", len(foundUrls))
}
