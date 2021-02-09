package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	rootUrl        = "https://doc.traefik.io"
	maxTitleLength = 65
)

var (
	versionRegex             = regexp.MustCompile(`^.*\/(v\d+\.\d+)\/.*$`)
	htmlFileRegex            = regexp.MustCompile(`^.*\.html$`)
	htmlUnderVersionRegex    = regexp.MustCompile(`^\/v\d+\.\d+\/.*\.html$`)
	sitemapUnderVersionRegex = regexp.MustCompile(`\/v\d+\.\d+\/.*sitemap\.xml(.gz)?`)
)

func run(cfg Config) error {
	// Extract product name
	productName := filepath.Base(cfg.Path)

	err := filepath.Walk(cfg.Path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			versions := versionRegex.FindAllString(path, 1)
			version := ""
			if len(versions) > 0 {
				version = versions[0]
			}

			if shouldProcessFile(path, htmlUnderVersionRegex) {
				return htmlFileUnderVersion(cfg.Path, productName, version)
			}

			if shouldProcessFile(path, htmlFileRegex) {
				return htmlFile(cfg.Path, productName, version)
			}

			if shouldProcessFile(path, sitemapUnderVersionRegex) {
				return sitemapUnderVersion(cfg.Path, productName, version)
			}

			return nil
		})

	return err
}

func htmlFileUnderVersion(path string, productName string, version string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(f))
	if err != nil {
		return err
	}

	doc.Find("head").Each(func(i int, s *goquery.Selection) {
		// Add link canonical URL
		link := s.Find(`link[rel="canonical"]`)
		if link == nil {
			s.AppendHtml(fmt.Sprintf(`<link rel="canonical" href="%s/%s/" />`, rootUrl, productName))
			fmt.Printf("[canonical] %s Adding canonical link\n", path)
		}
		// Add meta no follow
		meta := s.Find(`meta[name="robots"][content="index, nofollow"]`)
		if meta == nil {
			s.AppendHtml(`<meta name="robots" content="index, nofollow" />`)
			fmt.Printf("[robots] %s Adding meta robots\n", path)
		}

		// Adds a Suffix in a format | product-name | version
		title := s.Find(`title`)
		if title != nil {
			productNameTitleCase := strings.ToTitle(strings.ReplaceAll(productName, "-", " "))
			suffix := fmt.Sprintf(` | %s | %s`, productNameTitleCase, version)
			titleText := title.Text()
			// FIXME manage title length
			if !strings.Contains(titleText, suffix) {
				newTitle := strings.ReplaceAll(titleText, fmt.Sprintf(` - %s`, productNameTitleCase), "")
				if len(newTitle) > maxTitleLength {
					newTitle = newTitle[0:]
				}
			}
		}
	})
	return nil
}

func htmlFile(path string, basename string, version string) error {
	return nil
}

func sitemapUnderVersion(path string, basename string, version string) error {
	return nil
}

func shouldProcessFile(filePath string, includePattern *regexp.Regexp) bool {
	if includePattern == nil {
		return true
	}

	return includePattern.MatchString(filePath)
}
