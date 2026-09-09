package urlprocessor

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

type Processor interface {
	ProcessUrl(crawledUrl *url.URL, targetUrl *url.URL) (*url.URL, error)
}

type Option func(*UrlProcessor)

type UrlProcessor struct {
	supportScheme map[string]struct{}
	logger        *log.Logger
}

func NewUrlProcessor(opts ...Option) *UrlProcessor {
	p := &UrlProcessor{
		supportScheme: map[string]struct{}{},
	}
	for _, opt := range opts {
		opt(p)
	}

	return p
}

func WithLogger(logger *log.Logger) Option {
	return func(p *UrlProcessor) {
		p.logger = logger
	}
}

func WithSupportedScheme(scheme string) Option {
	return func(p *UrlProcessor) {
		p.supportScheme[scheme] = struct{}{}
	}
}

func (p *UrlProcessor) ProcessUrl(crawledUrl *url.URL, baseUrl *url.URL) (*url.URL, error) {
	resolvedHref := baseUrl.ResolveReference(crawledUrl)

	scheme := strings.ToLower(resolvedHref.Scheme)
	if _, ok := p.supportScheme[scheme]; !ok {
		p.logger.Printf("Skipping url as its not supported by scheme configured: %s", crawledUrl)
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
