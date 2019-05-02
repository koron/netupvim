package netup

import (
	"github.com/koron/go-github"
)

// Update updates or installs a package into target directory.
func Update(targetDir, workDir string, srcPack SourcePack, arch Arch, restoreFlag bool) error {
	ctx, err := newContext(targetDir, workDir, srcPack, arch)
	if err != nil {
		return err
	}
	err = ctx.mkdirAll()
	if err != nil {
		return err
	}

	logSetup(ctx.logDir, LogRotateCount)
	if GithubUser != "" {
		logWarn("GithubUser (from config or env) is deprecated and ignored")
	}
	if GithubVerbose {
		github.DefaultClient.Logger = logger
	}
	logInfo("context: target=%s source=%s", ctx.targetDir, ctx.source)

	// Run update.
	proc := update
	if restoreFlag {
		proc = restore
	}
	if err := proc(ctx); err != nil {
		return err
	}

	return nil
}
