package main

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/koron/go-arch"
)

// TODO: better messaging
// TODO: better logging

const (
	urlWin32 = "http://files.kaoriya.net/vim/vim74-kaoriya-win32.zip"
	urlWin64 = "http://files.kaoriya.net/vim/vim74-kaoriya-win64.zip"
)

type config struct {
	name      string
	url       string
	targetDir string
	dataDir   string
	logDir    string
	tmpDir    string
	varDir    string
}

func (c config) downloadPath() (string, error) {
	u, err := url.Parse(c.url)
	if err != nil {
		return "", err
	}
	return filepath.Join(c.tmpDir, filepath.Base(u.Path)), nil
}

func (c config) recipePath() string {
	return filepath.Join(c.varDir, c.name+"-recipe.txt")
}

func (c config) anchorPath() string {
	return filepath.Join(c.varDir, c.name+"-anchor.txt")
}

func (c config) anchor() (time.Time, error) {
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

func (c config) updateAnchor(t time.Time) error {
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

func (c config) dirs() []string {
	return []string{
		c.targetDir,
		c.dataDir,
		c.logDir,
		c.tmpDir,
		c.varDir,
	}
}

func (c config) prepare() error {
	for _, dir := range c.dirs() {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return err
		}
	}
	return nil
}

func newConfig(dir string) (config, error) {
	exe := filepath.Join(dir, "vim.exe")
	cpu, err := arch.Exe(exe)
	if err != nil {
		return config{}, err
	}
	var name, url string
	dataDir := filepath.Join(dir, "netupvim")
	switch cpu {
	case arch.X86:
		name = "vim74-win32"
		url = urlWin32
	case arch.AMD64:
		name = "vim74-win64"
		url = urlWin64
	}
	return config{
		name:      name,
		url:       url,
		targetDir: dir,
		dataDir:   dataDir,
		logDir:    filepath.Join(dataDir, "log"),
		tmpDir:    filepath.Join(dataDir, "tmp"),
		varDir:    filepath.Join(dataDir, "var"),
	}, nil
}

var errorNotModified = errors.New("not modified")

func download(url, outpath string, pivot time.Time) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	if !pivot.IsZero() {
		t := pivot.UTC().Format(http.TimeFormat)
		req.Header.Set("If-Modified-Since", t)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		f, err := os.Create(outpath)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := io.Copy(f, resp.Body); err != nil {
			return err
		}

	case http.StatusNotModified:
		return errorNotModified

	default:
		return fmt.Errorf("unexpected response: %s", resp.Status)
	}
	return nil
}

type fileInfo struct {
	name string
	size uint64
	hash uint32
}

type compareResult int

const (
	fileNotExist = iota + 1
	fileNotMatch
	fileIsMatch
)

func (info fileInfo) compareWithFile(name string) (compareResult, error) {
	fi, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return fileNotExist, nil
		}
		return 0, err
	}
	if (uint64)(fi.Size()) != info.size {
		return fileNotMatch, nil
	}
	v, err := calcCRC32(name)
	if err != nil {
		if os.IsNotExist(err) {
			return fileNotExist, nil
		}
		return 0, err
	}
	if v != info.hash {
		return fileNotMatch, nil
	}
	return fileIsMatch, nil
}

type fileInfoTable map[string]fileInfo

const fileInfoFormat = "%s\t%d\t%08x\n"

func loadFileInfo(fname string) (fileInfoTable, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b := bufio.NewReader(f)
	t := make(fileInfoTable)
	for {
		l, err := b.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		fi := fileInfo{}
		if _, err := fmt.Sscanf(l, fileInfoFormat, &fi.name, &fi.size, &fi.hash); err != nil {
			return nil, err
		}
		t[fi.name] = fi
	}
	return t, nil
}

func saveFileInfo(fname string, t fileInfoTable) error {
	f, err := os.Create(fname)
	if err != nil {
		return nil
	}
	defer f.Close()
	for _, v := range t {
		fmt.Fprintf(f, fileInfoFormat, v.name, v.size, v.hash)
	}
	return f.Sync()
}

func stripPath(name string, count int) string {
	s := strings.Split(name, "/")
	return path.Join(s[count:]...)
}

func calcCRC32(name string) (uint32, error) {
	r, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer r.Close()
	h := crc32.NewIEEE()
	if _, err := io.Copy(h, r); err != nil {
		return 0, err
	}
	return h.Sum32(), nil
}

func splitExt(name string) (base, ext string) {
	ext = path.Ext(name)
	base = name[:len(name)-len(ext)]
	return base, ext

}

func evacuateName(name string) string {
	base, ext := splitExt(name)
	return base + ".orig" + ext
}

func rotateName(name string, index int) string {
	if index == 0 {
		return name
	}
	base, ext := splitExt(name)
	return base + "." + strconv.Itoa(index) + ext
}

func rotateFiles(name string, max int) error {
	last := rotateName(name, max)
	err := os.Remove(last)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for i := max - 1; i >= 0; i-- {
		curr := rotateName(name, i)
		err := os.Rename(curr, last)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		last = curr
	}
	return nil
}

func extractZipFile(zf *zip.File, name string) error {
	if err := os.MkdirAll(filepath.Dir(name), 0777); err != nil {
		return err
	}
	r, err := zf.Open()
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := os.Create(name)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, r)
	return err
}

func extractZip(zipName, dir string, prev fileInfoTable) (fileInfoTable, error) {
	// extract zip file.
	zr, err := zip.OpenReader(zipName)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	curr := make(fileInfoTable)
	for _, zf := range zr.File {
		if zf.Mode().IsDir() {
			continue
		}
		zfName := stripPath(zf.Name, 1)
		curr[zfName] = fileInfo{
			name: zfName,
			size: zf.UncompressedSize64,
			hash: zf.CRC32,
		}
		outName := filepath.Join(dir, zfName)
		// evacuation and optimization.
		if p, ok := prev[zfName]; ok {
			r, err := p.compareWithFile(outName)
			if err != nil {
				// TODO: log an error.
				continue
			}
			switch r {
			case fileNotMatch:
				outName = evacuateName(outName)
			case fileIsMatch:
				// skip un-changed files.
				if p.hash == zf.CRC32 {
					continue
				}
			}
		}
		// rotation.
		ext := strings.ToLower(path.Ext(zfName))
		if ext == ".exe" || ext == ".dll" {
			if err := rotateFiles(outName, 5); err != nil {
				return nil, err
			}
		}
		if err := extractZipFile(zf, outName); err != nil {
			return nil, err
		}
	}
	return curr, nil
}

// cleanFiles removes unused/untracked files.
func cleanFiles(dir string, prev, curr fileInfoTable) {
	for _, p := range prev {
		if _, ok := curr[p.name]; ok {
			continue
		}
		fpath := filepath.Join(dir, p.name)
		if r, _ := p.compareWithFile(fpath); r != fileIsMatch {
			continue
		}
		os.Remove(fpath)
	}
}

func extract(dir, zipName, recipeName string) error {
	prev, err := loadFileInfo(recipeName)
	if err != nil {
		log.Printf("WARN: failed to load recipe: %s", err)
		log.Println("INFO: try to extract all files")
		prev = make(fileInfoTable)
	}
	curr, err := extractZip(zipName, dir, prev)
	if err != nil {
		return err
	}
	if err := saveFileInfo(recipeName, curr); err != nil {
		log.Printf("WARN: failed to save recipe: %s", err)
	}
	cleanFiles(dir, prev, curr)
	return nil
}

func update(c config) error {
	dp, err := c.downloadPath()
	if err != nil {
		return err
	}
	anchor, err := c.anchor()
	if err != nil {
		return err
	}
	if err := download(c.url, dp, anchor); err != nil {
		if err == errorNotModified {
			return nil
		}
		return err
	}
	anchor = time.Now()
	if err := extract(c.targetDir, dp, c.recipePath()); err != nil {
		return err
	}
	if err := c.updateAnchor(anchor); err != nil {
		os.Remove(c.anchorPath())
		return err
	}
	if err := os.Remove(dp); err != nil {
		log.Printf("WARN: failed to remove: %s", err)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("USAGE: netupvim {TARGET_DIR}")
		os.Exit(1)
	}
	c, err := newConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	if err := c.prepare(); err != nil {
		log.Fatal(err)
	}
	if err := update(c); err != nil {
		log.Fatal(err)
	}
}
