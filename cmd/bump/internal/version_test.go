package internal_test

import (
	"fmt"
	"testing"

	"github.com/mrvinkel/bump/cmd/bump/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVersion(t *testing.T) {
	type test struct {
		name     string
		version  string
		expected *internal.Version
		err      bool
	}

	tests := []test{
		{
			name:    "valid version",
			version: "1.2.3",
			expected: &internal.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
		},
		{
			name:    "valid version prefix",
			version: "v1.2.3",
			expected: &internal.Version{
				Prefix: internal.Ptr("v"),
				Major:  1,
				Minor:  2,
				Patch:  3,
			},
		},
		{
			name:    "valid version component prefix",
			version: "some-component-1.2.3",
			expected: &internal.Version{
				Prefix: internal.Ptr("some-component-"),
				Major:  1,
				Minor:  2,
				Patch:  3,
			},
		},
		{
			name:    "valid version pre release",
			version: "1.2.3-alpha.1",
			expected: &internal.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: []string{"alpha", "1"},
			},
		},
		{
			name:    "valid version build metadata",
			version: "1.2.3+build.123",
			expected: &internal.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: internal.Ptr("build.123"),
			},
		},
		{
			name:    "valid version all fields",
			version: "v1.2.3-beta.2+build.123",
			expected: &internal.Version{
				Prefix:     internal.Ptr("v"),
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: []string{"beta", "2"},
				Build:      internal.Ptr("build.123"),
			},
		},
		{
			name:    "valid version double digit",
			version: "56.43.32",
			expected: &internal.Version{
				Major: 56,
				Minor: 43,
				Patch: 32,
			},
		},
		{
			name:    "valid version with prefix",
			version: "v1.2.3",
			expected: &internal.Version{
				Prefix: internal.Ptr("v"),
				Major:  1,
				Minor:  2,
				Patch:  3,
			},
		},
		{
			name:    "valid version with prefix special character",
			version: "abc!\"#¤%&/()=?-_,.'¨^1.2.3",
			expected: &internal.Version{
				Prefix: internal.Ptr("abc!\"#¤%&/()=?-_,.'¨^"),
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
			actual, err := internal.ParseVersion(tc.version)
			if tc.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.String(), actual.String())
			}
		})
	}
}

func TestCompare(t *testing.T) {
	greater := 1
	less := -1
	equal := 0
	type test struct {
		version1 string
		version2 string
		expected int
	}

	tests := []test{
		{
			version1: "0.0.0",
			version2: "0.0.0",
			expected: equal,
		},
		{
			version1: "0.0.1",
			version2: "0.0.0",
			expected: greater,
		},
		{
			version1: "0.1.0",
			version2: "0.0.1",
			expected: greater,
		},
		{
			version1: "1.0.0",
			version2: "0.1.1",
			expected: greater,
		},
		{
			version1: "1.2.4",
			version2: "1.2.3",
			expected: greater,
		},
		{
			version1: "1.2.3",
			version2: "1.2.4",
			expected: less,
		},
		{
			version1: "1.2.3",
			version2: "1.2.3-alpha",
			expected: greater,
		},
		{
			version1: "1.2.3-alpha",
			version2: "1.2.3",
			expected: less,
		},
		{
			version1: "1.2.3-alpha.1",
			version2: "1.2.3-alpha.1",
			expected: equal,
		},
		{
			version1: "1.2.3-alpha.1",
			version2: "1.2.3-alpha.1+build",
			expected: equal,
		},
		{
			version1: "1.2.3-alpha.1+build",
			version2: "1.2.3-alpha.1",
			expected: equal,
		},
		{
			version1: "1.2.3-alpha.1+build",
			version2: "1.2.3-alpha.1+build.1",
			expected: equal,
		},
		{
			version1: "1.2.3-alpha.1+build.1",
			version2: "1.2.3-alpha.1+build",
			expected: equal,
		},
		{
			version1: "1.2.3-alpha.1",
			version2: "1.2.3-alpha.2",
			expected: less,
		},
		{
			version1: "1.2.3-alpha.2",
			version2: "1.2.3-alpha.1",
			expected: greater,
		},
		{
			version1: "1.2.3-alpha.1",
			version2: "1.2.3-beta.1",
			expected: less,
		},
		{
			version1: "1.2.3-beta.1",
			version2: "1.2.3-alpha.1",
			expected: greater,
		},
		{
			version1: "1.2.3+build",
			version2: "1.2.3",
			expected: equal,
		},
		{
			version1: "1.0.0-alpha",
			version2: "1.0.0-alpha.1",
			expected: less,
		},
		{
			version1: "1.0.0-alpha.1",
			version2: "1.0.0-alpha.beta",
			expected: less,
		},
		{
			version1: "1.0.0-alpha.beta",
			version2: "1.0.0-beta",
			expected: less,
		},
		{
			version1: "1.0.0-beta",
			version2: "1.0.0-beta.2",
			expected: less,
		},
		{
			version1: "1.0.0-beta.2",
			version2: "1.0.0-beta.11",
			expected: less,
		},
		{
			version1: "1.0.0-beta.11",
			version2: "1.0.0-rc.1",
			expected: less,
		},
		{
			version1: "1.0.0-rc.1",
			version2: "1.0.0",
			expected: less,
		},
	}

	for i, tc := range tests {
		expectedStr := ""
		switch tc.expected {
		case greater:
			expectedStr = "greater than"
		case less:
			expectedStr = "less than"
		case equal:
			expectedStr = "equal ot"
		}

		name := fmt.Sprintf("%d: %s is %s %s", i, tc.version1, expectedStr, tc.version2)
		t.Run(name, func(t *testing.T) {
			v1, err := internal.ParseVersion(tc.version1)
			require.NoError(t, err)
			v2, err := internal.ParseVersion(tc.version2)
			require.NoError(t, err)

			actual := internal.Compare(*v1, *v2)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestVersionString(t *testing.T) {
	type test struct {
		version  *internal.Version
		expected string
	}
	tests := []test{
		{
			version: &internal.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			expected: "1.2.3",
		},
		{
			version: &internal.Version{
				Prefix: internal.Ptr("v"),
				Major:  1,
				Minor:  2,
				Patch:  3,
			},
			expected: "v1.2.3",
		},
		{
			version: &internal.Version{
				Prefix:     internal.Ptr("v"),
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: []string{"alpha", "1"},
			},
			expected: "v1.2.3-alpha.1",
		},
		{
			version: &internal.Version{
				Prefix:     internal.Ptr("v"),
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: []string{"alpha", "1"},
				Build:      internal.Ptr("build.123"),
			},
			expected: "v1.2.3-alpha.1+build.123",
		},
		{
			version: &internal.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: internal.Ptr("build.123"),
			},
			expected: "1.2.3+build.123",
		},
		{
			version: &internal.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: []string{"beta"},
			},
			expected: "1.2.3-beta",
		},
	}

	for _, tc := range tests {
		t.Run(tc.expected, func(t *testing.T) {
			actual := tc.version.String()
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestBumpPreRelease(t *testing.T) {
	type test struct {
		version  *internal.Version
		expected *internal.Version
		err      bool
	}

	tests := []test{
		// 1.2.3 -> err
		// 1.2.3-alpha -> 1.2.3-alpha.1
		// 1.2.3-alpha.1 -> 1.2.3-alpha.2
		// 1.2.3-alpha.beta -> 1.2.3-alpha.beta.1
		// 1.2.3-alpha.beta.1 -> 1.2.3-alpha.beta.2
		{
			version: &internal.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			expected: nil,
			err:      true,
		},
		{
			version: &internal.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: []string{"alpha"},
			},
			expected: &internal.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: []string{"alpha", "1"},
			},
			err: false,
		},
		{
			version: &internal.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: []string{"alpha", "1"},
			},
			expected: &internal.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: []string{"alpha", "2"},
			},
			err: false,
		},
		{
			version: &internal.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: []string{"alpha", "beta"},
			},
			expected: &internal.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: []string{"alpha", "beta", "1"},
			},
			err: false,
		},
		{
			version: &internal.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: []string{"alpha", "beta", "1"},
			},
			expected: &internal.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: []string{"alpha", "beta", "2"},
			},
			err: false,
		},
	}

	for i, tc := range tests {
		name := fmt.Sprintf("%d: %s is bumped to %s", i, tc.version, tc.expected)
		t.Run(name, func(t *testing.T) {
			actual, err := internal.BumpPreRelease(tc.version)
			if tc.err {
				require.Error(t, err)
				assert.Nil(t, actual)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.String(), actual.String())
			}
		})
	}
}
