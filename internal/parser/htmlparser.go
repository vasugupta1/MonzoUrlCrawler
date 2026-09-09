package parser

import (
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type Parser interface {
	Parse(body io.ReadCloser) ([]*url.URL, error)
}

type HtmlParser struct {
}

func NewHtmlParser() *HtmlParser {

	return &HtmlParser{}
}

func (hp *HtmlParser) Parse(body io.ReadCloser) ([]*url.URL, error) {
	var urls []*url.URL
	seenUrls := make(map[string]struct{})
	tokenzer := html.NewTokenizer(body)
	for {
		tokenType := tokenzer.Next()
		switch tokenType {
		case html.ErrorToken:
			if err := tokenzer.Err(); err != io.EOF {
				return urls, err
			}
			return urls, nil
		case html.StartTagToken, html.SelfClosingTagToken:
			if token := tokenzer.Token(); isATag(token) {
				hrefUrl, err := hp.extractHrefUrl(token.Attr)
				if err != nil {
					//log and continue
					continue
				}
				if hrefUrl == nil {
					//log and continue
					continue
				}
				key := hrefUrl.String()
				if _, exists := seenUrls[key]; !exists {
					seenUrls[key] = struct{}{}
					urls = append(urls, hrefUrl)
				}

			}

		}
	}
}

// First determine if the tag starts with a
func isATag(token html.Token) bool {
	return token.DataAtom.String() == "a" || token.Data == "a"
}

// Extract href url out of the attribute
func (hp *HtmlParser) extractHrefUrl(attr []html.Attribute) (*url.URL, error) {
	for _, attr := range attr {
		if key := attr.Key; key == "href" && attr.Val != "" {
			rawhref := attr.Val

			processedHrefUrl, err := hp.processHref(rawhref)
			if err != nil {
				continue
			}
			if processedHrefUrl != nil {
				return processedHrefUrl, nil
			}
		}

	}
	return nil, nil
}

// process href url
func (hp *HtmlParser) processHref(hrefUrl string) (*url.URL, error) {
	parsedHref, err := url.Parse(strings.TrimSpace(hrefUrl))
	if err != nil {
		return nil, err
	}

	return parsedHref, nil
}
