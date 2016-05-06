package netup

import (
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/koron/go-arch"
)

// Config represents netup's configuration.
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

	// DownloadTimeout is timeout for downloading archive (default: "5min")
	DownloadTimeout string
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

func (c *Config) getSource() string {
	if c.Source != "" {
		return c.Source
	}
	return "release"
}

func (c *Config) getTargetDir() string {
	if c.TargetDir != "" {
		return c.TargetDir
	}
	dir, err := os.Getwd()
	if err != nil {
		logFatal(err)
	}
	return dir
}

func (c *Config) getCPU() arch.CPU {
	return arch.ParseCPU(c.CPU)
}

func (c *Config) getGithubUser() string {
	if c.GithubUser != "" {
		return c.GithubUser
	}
	v, _ := os.LookupEnv("NETUPVIM_GITHUB_USER")
	return v
}

func (c *Config) getGithubToken() string {
	if c.GithubToken != "" {
		return c.GithubToken
	}
	v, _ := os.LookupEnv("NETUPVIM_GITHUB_TOKEN")
	return v
}

func (c *Config) getDownloadTimeout() time.Duration {
	if c.DownloadTimeout == "" {
		return 5 * time.Minute
	}
	t, err := time.ParseDuration(c.DownloadTimeout)
	if err != nil {
		logFatal(err)
	}
	return t
}
