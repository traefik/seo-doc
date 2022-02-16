package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSitemapTransform_Match(t *testing.T) {
	transform := NewSitemapTransform("test")

	testCases := []struct {
		desc   string
		path   string
		assert assert.BoolAssertionFunc
	}{
		{
			desc:   "match sitemap under version",
			path:   "foo/v2.4/sitemap.xml",
			assert: assert.True,
		},
		{
			desc:   "match sitemap under version",
			path:   "foo/v2.4/sitemap.xml.gz",
			assert: assert.True,
		},
		{
			desc:   "does not match sitemap under version",
			path:   "foo/v2.4/powpow.xml.gz",
			assert: assert.False,
		},
		{
			desc:   "does not match sitemap under version",
			path:   "/powpow.xml.gz",
			assert: assert.False,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			test.assert(t, transform.Match(test.path))
		})
	}
}

func TestSitemapTransform_Apply(t *testing.T) {
	transform := NewSitemapTransform("test")

	testCases := []struct {
		desc string
		path string
	}{
		{
			desc: "",
			path: "./fixtures/input/sitemap.xml",
		},
		{
			desc: "",
			path: "./fixtures/input/sitemap.xml.gz",
		},
		{
			desc: "",
			path: "./fixtures/input/index.html",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			file := copyFile(t, test.path)

			err := transform.Apply(file)
			require.NoError(t, err)

			assert.NoFileExists(t, file)
		})
	}
}
