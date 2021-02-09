package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_validate(t *testing.T) {
	testCases := []struct {
		desc     string
		cfg      Config
		expected string
	}{
		{
			desc: "success",
			cfg: Config{
				Path: "path",
			},
		},
		{
			desc: "missing path",
			cfg: Config{
				Path: "",
			},
			expected: "path is required",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			err := validate(test.cfg)

			if test.expected != "" {
				assert.EqualError(t, err, test.expected)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
