package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/koron/go-arch"
	"github.com/koron/go-github"
)

const logRotateCount = 5

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
		return err
	}
	defer f.Close()
	for _, v := range t {
		_, err := fmt.Fprintf(f, fileInfoFormat, v.name, v.size, v.hash)
		if err != nil {
			return err
		}
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
		logInfo("remove unused file %s", fpath)
	}
}

func extract(dir, zipName, recipeName string) error {
	prev, err := loadFileInfo(recipeName)
	if err != nil {
		logLoadRecipeFailed(err)
		prev = make(fileInfoTable)
	}
	logInfo("extract archive: %s", zipName)
	msgPrintf("extract archive\n")
	last := -1
	curr, err := extractZip(zipName, dir, 1, prev, func(curr, max uint64) {
		v := int(curr * 100 / max)
		if v != last {
			msgPrintProgress(v)
			last = v
		}
	})
	msgPrintln()
	if err != nil {
		return err
	}
	if err := saveFileInfo(recipeName, curr); err != nil {
		logSaveRecipeFailed(err)
	}
	cleanFiles(dir, prev, curr)
	logInfo("extract completed successfully")
	return nil
}

func update(c *context) error {
	t, err := c.anchor()
	if err != nil {
		return err
	}
	s, err := determineSource(c.source, c.cpu)
	if err != nil {
		return err
	}
	logInfo("determined source: %s", s.String())
	last := -1
	p, err := s.download(c.tmpDir, t, func(curr, max int64) {
		v := int(curr * 100 / max)
		if v != last {
			msgPrintProgress(v)
			last = v
		}
	})
	msgPrintln()
	if err != nil {
		if err == errSourceNotModified {
			logInfo("no updates found")
			err = nil
		}
		return err
	}
	logInfo("download completed successfully")
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
		logCleanArchiveFailed(err)
	}
	return nil
}

func restore(c *context) error {
	if err := os.Remove(c.anchorPath()); os.IsExist(err) {
		return err
	}
	if err := os.Remove(c.recipePath()); os.IsExist(err) {
		return err
	}
	logInfo("deleted anchor and recipe to restore")
	return update(c)
}

func showHelp() {
	fmt.Fprintf(os.Stderr, `%[1]s is tool to upgrade/install Vim (+kaoriya) in/to target dir.

Usage: %[1]s [options]

Options are:
`, filepath.Base(os.Args[0]))
	flag.PrintDefaults()
}

func main() {
	// Fetch environment variables.
	var (
		githubUser  string
		githubToken string
	)
	githubUser, _ = os.LookupEnv("NETUPVIM_GITHUB_USER")
	githubToken, _ = os.LookupEnv("NETUPVIM_GITHUB_TOKEN")

	// Load config.
	conf, err := loadConfig("netupvim.ini")
	if err != nil {
		logFatal(err)
	}
	var (
		defaultSource = "release"
		defaultTarget = mustGetwd()
		cpu           arch.CPU
	)
	if conf.Source != "" {
		defaultSource = conf.Source
	}
	if conf.TargetDir != "" {
		defaultTarget = conf.TargetDir
	}
	if conf.CPU != "" {
		cpu = arch.ParseCPU(conf.CPU)
	}
	if conf.GithubUser != "" {
		githubUser = conf.GithubUser
	}
	if conf.GithubToken != "" {
		githubToken = conf.GithubToken
	}

	// Parse options.
	var (
		helpOpt    = flag.Bool("h", false, "show this message")
		targetOpt  = flag.String("t", defaultTarget, "target dir to upgrade/install")
		sourceOpt  = flag.String("s", defaultSource, "source of update: release,develop,canary")
		restoreOpt = flag.Bool("restore", false, "force download & extract all files")
	)
	flag.Parse()
	if *helpOpt {
		showHelp()
		os.Exit(1)
	}

	// Setup context and other components.
	c, err := newContext(*targetOpt, *sourceOpt)
	if err != nil {
		logFatal(err)
	}
	if cpu != 0 {
		c.cpu = cpu
	}
	if githubUser != "" && githubToken != "" {
		github.DefaultClient.Username = githubUser
		github.DefaultClient.Token = githubToken
	}
	if err := c.prepare(); err != nil {
		logFatal(err)
	}
	logSetup(c.logDir, logRotateCount)
	logInfo("context: CPU=%s source=%s dir=%s",
		c.cpu.String(), c.source.String(), c.targetDir)
	github.DefaultClient.Logger = logger

	// Run update.
	proc := update
	if *restoreOpt {
		proc = restore
	}
	if err := proc(c); err != nil {
		logFatal(err)
	}
}
