package parser

import (
	"io"
	"strings"
	"testing"
)

func body(html string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(html))
}

func TestParse(t *testing.T) {
	parser := NewHtmlParser()

	t.Run("extracts href from anchor tag", func(t *testing.T) {
		urls, err := parser.Parse(body(`<html><body><a href="/about">About</a></body></html>`))
		if err != nil {
			t.Fatal(err)
		}
		if len(urls) != 1 || urls[0].String() != "/about" {
			t.Errorf("got %v, want [/about]", urls)
		}
	})

	t.Run("extracts multiple hrefs", func(t *testing.T) {
		urls, err := parser.Parse(body(`<a href="/a">A</a><a href="/b">B</a>`))
		if err != nil {
			t.Fatal(err)
		}
		if len(urls) != 2 {
			t.Errorf("got %d urls, want 2", len(urls))
		}
	})

	t.Run("deduplicates same href", func(t *testing.T) {
		urls, err := parser.Parse(body(`<a href="/dup">1</a><a href="/dup">2</a>`))
		if err != nil {
			t.Fatal(err)
		}
		if len(urls) != 1 {
			t.Errorf("got %d urls, want 1", len(urls))
		}
	})

	t.Run("skips empty href", func(t *testing.T) {
		urls, err := parser.Parse(body(`<a href="">empty</a>`))
		if err != nil {
			t.Fatal(err)
		}
		if len(urls) != 0 {
			t.Errorf("got %d urls, want 0", len(urls))
		}
	})

	t.Run("skips anchor without href", func(t *testing.T) {
		urls, err := parser.Parse(body(`<a name="top">no href</a>`))
		if err != nil {
			t.Fatal(err)
		}
		if len(urls) != 0 {
			t.Errorf("got %d urls, want 0", len(urls))
		}
	})

	t.Run("ignores non-anchor tags", func(t *testing.T) {
		urls, err := parser.Parse(body(`<div href="/nope"></div><link href="/style.css">`))
		if err != nil {
			t.Fatal(err)
		}
		if len(urls) != 0 {
			t.Errorf("got %d urls, want 0", len(urls))
		}
	})

	t.Run("returns empty for no links", func(t *testing.T) {
		urls, err := parser.Parse(body(`<html><body><p>Hello</p></body></html>`))
		if err != nil {
			t.Fatal(err)
		}
		if len(urls) != 0 {
			t.Errorf("got %d urls, want 0", len(urls))
		}
	})

	t.Run("trims whitespace from href", func(t *testing.T) {
		urls, err := parser.Parse(body(`<a href="  /spaced  ">link</a>`))
		if err != nil {
			t.Fatal(err)
		}
		if len(urls) != 1 || urls[0].String() != "/spaced" {
			t.Errorf("got %v, want [/spaced]", urls)
		}
	})
}
