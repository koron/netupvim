package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// TODO: better messaging

var (
	help   = flag.Bool("h", false, "show this message")
	target = flag.String("t", mustGetwd(), "target dir to upgrade/install")
	src    = flag.String("s", "release", "source of update: release,develop,canary")
)

func mustGetwd() string {
	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return d
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
	t, err := c.anchor()
	if err != nil {
		return err
	}
	s, err := determineSource(c.source, c.cpu)
	if err != nil {
		return err
	}
	p, err := s.download(c.tmpDir, t)
	if err != nil {
		if err == errSourceNotModified {
			err = nil
		}
		return err
	}
	// capture anchor's new value.
	t = time.Now()
	if err := extract(c.targetDir, p, c.recipePath()); err != nil {
		return err
	}
	if err := c.updateAnchor(t); err != nil {
		os.Remove(c.anchorPath())
		return err
	}
	if err := os.Remove(p); err != nil {
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
	c, err := newConfig(*target, *src)
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
