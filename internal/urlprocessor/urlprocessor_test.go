package urlprocessor

import (
	"net/url"
	"strings"
	"testing"
)

func mustParse(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}

func TestProcessUrl(t *testing.T) {
	processor := NewUrlProcessor("http", "https")
	base := mustParse("https://monzo.com/")

	t.Run("resolves relative path", func(t *testing.T) {
		got, err := processor.ProcessUrl(mustParse("/about"), base)
		if err != nil {
			t.Fatal(err)
		}
		if got.String() != "https://monzo.com/about" {
			t.Errorf("got %q, want %q", got.String(), "https://monzo.com/about")
		}
	})

	t.Run("rejects unsupported scheme", func(t *testing.T) {
		_, err := processor.ProcessUrl(mustParse("mailto:help@monzo.com"), base)
		if err == nil || !strings.Contains(err.Error(), "unsupported scheme") {
			t.Errorf("expected unsupported scheme error, got %v", err)
		}
	})

	t.Run("rejects external domain", func(t *testing.T) {
		_, err := processor.ProcessUrl(mustParse("https://google.com/page"), base)
		if err == nil || !strings.Contains(err.Error(), "Host Domain doesn't match") {
			t.Errorf("expected host mismatch error, got %v", err)
		}
	})

	t.Run("strips fragment", func(t *testing.T) {
		got, err := processor.ProcessUrl(mustParse("/about#team"), base)
		if err != nil {
			t.Fatal(err)
		}
		if got.String() != "https://monzo.com/about" {
			t.Errorf("got %q, want %q", got.String(), "https://monzo.com/about")
		}
	})

	t.Run("normalises empty path to root", func(t *testing.T) {
		got, err := processor.ProcessUrl(mustParse("https://monzo.com"), base)
		if err != nil {
			t.Fatal(err)
		}
		if got.String() != "https://monzo.com/" {
			t.Errorf("got %q, want %q", got.String(), "https://monzo.com/")
		}
	})

	t.Run("host comparison is case insensitive", func(t *testing.T) {
		got, err := processor.ProcessUrl(mustParse("https://MONZO.COM/page"), base)
		if err != nil {
			t.Fatal(err)
		}
		if got.String() != "https://monzo.com/page" {
			t.Errorf("got %q, want %q", got.String(), "https://monzo.com/page")
		}
	})
}
