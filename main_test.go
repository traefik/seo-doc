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

func Test_getProductName(t *testing.T) {
	testCases := []struct {
		desc     string
		cfg      Config
		expected string
	}{
		{
			desc: "take product option if provided",
			cfg: Config{
				Path:    "/path/to/doc/traefik",
				Product: "test",
			},
			expected: "test",
		},
		{
			desc: "take path option if no product option provided",
			cfg: Config{
				Path:    "/path/to/doc/traefik",
				Product: "",
			},
			expected: "traefik",
		},
		{
			desc: "no product option, no path option",
			cfg: Config{
				Path:    "",
				Product: "",
			},
			expected: ".",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			productName := getProductName(test.cfg)
			assert.Equal(t, test.expected, productName)
		})
	}
}
