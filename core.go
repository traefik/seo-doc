package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
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
	htmlFileRegex            = regexp.MustCompile(`^.*\.html$`)
	htmlUnderVersionRegex    = regexp.MustCompile(`^.*\/v\d+\.\d+\/.*\.html$`)
	sitemapUnderVersionRegex = regexp.MustCompile(`\.*/v\d+\.\d+\/.*sitemap\.xml(.gz)?`)
)

func run(cfg Config) error {
	// Extract product name
	productName := filepath.Base(cfg.Path)

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

			if shouldProcessFile(path, htmlFileRegex) {
				return htmlFile(path)
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
			fmt.Printf("[canonical] %s Adding canonical link\n", path)
		}

		// Add meta no follow
		meta := s.Find(`meta[name="robots"][content="index, nofollow"]`)
		if meta != nil && len(meta.Nodes) == 0 {
			s.AppendHtml(`<meta name="robots" content="index, nofollow" />`)
			fmt.Printf("[robots] %s Adding meta robots\n", path)
		}

		// Adds a Suffix in a format | product-name | version
		title := s.Find(`title`)
		if title != nil {
			productNameTitleCase := strings.Title(strings.ReplaceAll(productName, "-", " "))
			suffix := fmt.Sprintf(`| %s | %s`, productNameTitleCase, version)
			titleText := title.Text()
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

func htmlFile(path string) error {
	doc, err := readFile(path)
	if err != nil {
		return err
	}

	// Adds to a document a meta description if available as a hidden input.
	doc.Find(`#meta-description`).Each(func(i int, content *goquery.Selection) {
		doc.Find("head").Each(func(i int, s *goquery.Selection) {
			desc := s.Find(`meta[name="description"]`)
			if desc != nil {
				fmt.Printf("[description] %s Updating meta description\n", path)
				if v, ok := content.Attr("value"); ok {
					desc.SetAttr("content", v)
				}
			} else {
				fmt.Printf("[description] %s Creating meta description\n", path)
				if v, ok := content.Attr("value"); ok {
					s.AppendHtml(fmt.Sprintf(`<meta name="description" content="%s" />`, v))
				}
			}
		})
	})

	return writeFile(path, doc)
}

func sitemapUnderVersion(path string) error {
	fmt.Printf("[sitemap] %s deleted\n", path)
	return os.Remove(path)
}

func writeFile(path string, doc *goquery.Document) error {
	html, err := doc.Html()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, []byte(html), os.ModeAppend)
}

func readFile(path string) (*goquery.Document, error) {
	f, err := os.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return nil, err
	}

	return goquery.NewDocumentFromReader(bufio.NewReader(f))
}

func shouldProcessFile(filePath string, includePattern *regexp.Regexp) bool {
	if includePattern == nil {
		return true
	}

	return includePattern.MatchString(filePath)
}
