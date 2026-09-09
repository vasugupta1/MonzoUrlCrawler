package crawler

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"testing"
	"time"

	"github.com/vasugupta1/MonzoUrlCrawler/internal/fetcher"
	"github.com/vasugupta1/MonzoUrlCrawler/internal/parser"
	"github.com/vasugupta1/MonzoUrlCrawler/internal/urlprocessor"
)

func setupCrawler(srv *httptest.Server) (*Crawler, *url.URL) {
	base, _ := url.Parse(srv.URL + "/")
	f := fetcher.NewHttpFetcher(5 * time.Second)
	p := parser.NewHtmlParser()
	proc := urlprocessor.NewUrlProcessor("http", "https")
	c := NewCrawler(f, p, proc,
		WithWorkercount(2),
		WithLogger(log.New(io.Discard, "", 0)),
	)
	return c, base
}

func TestCrawl(t *testing.T) {
	t.Run("crawls linked pages", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/":
				w.Write([]byte(`<a href="/about">about</a><a href="/contact">contact</a>`))
			case "/about", "/contact":
				w.Write([]byte(`<p>no links</p>`))
			}
		}))
		defer srv.Close()

		c, base := setupCrawler(srv)
		results, err := c.Crawl(context.Background(), base)
		if err != nil {
			t.Fatal(err)
		}

		sort.Strings(results)
		want := []string{srv.URL + "/", srv.URL + "/about", srv.URL + "/contact"}
		sort.Strings(want)

		if len(results) != len(want) {
			t.Fatalf("got %d urls, want %d: %v", len(results), len(want), results)
		}
		for i := range want {
			if results[i] != want[i] {
				t.Errorf("results[%d] = %q, want %q", i, results[i], want[i])
			}
		}
	})

	t.Run("does not revisit same page", func(t *testing.T) {
		visitCount := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			visitCount++
			// every page links back to root
			w.Write([]byte(`<a href="/">home</a>`))
		}))
		defer srv.Close()

		c, base := setupCrawler(srv)
		_, err := c.Crawl(context.Background(), base)
		if err != nil {
			t.Fatal(err)
		}

		if visitCount != 1 {
			t.Errorf("visited root %d times, want 1", visitCount)
		}
	})

	t.Run("skips external links", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<a href="https://external.com/page">ext</a>`))
		}))
		defer srv.Close()

		c, base := setupCrawler(srv)
		results, err := c.Crawl(context.Background(), base)
		if err != nil {
			t.Fatal(err)
		}

		if len(results) != 1 {
			t.Errorf("got %d urls, want 1 (just base): %v", len(results), results)
		}
	})

	t.Run("handles fetch errors gracefully", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/":
				w.Write([]byte(`<a href="/broken">link</a>`))
			case "/broken":
				w.WriteHeader(http.StatusInternalServerError)
			}
		}))
		defer srv.Close()

		c, base := setupCrawler(srv)
		results, err := c.Crawl(context.Background(), base)
		if err != nil {
			t.Fatal(err)
		}

		// base + /broken (submitted before fetch fails)
		if len(results) < 1 {
			t.Errorf("expected at least base url in results, got %v", results)
		}
	})

	t.Run("stops on cancelled context", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// infinite chain of links
			w.Write([]byte(`<a href="/a">a</a><a href="/b">b</a><a href="/c">c</a>`))
		}))
		defer srv.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		c, base := setupCrawler(srv)
		_, err := c.Crawl(ctx, base)
		if err != nil {
			t.Fatal(err)
		}
		// just verifying it terminates and doesn't hang
	})
}
