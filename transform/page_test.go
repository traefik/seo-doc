package transform

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
			path:   "/v2.4/index.html",
			assert: assert.True,
		},
		{
			path:   "v2.4/index.html",
			assert: assert.False,
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

	testCases := []struct {
		desc    string
		src     string
		product string
		update  bool
	}{
		{
			desc:    "without canonical",
			src:     "index.html",
			product: "test",
		},
		{
			desc:    "without canonical, add sub-folder",
			src:     "foo/index.html",
			product: "test",
		},
		{
			desc:    "with canonical",
			src:     "index1.html",
			product: "test",
		},
		{
			desc:    "with long title",
			src:     "index2.html",
			product: "test",
		},
		{
			desc:    "middlewares rule",
			src:     "middlewares/foo/index.html",
			product: "traefik",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			root := t.TempDir()

			// Creates a fake latest version.
			copyFile(t, "index.html", "", root)
			copyFile(t, "foo/index.html", "", root)
			copyFile(t, "middlewares/foo/index.html", "", root)
			copyFile(t, "middlewares/http/foo/index.html", "", root)

			file := copyFile(t, test.src, "v1.0", root)

			transform := NewPageTransform(test.product)

			err := transform.Apply(file)
			require.NoError(t, err)

			compareFile(t, filepath.Join("./fixtures/output/", test.src), file, test.update)
		})
	}
}

func copyFile(t *testing.T, src, v, root string) string {
	t.Helper()

	if root == "" {
		root = t.TempDir()
	}

	if v != "" {
		root = filepath.Join(root, v)
	}

	dst := filepath.Join(root, src)

	err := os.MkdirAll(filepath.Dir(dst), 0o700)
	require.NoError(t, err)

	in, err := os.Open(filepath.Join("./fixtures/input/", src))
	require.NoError(t, err)
	defer func() { _ = in.Close() }()

	out, err := os.Create(dst)
	require.NoError(t, err)
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, in)
	require.NoError(t, err)

	return dst
}

func compareFile(t *testing.T, src, dst string, update bool) {
	t.Helper()

	a, err := os.ReadFile(src)
	require.NoError(t, err)

	if update {
		err = os.MkdirAll(filepath.Dir(dst), 0o700)
		require.NoError(t, err)

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
		FromFile: "Actual",
		ToFile:   "Expected",
		Context:  2,
	}
	text, err := difflib.GetUnifiedDiffString(diff)
	require.NoError(t, err)

	assert.Empty(t, text)
}
