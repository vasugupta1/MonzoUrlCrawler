package parser

import (
	"io"

	"golang.org/x/net/html"
)

type Parser interface {
	Parse(body io.ReadCloser, url string) ([]string, error)
}

type HtmlParser struct {
}

func NewHtmlParser() *HtmlParser {
	return &HtmlParser{}
}

// Need to figure out the best wato extract hrefs from the body, this needs to be somewhat efficent
func (hp *HtmlParser) Parse(body io.ReadCloser, url string) ([]string, error) {
	var links []string
	tokenzer := html.NewTokenizer(body)

	return links, nil
}
