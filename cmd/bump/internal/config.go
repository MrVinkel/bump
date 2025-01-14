package internal

import (
	"encoding/json"
	"io/fs"

	_ "embed"
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
	Debug   *bool    `json:"debug"`
	Shell   *string  `json:"shell"`
	PreHook []string `json:"preHook"`
}

func ReadConfig(fs fs.FS) (bool, *Config, error) {
	file, err := fs.Open(CONFIG_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil, nil
		}
		return false, nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return false, nil, err
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return false, nil, err
	}

	setDefaults(&config)

	return true, &config, nil
}

func setDefaults(config *Config) {
	if config.Commit == nil {
		config.Commit = Ptr(true)
	}
	if config.Message == nil {
		config.Message = Ptr("release ${version}")
	}
	if config.Prefix == nil {
		config.Prefix = Ptr("v")
	}
	if config.Fetch == nil {
		config.Fetch = Ptr(true)
	}
	if config.Verify == nil {
		config.Verify = Ptr(true)
	}
	if config.Debug == nil {
		config.Debug = Ptr(false)
	}
	if config.Shell == nil {
		config.Shell = Ptr("/bin/bash")
	}
}
