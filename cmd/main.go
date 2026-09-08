package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/vasugupta1/MonzoUrlCrawler/internal/fetcher"
	"github.com/vasugupta1/MonzoUrlCrawler/internal/parser"
)

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

	links, err := hf.Fetch(ctx, startURL)
	if err != nil {
		log.Fatalf("Fetch failed: %v", err)
	}

	fmt.Printf("Discovered %d links on %s:\n", len(links), startURL)
	for _, link := range links {
		fmt.Println(" -", link.String())
	}
}
