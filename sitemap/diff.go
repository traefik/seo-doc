package sitemap

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Item a sitemap item.
type Item struct {
	Status  string    `json:"status,omitempty"`
	Path    string    `json:"path,omitempty"`
	Date    time.Time `json:"date,omitempty"`
	Product string    `json:"product,omitempty"`
	Version string    `json:"version,omitempty"`
}

// NewItem creates a new Item.
func NewItem(status, path string, currentDate time.Time) Item {
	return Item{
		Status: status,
		Path:   baseURL + path,
		Date:   currentDate,
	}
}

// FromDiff creates a sitemap from a diff.
func FromDiff(src string) (URLSet, error) {
	// Reads existing sitemap.xml file.
	file, err := os.Open(src)
	if err != nil {
		return URLSet{}, err
	}

	var us URLSet
	err = xml.NewDecoder(file).Decode(&us)
	if err != nil {
		return URLSet{}, err
	}

	// Extract new items.
	date := time.Now().Add(-48 * time.Hour)

	data, err := gitLog(filepath.Dir(src), date)
	if err != nil {
		return URLSet{}, err
	}

	items, err := extractNewItems(data)
	if err != nil {
		return URLSet{}, err
	}

	log.Println("current", len(us.URL), len(items))

	// Merge.
	set := merge(us, items)

	log.Println("new", len(set.URL), len(items))

	return set, nil
}

func gitLog(root string, after time.Time) (*bytes.Reader, error) {
	cmd := exec.Command("git", "log",
		"--name-status",
		"--oneline",
		"--reverse",
		"--after="+after.Format("2006-01-02"),
		"--format=%h %cd",
		"--date=format:%Y-%m-%dT%H:%M:%S",
	)

	cmd.Dir = root

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(output), nil
}

func extractNewItems(data io.Reader) (map[string]Item, error) {
	expDate := regexp.MustCompile(`^[^\s]+\s+(\d{4}-\d{2}-\d{2}T\d{1,2}:\d{1,2}:\d{1,2})$`)
	exp := regexp.MustCompile(`^([MADU])\s+([a-z][^/]+.+/)index\.html$`)
	expSpe := regexp.MustCompile(`^([CR])\d+\s+([a-z][^/]+.+/)index\.html\s+([a-z][^/]+.+/)index\.html$`)
	expVersion := regexp.MustCompile(`^([^/]+)/(master|v\d+\.\d+)/(?:.+/)?$`)

	uniqStatus := make(map[string]Item)

	var currentDate time.Time

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if expDate.MatchString(line) {
			d, err := time.Parse("2006-01-02T15:04:05", expDate.FindStringSubmatch(line)[1])
			if err != nil {
				return nil, err
			}
			currentDate = d
		}

		if !strings.HasSuffix(line, "index.html") {
			continue
		}

		if expSpe.MatchString(line) {
			submatch := expSpe.FindStringSubmatch(line)

			if !expVersion.MatchString(submatch[2]) {
				uniqStatus[baseURL+submatch[2]] = NewItem("D", submatch[2], currentDate)
			}
			if !expVersion.MatchString(submatch[3]) {
				uniqStatus[baseURL+submatch[3]] = NewItem("A", submatch[3], currentDate)
			}

			continue
		}

		if !exp.MatchString(line) {
			continue
		}

		submatch := exp.FindStringSubmatch(line)

		if expVersion.MatchString(submatch[2]) {
			continue
		}

		uniqStatus[baseURL+submatch[2]] = NewItem(submatch[1], submatch[2], currentDate)
	}

	return uniqStatus, nil
}

func merge(us URLSet, items map[string]Item) URLSet {
	var smurls []SMUrl

	for _, u := range us.URL {
		item, ok := items[u.Loc]
		if !ok {
			u.ChangeFreq = changeFreqDaily

			smurls = append(smurls, u)
			continue
		}

		if strings.EqualFold(item.Status, "D") {
			// Deleted item.
			delete(items, u.Loc)
			continue
		}

		smurl := SMUrl{
			Loc:        item.Path,
			LastMod:    item.Date.Format("2006-01-02"),
			ChangeFreq: changeFreqDaily,
		}

		delete(items, u.Loc)
		smurls = append(smurls, smurl)
	}

	for _, item := range items {
		if strings.EqualFold(item.Status, "D") {
			// Deleted item.
			continue
		}

		smurl := SMUrl{
			Loc:        item.Path,
			LastMod:    item.Date.Format("2006-01-02"),
			ChangeFreq: changeFreqDaily,
		}

		smurls = append(smurls, smurl)
	}

	sort.Slice(smurls, func(i, j int) bool {
		return smurls[i].Loc < smurls[j].Loc
	})

	us.URL = smurls

	return us
}
