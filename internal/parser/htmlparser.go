package parser

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type Parser interface {
	Parse(body io.ReadCloser, baseUrl *url.URL) ([]*url.URL, error)
}

type HtmlParser struct {
	supportScheme map[string]struct{}
}

func NewHtmlParser(supportScheme map[string]struct{}) *HtmlParser {
	if len(supportScheme) == 0 {
		supportScheme = map[string]struct{}{
			"http":  {},
			"https": {},
		}
	}
	return &HtmlParser{
		supportScheme: supportScheme,
	}
}

func (hp *HtmlParser) Parse(body io.ReadCloser, baseUrl *url.URL) ([]*url.URL, error) {
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
				hrefUrl, err := hp.extractHrefUrl(token.Attr, baseUrl)
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
func (hp *HtmlParser) extractHrefUrl(attr []html.Attribute, baseUrl *url.URL) (*url.URL, error) {
	for _, attr := range attr {
		if key := attr.Key; key == "href" && attr.Val != "" {
			rawhref := attr.Val

			processedHrefUrl, err := hp.processHref(rawhref, baseUrl)
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

func (hp *HtmlParser) processHref(hrefUrl string, baseUrl *url.URL) (*url.URL, error) {
	parsedHref, err := url.Parse(strings.TrimSpace(hrefUrl))
	if err != nil {
		return nil, err
	}

	resolvedHref := baseUrl.ResolveReference(parsedHref)

	scheme := strings.ToLower(resolvedHref.Scheme)

	if _, ok := hp.supportScheme[scheme]; !ok {
		return nil, fmt.Errorf("unsupported scheme: %s", scheme)
	}

	resolvedHref.Host = strings.ToLower(resolvedHref.Host)

	if resolvedHref.Host != strings.ToLower(baseUrl.Host) {
		return nil, fmt.Errorf("Base doesn't match hence skip")
	}

	resolvedHref.Fragment = ""
	resolvedHref.RawFragment = ""

	if resolvedHref.Path == "" {
		resolvedHref.Path = "/"
	}

	return resolvedHref, nil
}
