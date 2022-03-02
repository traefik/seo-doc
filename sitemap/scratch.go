package sitemap

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// FromScratch creates a sitemap from scratch.
func FromScratch(root string) (URLSet, error) {
	exp := regexp.MustCompile(`^(.+/)index\.html$`)
	expVersion := regexp.MustCompile(`^([^/]+)/(master|v\d+\.\d+)/.*$`)

	us := URLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}

	errW := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		if !exp.MatchString(path) {
			return nil
		}

		urlPath := strings.TrimPrefix(strings.TrimSuffix(path, "index.html"), root)

		if expVersion.MatchString(urlPath) {
			return nil
		}

		us.URL = append(us.URL, SMUrl{
			Loc:        baseURL + urlPath,
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: changeFreqDaily,
		})

		return nil
	})
	if errW != nil {
		return URLSet{}, errW
	}

	sort.Slice(us.URL, func(i, j int) bool {
		return us.URL[i].Loc < us.URL[j].Loc
	})

	return us, nil
}
