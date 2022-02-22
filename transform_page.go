package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const rootURL = "https://doc.traefik.io"

const maxTitleLength = 65

// envVarLatestTag it's an env vars used by structor to know the latest version.
const envVarLatestTag = "STRUCTOR_LATEST_TAG"

// PageTransform transforms HTML files under a versioned folder.
type PageTransform struct {
	product string
	pattern *regexp.Regexp
	latest  string
}

// NewPageTransform created a new PageTransform.
func NewPageTransform(product string) *PageTransform {
	return &PageTransform{
		product: product,
		pattern: regexp.MustCompile(`^.*/(v\d+\.\d+)/(.*\.html)$`),
		latest:  os.Getenv(envVarLatestTag),
	}
}

// Match return true if the file is under a versioned folder.
func (t PageTransform) Match(filename string) bool {
	return t.pattern.MatchString(filename)
}

// Apply applies HTML transformations.
func (t PageTransform) Apply(filename string) error {
	versions := t.pattern.FindStringSubmatch(filename)
	if len(versions) < 3 {
		return fmt.Errorf("version not found: %s", filename)
	}

	v := versions[1]

	doc, err := readDocument(filename)
	if err != nil {
		return err
	}

	doc.Find("head").Each(func(i int, s *goquery.Selection) {
		// Add link canonical URL
		if strings.HasPrefix(t.latest, v) {
			t.addCanonical(s, filename, versions[2])
		}

		// Add meta no follow
		meta := s.Find(`meta[name="robots"][content="index, nofollow"]`)
		if meta != nil && len(meta.Nodes) == 0 {
			s.AppendHtml(`<meta name="robots" content="index, nofollow" />`)
			log.Printf("[robots] %s Adding meta robots", filename)
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

	return writeFile(filename, doc)
}

func (t PageTransform) addCanonical(s *goquery.Selection, filename, p string) {
	link := s.Find(`link[rel="canonical"]`)
	if link == nil || len(link.Nodes) != 0 {
		return
	}

	r, err := url.Parse(rootURL)
	if err != nil {
		log.Printf("ERROR: unable to parse the root URL: %s", rootURL)
		return
	}

	cano, err := r.Parse(path.Join(t.product, filepath.Dir(p), "/"))
	if err != nil {
		log.Printf("ERROR: unable to create canonical path: %s %s %s", rootURL, t.product, p)
		return
	}

	s.AppendHtml(fmt.Sprintf(`<link rel="canonical" href=%q />`, strings.TrimSuffix(cano.String(), "/")+"/"))
	log.Printf("[canonical] %s Adding canonical link", filename)
}

func writeFile(filename string, doc *goquery.Document) error {
	html, err := doc.Html()
	if err != nil {
		return err
	}

	replacer := strings.NewReplacer(
		`src="http://`, `src="https://`,
		`href="http://`, `href="https://`,
	)
	html = replacer.Replace(html)

	return os.WriteFile(filename, []byte(html), os.ModeAppend)
}

func readDocument(filename string) (*goquery.Document, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer func() { _ = f.Close() }()

	return goquery.NewDocumentFromReader(bufio.NewReader(f))
}
