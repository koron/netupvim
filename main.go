package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/koron/go-github"
	"github.com/koron/netupvim/netup"
)

var (
	version = "none"
)

var (
	targetDir  = "."
	sourceName = "release"
	cpu        string
	selfUpdate = true

	restoreMode bool
	sweepMode   bool
)

func setup() error {
	conf, err := loadConfig("netupvim.ini")
	if err != nil {
		return err
	}

	// Parse options.
	var (
		helpOpt    = flag.Bool("h", false, "show this message")
		targetOpt  = flag.String("t", conf.getTargetDir(), "target dir to upgrade/install")
		sourceOpt  = flag.String("s", conf.getSource(), "source of update: release,develop,canary")
		versionOpt = flag.Bool("version", false, "show version")
	)
	flag.BoolVar(&restoreMode, "restore", false, "force download & extract all files")
	flag.BoolVar(&sweepMode, "sweep", false, "sweep rotated files (.1.exe, .2.exe or so)")
	flag.Parse()
	if *helpOpt {
		showHelp()
		os.Exit(1)
	}
	if *versionOpt {
		showVersion()
		os.Exit(1)
	}

	if sweepMode && restoreMode {
		return errors.New(`can't be combined with "-sweep" and "-restore" flags`)
	}

	// setup context.
	targetDir = *targetOpt
	sourceName = *sourceOpt
	cpu = conf.CPU
	selfUpdate = !conf.DisableSelfUpdate

	netup.Version = version
	netup.DownloadTimeout = conf.getDownloadTimeout()
	netup.GithubUser = conf.getGithubUser()
	github.DefaultClient.Token = conf.getGithubToken()
	netup.GithubVerbose = conf.GithubVerbose
	if conf.LogRotateCount > 0 {
		netup.LogRotateCount = conf.LogRotateCount
	}
	if conf.ExeRotateCount > 0 {
		netup.ExeRotateCount = conf.ExeRotateCount
	}

	return nil
}

func shouldSelfUpdate() bool {
	if !selfUpdate {
		return false
	}
	_, err := os.Stat(filepath.Join(targetDir, "netupvim.exe"))
	return err == nil
}

func run() error {
	if err := setup(); err != nil {
		return err
	}
	workDir := filepath.Join(targetDir, "netupvim")
	// update vim
	vimPack, ok := vimSet[sourceName]
	if !ok {
		return fmt.Errorf("invalid source: %s", sourceName)
	}

	if sweepMode {
		err := netup.Sweep(targetDir, workDir, vimPack,
			netup.Arch{Name: cpu, Hint: "vim.exe"})
		if err != nil {
			return err
		}
		err = netup.Sweep(targetDir, workDir, netupPack,
			netup.Arch{Name: "X86"})
		if err != nil {
			return err
		}
		return nil
	}

	err := netup.Update(
		targetDir,
		workDir,
		vimPack,
		netup.Arch{Name: cpu, Hint: "vim.exe"},
		restoreMode)
	if err != nil {
		return err
	}
	// try to update netupvim
	if shouldSelfUpdate() {
		netup.LogInfo("trying to update netupvim")
		err := netup.Update(
			targetDir,
			workDir,
			netupPack,
			netup.Arch{Name: "X86"},
			restoreMode)
		if err != nil {
			netup.LogInfo("failed to udate netupvim: %s", err)
		}
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

func showVersion() {
	fmt.Fprintf(os.Stderr, "netupvim version %s\n", version)
}

func main() {
	if err := run(); err != nil {
		netup.LogFatal(err)
	}
}
