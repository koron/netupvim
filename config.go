package main

import (
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/koron/go-arch"
	"github.com/koron/netupvim/netup"
)

// config represents netup's configuration.
type config struct {

	// Source is source of update: release, develop and canary.
	// Default is "release"
	Source string `toml:"source"`

	// TargetDir is target directory to update.  Default is current working
	// directory.
	TargetDir string `toml:"target_dir"`

	// CPU is target CPU architecture: "x86" or "amd64"
	CPU string `toml:"cpu"`

	// GithubUser is username which be used for github's basic auth.
	// DEPRECATED.
	GithubUser string `toml:"github_user"`

	// GithubToken is token which be used for github's basic auth.
	GithubToken string `toml:"github_token"`

	// GithubVerbose enables log for github related operation.
	GithubVerbose bool `toml:"github_verbose"`

	// DownloadTimeout is timeout for downloading archive (default: "5min")
	DownloadTimeout string

	// LogRotateCount is used for log rotation.
	LogRotateCount int `toml:"log_rotate_count"`

	// ExeRotateCount is used for executable files rotation.
	ExeRotateCount int `toml:"exe_rotate_count"`

	// DisableSelfUpdate disables netupvim's self update.
	DisableSelfUpdate bool `toml:"disable_self_update"`
}

func loadConfig(name string) (*config, error) {
	var conf config
	_, err := toml.DecodeFile(name, &conf)
	if err != nil {
		if os.IsNotExist(err) {
			return &config{}, nil
		}
		return nil, err
	}
	return &conf, nil
}

func (c *config) getSource() string {
	if c.Source != "" {
		return c.Source
	}
	return "release"
}

func (c *config) getTargetDir() string {
	if c.TargetDir != "" {
		return c.TargetDir
	}
	dir, err := os.Getwd()
	if err != nil {
		netup.LogFatal(err)
	}
	return dir
}

func (c *config) getCPU() arch.CPU {
	return arch.ParseCPU(c.CPU)
}

func (c *config) getGithubUser() string {
	if c.GithubUser != "" {
		return c.GithubUser
	}
	v, _ := os.LookupEnv("NETUPVIM_GITHUB_USER")
	return v
}

func (c *config) getGithubToken() string {
	if c.GithubToken != "" {
		return c.GithubToken
	}
	v, _ := os.LookupEnv("NETUPVIM_GITHUB_TOKEN")
	return v
}

func (c *config) getDownloadTimeout() time.Duration {
	if c.DownloadTimeout == "" {
		return 5 * time.Minute
	}
	t, err := time.ParseDuration(c.DownloadTimeout)
	if err != nil {
		netup.LogFatal(err)
	}
	return t
}
