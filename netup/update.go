package netup

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/koron/go-github"
)

var (
	// DownloadTimeout is timeout for download file.
	DownloadTimeout = 5 * time.Minute

	// GithubUser is username which be used for github's basic auth.
	GithubUser string

	// GithubToken is token which be used for github's basic auth.
	GithubToken string

	// GithubVerbose enables log for github related operation.
	GithubVerbose bool

	// LogRotateCount is used for log rotation.
	LogRotateCount = 5

	// ExeRotateCount is used for executable files rotation.
	ExeRotateCount = 5
)

// Update updates or installs a package into target directory.
func Update(targetDir, workDir string, srcPack SourcePack, arch Arch, restoreFlag bool) error {
	// deterine source.
	cpu, err := arch.detectCPU(targetDir)
	if err != nil {
		return err
	}
	src, ok := srcPack[cpu]
	if !ok {
		return fmt.Errorf("unsupported arch: %+v", arch)
	}

	// setup environment.
	downloadTimeout = DownloadTimeout
	if GithubToken != "" {
		github.DefaultClient.Username = GithubUser
		github.DefaultClient.Token = GithubToken
	}

	ctx := &context{
		targetDir: targetDir,
		dataDir:   workDir,
		logDir:    filepath.Join(workDir, "log"),
		tmpDir:    filepath.Join(workDir, "tmp"),
		varDir:    filepath.Join(workDir, "var", src.name()),
		source:    src,
	}
	if err := ctx.mkdirAll(); err != nil {
		return err
	}

	logSetup(ctx.logDir, LogRotateCount)
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
