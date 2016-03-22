package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/koron/go-arch"
)

// TODO: better messaging
// TODO: better logging

const (
	urlWin32 = "http://files.kaoriya.net/vim/vim74-kaoriya-win32.zip"
	urlWin64 = "http://files.kaoriya.net/vim/vim74-kaoriya-win64.zip"
)

var (
	errorNotModified = errors.New("not modified")
)

func determineSourceURL(c *config) (string, error) {
	switch c.cpu {
	case arch.X86:
		return urlWin32, nil
	case arch.AMD64:
		return urlWin64, nil
	default:
		return "", nil
	}
}

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
	curr, err := extractZip(zipName, dir, 1, prev)
	if err != nil {
		return err
	}
	if err := saveFileInfo(recipeName, curr); err != nil {
		log.Printf("WARN: failed to save recipe: %s", err)
	}
	cleanFiles(dir, prev, curr)
	return nil
}

func update(c *config) error {
	srcURL, err := determineSourceURL(c)
	if err != nil {
		return err
	}
	dp, err := c.downloadPath(srcURL)
	if err != nil {
		return err
	}
	anchor, err := c.anchor()
	if err != nil {
		return err
	}
	if err := download(srcURL, dp, anchor); err != nil {
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
