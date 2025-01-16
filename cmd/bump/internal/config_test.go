package internal_test

import (
	"testing"

	"testing/fstest"

	"github.com/mrvinkel/bump/cmd/bump/internal"
	"github.com/stretchr/testify/assert"
)

func TestReadConfigNotExist(t *testing.T) {
	fs := fstest.MapFS{}

	config, err := internal.ReadConfig(fs)

	assert.Nil(t, err)
	assert.Nil(t, config)
}

func TestReadConfig(t *testing.T) {
	fs := fstest.MapFS{
		internal.CONFIG_FILE: &fstest.MapFile{
			Data: []byte(`{"commit": false}`),
		},
	}

	config, err := internal.ReadConfig(fs)

	assert.Nil(t, err)
	assert.NotNil(t, config)

	assert.False(t, *config.Commit)

	// defaults
	assert.Equal(t, "release ${version}", *config.Message)
	assert.Nil(t, config.Prefix)
	assert.Nil(t, config.Fetch)
	assert.Nil(t, config.Verify)
	assert.Equal(t, "/bin/bash -c", *config.Shell)
	assert.Empty(t, config.PreHook)
}
