package urlprocessor

import (
	"fmt"
	"net/url"
	"strings"
)

type Processor interface {
	ProcessUrl(crawledUrl *url.URL, targetUrl *url.URL) (*url.URL, error)
}

type Option func(*UrlProcessor)

type UrlProcessor struct {
	supportScheme map[string]struct{}
}

func NewUrlProcessor(supportedSchems ...string) *UrlProcessor {
	schemes := make(map[string]struct{})
	for _, scheme := range supportedSchems {
		schemes[scheme] = struct{}{}
	}
	p := &UrlProcessor{
		supportScheme: schemes,
	}

	return p
}

func (p *UrlProcessor) ProcessUrl(crawledUrl *url.URL, baseUrl *url.URL) (*url.URL, error) {
	resolvedHref := baseUrl.ResolveReference(crawledUrl)

	scheme := strings.ToLower(resolvedHref.Scheme)
	if _, ok := p.supportScheme[scheme]; !ok {
		return nil, fmt.Errorf("unsupported scheme: %s", scheme)
	}

	resolvedHref.Host = strings.ToLower(resolvedHref.Host)
	if resolvedHref.Host != strings.ToLower(baseUrl.Host) {
		return nil, fmt.Errorf("Host Domain doesn't match")
	}

	// page could have fragments which tell the browser where to scroll on the page so its probably best to set them empty
	resolvedHref.Fragment = ""
	resolvedHref.RawFragment = ""

	// normalise the the root incase we get into a loop with <base_url> and <base_url>/
	if resolvedHref.Path == "" {
		resolvedHref.Path = "/"
	}

	return resolvedHref, nil
}
