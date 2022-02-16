package main

import (
	"log"
	"os"
	"regexp"
)

// SitemapTransform transforms sitemap files.
type SitemapTransform struct {
	pattern *regexp.Regexp
	product string
}

// NewSitemapTransform created a new SitemapTransform.
func NewSitemapTransform(product string) *SitemapTransform {
	return &SitemapTransform{
		product: product,
		pattern: regexp.MustCompile(`\.*/v\d+\.\d+/.*sitemap\.xml(\.gz)?$`),
	}
}

// Match return true if the file is a sitemap related file.
func (t SitemapTransform) Match(path string) bool {
	return t.pattern.MatchString(path)
}

// Apply removes a file.
func (t SitemapTransform) Apply(path string) error {
	// Remove sitemap files for versioned documentation.
	log.Printf("[sitemap] %s deleted", path)
	return os.Remove(path)
}
