package main

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_shouldProcessFile(t *testing.T) {
	testCases := []struct {
		desc           string
		filePath       string
		includePattern *regexp.Regexp
		expected       bool
	}{
		{
			desc:           "match html file",
			filePath:       "index.html",
			includePattern: htmlFileRegex,
			expected:       true,
		},
		{
			desc:           "does not match html file",
			filePath:       "index.htm",
			includePattern: htmlFileRegex,
			expected:       false,
		},
		{
			desc:           "match html file under version",
			filePath:       "/v2.4/index.html",
			includePattern: htmlFileRegex,
			expected:       true,
		},
		{
			desc:           "does not match html file under version",
			filePath:       "/v2.4/index.htm",
			includePattern: htmlFileRegex,
			expected:       false,
		},
		{
			desc:           "match sitemap under version",
			filePath:       "foo/v2.4/sitemap.xml",
			includePattern: sitemapUnderVersionRegex,
			expected:       true,
		},
		{
			desc:           "match sitemap under version",
			filePath:       "foo/v2.4/sitemap.xml.gz",
			includePattern: sitemapUnderVersionRegex,
			expected:       true,
		},
		{
			desc:           "does not match sitemap under version",
			filePath:       "foo/v2.4/powpow.xml.gz",
			includePattern: sitemapUnderVersionRegex,
			expected:       false,
		},
		{
			desc:           "does not match sitemap under version",
			filePath:       "/powpow.xml.gz",
			includePattern: sitemapUnderVersionRegex,
			expected:       false,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			match := shouldProcessFile(test.filePath, test.includePattern)

			assert.Equal(t, test.expected, match)
		})
	}
}
