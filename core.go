package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	rootURL        = "https://doc.traefik.io"
	maxTitleLength = 65
)

var (
	versionRegex             = regexp.MustCompile(`^.*\/(v\d+\.\d+)\/.*$`)
	htmlUnderVersionRegex    = regexp.MustCompile(`^.*\/v\d+\.\d+\/.*\.html$`)
	sitemapUnderVersionRegex = regexp.MustCompile(`\.*/v\d+\.\d+\/.*sitemap\.xml(.gz)?`)
)

func run(cfg Config) error {
	// Extract product name
	productName := cfg.Product
	if productName == "" {
		productName = filepath.Base(cfg.Path)
	}

	err := filepath.Walk(cfg.Path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			versions := versionRegex.FindStringSubmatch(path)

			version := ""
			if len(versions) > 1 {
				version = versions[1]
			}

			if shouldProcessFile(path, htmlUnderVersionRegex) {
				return htmlFileUnderVersion(path, productName, version)
			}

			if shouldProcessFile(path, sitemapUnderVersionRegex) {
				return sitemapUnderVersion(path)
			}

			return nil
		})

	return err
}

func htmlFileUnderVersion(path string, productName string, version string) error {
	doc, err := readFile(path)
	if err != nil {
		return err
	}

	doc.Find("head").Each(func(i int, s *goquery.Selection) {
		// Add link canonical URL
		link := s.Find(`link[rel="canonical"]`)
		if link != nil && len(link.Nodes) == 0 {
			s.AppendHtml(fmt.Sprintf(`<link rel="canonical" href="%s/%s/" />`, rootURL, productName))
			log.Printf("[canonical] %s Adding canonical link", path)
		}

		// Add meta no follow
		meta := s.Find(`meta[name="robots"][content="index, nofollow"]`)
		if meta != nil && len(meta.Nodes) == 0 {
			s.AppendHtml(`<meta name="robots" content="index, nofollow" />`)
			log.Printf("[robots] %s Adding meta robots", path)
		}

		// Adds a Suffix in a format | product-name | version
		title := s.Find(`title`)
		if title != nil {
			titleText := title.Text()

			productNameTitleCase := strings.Title(strings.ReplaceAll(productName, "-", " "))
			suffix := fmt.Sprintf("| %s | %s", productNameTitleCase, version)

			if !strings.Contains(titleText, suffix) {
				newTitle := fmt.Sprintf("%s %s", strings.ReplaceAll(titleText, fmt.Sprintf(` - %s`, productNameTitleCase), ""), suffix)
				if len(newTitle) > maxTitleLength {
					maxNewTitleLength := maxTitleLength - len(suffix)
					newTitle = fmt.Sprintf("%s... %s", titleText[:maxNewTitleLength-4], suffix)
				}

				title.SetText(newTitle)
			}
		}
	})

	return writeFile(path, doc)
}

func sitemapUnderVersion(path string) error {
	log.Printf("[sitemap] %s deleted", path)
	return os.Remove(path)
}

func writeFile(path string, doc *goquery.Document) error {
	html, err := doc.Html()
	if err != nil {
		return err
	}

	html = strings.ReplaceAll(html, `src="http://`, `src="https://`)
	html = strings.ReplaceAll(html, `href="http://`, `href="https://`)

	return ioutil.WriteFile(path, []byte(html), os.ModeAppend)
}

func readFile(path string) (*goquery.Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() { _ = f.Close() }()

	return goquery.NewDocumentFromReader(bufio.NewReader(f))
}

func shouldProcessFile(filePath string, includePattern *regexp.Regexp) bool {
	if includePattern == nil {
		return true
	}

	return includePattern.MatchString(filePath)
}
