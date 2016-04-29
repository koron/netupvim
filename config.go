package main

import (
	"os"

	"github.com/BurntSushi/toml"
)

// Config represents netupvim's configuration.
type Config struct {

	// Source is source of update: release, develop and canary.
	// Default is "release"
	Source string `toml:"source"`

	// TargetDir is target directory to update.  Default is current working
	// directory.
	TargetDir string `toml:"target_dir"`
}

func loadConfig(name string) (*Config, error) {
	var conf Config
	_, err := toml.DecodeFile(name, &conf)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	return &conf, nil
}
