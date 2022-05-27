package transform

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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const rootURL = "https://doc.traefik.io"

const maxTitleLength = 65

// PageTransform transforms HTML files under a versioned folder.
type PageTransform struct {
	product string
	pattern *regexp.Regexp
}

// NewPageTransform created a new PageTransform.
func NewPageTransform(product string) *PageTransform {
	return &PageTransform{
		product: product,
		pattern: regexp.MustCompile(`^(.*)/(v\d+\.\d+)/(.*\.html)$`),
	}
}

// Match return true if the file is under a versioned folder.
func (t PageTransform) Match(filename string) bool {
	return t.pattern.MatchString(filename)
}

// Apply applies HTML transformations.
func (t PageTransform) Apply(filename string) error {
	versions := t.pattern.FindStringSubmatch(filename)
	if len(versions) < 4 {
		return fmt.Errorf("version not found: %s", filename)
	}

	v := versions[2]

	doc, err := readDocument(filename)
	if err != nil {
		return err
	}

	doc.Find("head").Each(func(i int, s *goquery.Selection) {
		expectedRelPath := versions[3]
		expectedAbs := filepath.Join(versions[1], expectedRelPath)

		if t.product == "traefik" {
			// exception your middlewares (middlewares/foo/index.html -> (middlewares/http/foo/index.html)
			midExp := regexp.MustCompile(`^middlewares/([^/]+/[^/]+.html)`)
			if midExp.MatchString(expectedRelPath) {
				expectedRelPath = "middlewares/http/" + midExp.FindStringSubmatch(expectedRelPath)[1]
			}
		}

		// Add link canonical URL
		if _, err = os.Stat(expectedAbs); err == nil {
			t.addCanonical(s, expectedRelPath)
			log.Printf("[canonical] %s Adding canonical link", filename)
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

			productNameTitleCase := cases.Title(language.English).String(strings.ReplaceAll(t.product, "-", " "))
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

func (t PageTransform) addCanonical(s *goquery.Selection, fp string) {
	link := s.Find(`link[rel="canonical"]`)
	if link == nil || len(link.Nodes) != 0 {
		return
	}

	r, err := url.Parse(rootURL)
	if err != nil {
		log.Printf("ERROR: unable to parse the root URL: %s", rootURL)
		return
	}

	cano, err := r.Parse(path.Join(t.product, filepath.Dir(fp), "/"))
	if err != nil {
		log.Printf("ERROR: unable to create canonical path: %s %s %s", rootURL, t.product, fp)
		return
	}

	s.AppendHtml(fmt.Sprintf(`<link rel="canonical" href=%q />`, strings.TrimSuffix(cano.String(), "/")+"/"))
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
