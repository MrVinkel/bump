package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVersion(t *testing.T) {
	type test struct {
		name     string
		version  string
		expected *Version
		err      bool
	}

	tests := []test{
		{
			name:    "valid version",
			version: "1.2.3",
			expected: &Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
		},
		{
			name:    "valid version double digit",
			version: "56.43.32",
			expected: &Version{
				Major: 56,
				Minor: 43,
				Patch: 32,
			},
		},
		{
			name:    "valid version with prefix",
			version: "v1.2.3",
			expected: &Version{
				Prefix: Ptr("v"),
				Major:  1,
				Minor:  2,
				Patch:  3,
			},
		},
		{
			name:    "valid version with prefix special character",
			version: "abc!\"#¤%&/()=?-_,.'¨^1.2.3",
			expected: &Version{
				Prefix: Ptr("abc!\"#¤%&/()=?-_,.'¨^"),
				Major:  1,
				Minor:  2,
				Patch:  3,
			},
		},
		{
			name:    "invalid version",
			version: "asdf",
			err:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ParseVersion(tc.version)
			if tc.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.String(), actual.String())
			}
		})
	}
}
