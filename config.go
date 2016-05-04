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

	// CPU is target CPU architecture: "x86" or "amd64"
	CPU string `toml:"cpu"`

	// GithubUser is username which be used for github's basic auth.
	GithubUser string `toml:"github_user"`

	// GithubUser is token which be used for github's basic auth.
	GithubToken string `toml:"github_token"`
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
