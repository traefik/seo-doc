package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPageTransform_Match(t *testing.T) {
	transform := NewPageTransform("test")

	testCases := []struct {
		path   string
		assert assert.BoolAssertionFunc
	}{
		{
			path:   "foo/v2.4/index.html",
			assert: assert.True,
		},
		{
			path:   "foo/index.html",
			assert: assert.False,
		},
		{
			path:   "foo/v2/index.html",
			assert: assert.False,
		},
		{
			path:   "foo/v2.4-bar/index.html",
			assert: assert.False,
		},
		{
			path:   "foo/v2.4/index.js",
			assert: assert.False,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.path, func(t *testing.T) {
			t.Parallel()

			test.assert(t, transform.Match(test.path))
		})
	}
}

func TestPageTransform_Apply(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("windows")
	}

	transform := NewPageTransform("test")

	testCases := []struct {
		desc   string
		src    string
		dst    string
		update bool
	}{
		{
			desc: "without canonical",
			src:  "./fixtures/input/index.html",
			dst:  "./fixtures/output/index.html",
		},
		{
			desc: "with canonical",
			src:  "./fixtures/input/index1.html",
			dst:  "./fixtures/output/index1.html",
		},
		{
			desc: "with long title",
			src:  "./fixtures/input/index2.html",
			dst:  "./fixtures/output/index2.html",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			file := copyFile(t, test.src)

			err := transform.Apply(file)
			require.NoError(t, err)

			compareFile(t, file, test.dst, test.update)
		})
	}
}

func copyFile(t *testing.T, src string) string {
	t.Helper()

	source, err := os.Open(src)
	require.NoError(t, err)
	defer func() { _ = source.Close() }()

	temp := filepath.Join(t.TempDir(), "v1.0")
	err = os.MkdirAll(temp, 0o700)
	require.NoError(t, err)

	dst := filepath.Join(temp, filepath.Base(src))

	destination, err := os.Create(dst)
	require.NoError(t, err)
	defer func() { _ = destination.Close() }()

	_, err = io.Copy(destination, source)
	require.NoError(t, err)

	return dst
}

func compareFile(t *testing.T, src, dst string, update bool) {
	t.Helper()

	a, err := os.ReadFile(src)
	require.NoError(t, err)

	if update {
		var out *os.File
		out, err = os.Create(dst)
		require.NoError(t, err)

		t.Cleanup(func() {
			_ = out.Close()
		})

		_, err = out.Write(a)
		require.NoError(t, err)
	}

	b, err := os.ReadFile(dst)
	require.NoError(t, err)

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(bytes.TrimSpace(a))),
		B:        difflib.SplitLines(string(bytes.TrimSpace(b))),
		FromFile: "Original",
		ToFile:   "Current",
		Context:  2,
	}
	text, err := difflib.GetUnifiedDiffString(diff)
	require.NoError(t, err)

	assert.Empty(t, text)
}
