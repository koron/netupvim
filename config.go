package main

// Config represents netupvim's configuration.
type Config struct {

	// Source is source of update: release, develop and canary.
	// Default is "release"
	Source string

	// TargetDir is target directory to update.  Default is current working
	// directory.
	TargetDir string
}

func loadConfig(name string) (*Config, error) {
	// TODO:
	return &Config{}, nil
}

func (c *Config) source(defaultValue string) string {
	if c.Source == "" {
		return defaultValue
	}
	return c.Source
}

func (c *Config) targetDir(defaultValue string) string {
	if c.TargetDir == "" {
		return defaultValue
	}
	return c.TargetDir
}
