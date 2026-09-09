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

	resolvedHref.Fragment = ""
	resolvedHref.RawFragment = ""

	if resolvedHref.Path == "" {
		resolvedHref.Path = "/"
	}

	return resolvedHref, nil
}
