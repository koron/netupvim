package netup

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
	source    string
	targetDir string
	dataDir   string
	logDir    string
	tmpDir    string
	varDir    string
}

func newContext(dir, src, name, exe string) (*context, error) {
	cpu, err := arch.Exe(filepath.Join(dir, exe))
	if err != nil {
		return nil, err
	}
	dataDir := filepath.Join(dir, name)
	return &context{
		cpu:       cpu,
		source:    src,
		targetDir: dir,
		dataDir:   dataDir,
		logDir:    filepath.Join(dataDir, "log"),
		tmpDir:    filepath.Join(dataDir, "tmp"),
		varDir:    filepath.Join(dataDir, "var"),
	}, nil
}

func (c *context) downloadPath(targetURL string) (string, error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return "", err
	}
	return filepath.Join(c.tmpDir, filepath.Base(u.Path)), nil
}

func (c *context) recipePath() string {
	return filepath.Join(c.varDir, "recipe.txt")
}

func (c *context) anchorPath() string {
	return filepath.Join(c.varDir, "anchor.txt")
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

func (c *context) resetAnchor() error {
	if err := os.Remove(c.anchorPath()); os.IsExist(err) {
		return err
	}
	return nil
}

func (c *context) resetRecipe() error {
	if err := os.Remove(c.recipePath()); os.IsExist(err) {
		return err
	}
	return nil
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
