package parser

import "io"

type Parser interface {
	Parse(body io.ReadCloser, url string) ([]string, error)
}

type HtmlParser struct {
}

func NewHtmlParser() *HtmlParser {
	return &HtmlParser{}
}

func (hp *HtmlParser) Parse(body io.ReadCloser, url string) ([]string, error) {

}
