package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/vasugupta1/MonzoUrlCrawler/internal/crawler"
	"github.com/vasugupta1/MonzoUrlCrawler/internal/fetcher"
	"github.com/vasugupta1/MonzoUrlCrawler/internal/parser"
)

//1. Add options pattern so it builds out crawler with set of the methods correctly
//3: Improve Concurrency
//4: Add logging in crawler, fetcher and html parser

func main() {
	startURL, _ := url.Parse("https://crawlme.monzo.com/")

	allowedSchemes := map[string]struct{}{
		"http":  {},
		"https": {},
	}
	htmlParser := parser.NewHtmlParser(allowedSchemes)
	fetchTimeout := 10 * time.Second
	hf := fetcher.NewHttpFetcher(fetchTimeout, htmlParser)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	crawler := crawler.NewCrawler(hf, 10)
	foundUrls, err := crawler.Crawl(ctx, startURL)

	if err != nil {
		log.Fatalf("Fetch failed: %v", err)
	}

	fmt.Println("Found total of %d", len(foundUrls))
	for _, url := range foundUrls {
		fmt.Println(url)
	}
}
