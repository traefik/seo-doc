package sitemap

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_extractNewItems(t *testing.T) {
	file, err := os.Open("./fixtures/diff.txt")
	require.NoError(t, err)

	defer func() { _ = file.Close() }()

	items, err := extractNewItems(file)
	require.NoError(t, err)

	if os.Getenv("UPDATE_GOLDEN") != "" {
		gg, errG := os.Create("./fixtures/items.golden.json")
		require.NoError(t, errG)

		defer func() { _ = gg.Close() }()

		encoder := json.NewEncoder(gg)
		encoder.SetIndent("", "  ")
		errG = encoder.Encode(items)
		require.NoError(t, errG)
	}

	golden, err := os.Open("./fixtures/items.golden.json")
	require.NoError(t, err)

	defer func() { _ = golden.Close() }()

	expected := make(map[string]Item)
	err = json.NewDecoder(golden).Decode(&expected)
	require.NoError(t, err)

	assert.Equal(t, expected, items)
}

func Test_merge(t *testing.T) {
	itemsFile, err := os.Open("./fixtures/items.golden.json")
	require.NoError(t, err)

	items := make(map[string]Item)
	err = json.NewDecoder(itemsFile).Decode(&items)
	require.NoError(t, err)

	file, err := os.Open("./fixtures/sitemap-5dd7f130.xml")
	require.NoError(t, err)

	var us URLSet
	err = xml.NewDecoder(file).Decode(&us)
	require.NoError(t, err)

	set := merge(us, items)

	if os.Getenv("UPDATE_GOLDEN") != "" {
		gg, errG := os.Create("./fixtures/sitemap.golden.xml")
		require.NoError(t, errG)

		defer func() { _ = gg.Close() }()

		encoder := xml.NewEncoder(gg)
		encoder.Indent("", " ")
		errG = encoder.Encode(set)
		require.NoError(t, errG)
	}

	golden, err := os.Open("./fixtures/sitemap.golden.xml")
	require.NoError(t, err)

	defer func() { _ = golden.Close() }()

	var expected URLSet
	err = xml.NewDecoder(golden).Decode(&expected)
	require.NoError(t, err)

	assert.Equal(t, expected, set)
}
