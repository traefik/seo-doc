package sitemap

import (
	"compress/gzip"
	"encoding/xml"
	"io"
	"log"
	"os"
	"path/filepath"
)

const baseURL = "https://doc.traefik.io/"

const changeFreqDaily = "daily"

// URLSet root of a sitemap file.
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URL     []SMUrl  `xml:"url"`
}

// SMUrl item of a sitemap file.
type SMUrl struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   int    `xml:"priority,omitempty"`
}

// Generate generates sitemap files.
func Generate(root string) error {
	src := filepath.Join(root, "sitemap.xml")

	var set URLSet
	if _, err := os.Stat(src); err != nil {
		log.Println("From scratch", src)

		set, err = FromScratch(root)
		if err != nil {
			return err
		}
	} else {
		log.Println("From diff", src)

		set, err = FromDiff(src)
		if err != nil {
			return err
		}
	}

	return saveSitemap(src, set)
}

func saveSitemap(dst string, set URLSet) error {
	file, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	gz, err := os.Create(dst + ".gz")
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	zw := gzip.NewWriter(gz)
	zw.Name = "sitemap.xml"

	defer func() { _ = zw.Close() }()

	output := io.MultiWriter(file, zw)

	_, err = output.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>`))
	if err != nil {
		return err
	}

	_, err = output.Write([]byte("\n"))
	if err != nil {
		return err
	}

	encoder := xml.NewEncoder(output)
	encoder.Indent("", "  ")
	err = encoder.Encode(set)
	if err != nil {
		return err
	}

	_, err = output.Write([]byte("\n"))
	if err != nil {
		return err
	}

	return nil
}
