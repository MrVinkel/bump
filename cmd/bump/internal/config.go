package internal

import (
	"encoding/json"
	"io/fs"

	"io"
	"os"
)

const (
	CONFIG_FILE = ".bump.json"
)

type Config struct {
	Commit  *bool    `json:"commit"`
	Message *string  `json:"message"`
	Prefix  *string  `json:"prefix"`
	Fetch   *bool    `json:"fetch"`
	Verify  *bool    `json:"verify"`
	Shell   *string  `json:"shell"`
	PreHook []string `json:"preHook"`
}

func ReadConfig(fs fs.FS) (*Config, error) {
	file, err := fs.Open(CONFIG_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close() // nolint:errcheck

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	setDefaults(&config)

	return &config, nil
}

func setDefaults(config *Config) {
	if config.Message == nil {
		config.Message = Ptr("release ${version}")
	}
	if config.Shell == nil || *config.Shell == "" {
		config.Shell = Ptr("/bin/bash -c")
	}
}
