package netup

import (
	"io"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/xerrors"
)

type context struct {
	targetDir string
	dataDir   string
	logDir    string
	tmpDir    string
	varDir    string
	source    Source
}

func newContext(targetDir, workDir string, srcPack SourcePack, arch Arch) (*context, error) {
	// deterine source.
	cpu, err := arch.detectCPU(targetDir)
	if err != nil {
		return nil, xerrors.Errorf("failed to detect CPU: %w", err)
	}
	src, ok := srcPack[cpu]
	if !ok {
		return nil, xerrors.Errorf("unsupported arch: %+v", arch)
	}

	return &context{
		targetDir: targetDir,
		dataDir:   workDir,
		logDir:    filepath.Join(workDir, "log"),
		tmpDir:    filepath.Join(workDir, "tmp"),
		varDir:    filepath.Join(workDir, "var", src.name()),
		source:    src,
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

func (c *context) mkdirAll() error {
	for _, dir := range c.dirs() {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return err
		}
	}
	return nil
}
