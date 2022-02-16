package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const rootURL = "https://doc.traefik.io"

const maxTitleLength = 65

// VersionTransform transforms files under a versioned folder.
type VersionTransform struct {
	product        string
	pattern        *regexp.Regexp
	versionPattern *regexp.Regexp
}

// NewVersionTransform created a new VersionTransform.
func NewVersionTransform(product string) *VersionTransform {
	return &VersionTransform{
		product:        product,
		pattern:        regexp.MustCompile(`^.*/v\d+\.\d+/.*\.html$`),
		versionPattern: regexp.MustCompile(`^.*/(v\d+\.\d+)/.*$`),
	}
}

// Match return true if the file is under a versioned folder.
func (t VersionTransform) Match(path string) bool {
	return t.pattern.MatchString(path)
}

// Apply applies HTML transformations.
func (t VersionTransform) Apply(path string) error {
	versions := t.versionPattern.FindStringSubmatch(path)

	v := ""
	if len(versions) > 1 {
		v = versions[1]
	}

	doc, err := readDocument(path)
	if err != nil {
		return err
	}

	doc.Find("head").Each(func(i int, s *goquery.Selection) {
		// Add link canonical URL
		link := s.Find(`link[rel="canonical"]`)
		if link != nil && len(link.Nodes) == 0 {
			s.AppendHtml(fmt.Sprintf(`<link rel="canonical" href="%s/%s/" />`, rootURL, t.product))
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

			productNameTitleCase := strings.Title(strings.ReplaceAll(t.product, "-", " "))
			suffix := fmt.Sprintf("| %s | %s", productNameTitleCase, v)

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

func writeFile(path string, doc *goquery.Document) error {
	html, err := doc.Html()
	if err != nil {
		return err
	}

	replacer := strings.NewReplacer(
		`src="http://`, `src="https://`,
		`href="http://`, `href="https://`,
	)
	html = replacer.Replace(html)

	return os.WriteFile(path, []byte(html), os.ModeAppend)
}

func readDocument(path string) (*goquery.Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() { _ = f.Close() }()

	return goquery.NewDocumentFromReader(bufio.NewReader(f))
}
