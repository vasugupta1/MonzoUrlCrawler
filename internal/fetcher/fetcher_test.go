package fetcher

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestFetch(t *testing.T) {
	fetcher := NewHttpFetcher(5 * time.Second)

	t.Run("returns body on 200", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello"))
		}))
		defer srv.Close()

		u, _ := url.Parse(srv.URL)
		body, err := fetcher.Fetch(context.Background(), u)
		if err != nil {
			t.Fatal(err)
		}
		defer body.Close()

		data, _ := io.ReadAll(body)
		if string(data) != "hello" {
			t.Errorf("got %q, want %q", string(data), "hello")
		}
	})

	t.Run("returns error on non-200 status", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		u, _ := url.Parse(srv.URL)
		_, err := fetcher.Fetch(context.Background(), u)
		if err == nil {
			t.Fatal("expected error for 404 status")
		}
	})

	t.Run("returns error on cancelled context", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		}))
		defer srv.Close()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		u, _ := url.Parse(srv.URL)
		_, err := fetcher.Fetch(ctx, u)
		if err == nil {
			t.Fatal("expected error for cancelled context")
		}
	})

	t.Run("returns error for unreachable server", func(t *testing.T) {
		u, _ := url.Parse("http://127.0.0.1:1")
		_, err := fetcher.Fetch(context.Background(), u)
		if err == nil {
			t.Fatal("expected connection error")
		}
	})
}
