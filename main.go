package main

import (
	"errors"
	"flag"
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

const (
	urlWin32     = "http://files.kaoriya.net/vim/vim74-kaoriya-win32.zip"
	urlWin64     = "http://files.kaoriya.net/vim/vim74-kaoriya-win64.zip"
	urlWin32Test = "http://files.kaoriya.net/vim/vim74-kaoriya-win32-test.zip"
	urlWin64Test = "http://files.kaoriya.net/vim/vim74-kaoriya-win64-test.zip"
)

var (
	errNotModified     = errors.New("not modified")
	errUnsupportedArch = errors.New("unsupported architecture")
)

var (
	help    = flag.Bool("h", false, "show this message")
	target  = flag.String("t", mustGetwd(), "target dir to upgrade/install")
	restore = flag.Bool("r", false, "restore all files (WIP)")
	beta    = flag.Bool("b", false, "use beta release")
)

func mustGetwd() string {
	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return d
}

func determineSourceURL(c *config) (string, error) {
	if *beta {
		switch c.cpu {
		case arch.X86:
			return urlWin32Test, nil
		case arch.AMD64:
			return urlWin64Test, nil
		default:
			return "", errUnsupportedArch
		}
	}
	switch c.cpu {
	case arch.X86:
		return urlWin32, nil
	case arch.AMD64:
		return urlWin64, nil
	default:
		return "", errUnsupportedArch
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
		return errNotModified

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
		if err == errNotModified {
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

func showHelp() {
	fmt.Fprintf(os.Stderr, `%[1]s is tool to upgrade/install Vim (+kaoriya) in/to target dir.

Usage: %[1]s [options]

Options are:
`, filepath.Base(os.Args[0]))
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if *help {
		showHelp()
		os.Exit(1)
	}
	c, err := newConfig(*target)
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
