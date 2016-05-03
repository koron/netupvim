package main

import (
	"io"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/koron/go-arch"
)

type context struct {
	cpu       arch.CPU
	source    sourceType
	targetDir string
	dataDir   string
	logDir    string
	tmpDir    string
	varDir    string
}

func newContext(dir, src string) (*context, error) {
	exe := filepath.Join(dir, "vim.exe")
	cpu, err := arch.Exe(exe)
	if err != nil {
		return nil, err
	}
	dataDir := filepath.Join(dir, "netupvim")
	st, err := toSourceType(src)
	if err != nil {
		return nil, err
	}
	return &context{
		cpu:       cpu,
		source:    st,
		targetDir: dir,
		dataDir:   dataDir,
		logDir:    filepath.Join(dataDir, "log"),
		tmpDir:    filepath.Join(dataDir, "tmp"),
		varDir:    filepath.Join(dataDir, "var"),
	}, nil
}

func (c *context) name() string {
	switch c.cpu {
	case arch.X86:
		return "vim74-win32"
	case arch.AMD64:
		return "vim74-win64"
	default:
		return ""
	}
}

func (c *context) downloadPath(targetURL string) (string, error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return "", err
	}
	return filepath.Join(c.tmpDir, filepath.Base(u.Path)), nil
}

func (c *context) recipePath() string {
	return filepath.Join(c.varDir, c.name()+"-recipe.txt")
}

func (c *context) anchorPath() string {
	return filepath.Join(c.varDir, c.name()+"-anchor.txt")
}

func (c *context) anchor() (time.Time, error) {
	f, err := os.Open(c.anchorPath())
	if err != nil {
		return time.Time{}, nil
	}
	defer f.Close()
	buf := make([]byte, 25)
	if _, err := io.ReadFull(f, buf); err != nil {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, string(buf))
}

func (c *context) updateAnchor(t time.Time) error {
	f, err := os.Create(c.anchorPath())
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.WriteString(f, t.Format(time.RFC3339)); err != nil {
		return err
	}
	return f.Sync()
}

func (c *context) dirs() []string {
	return []string{
		c.targetDir,
		c.dataDir,
		c.logDir,
		c.tmpDir,
		c.varDir,
	}
}

func (c *context) prepare() error {
	for _, dir := range c.dirs() {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return err
		}
	}
	return nil
}
