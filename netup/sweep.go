package netup

import (
	"github.com/koron/netupvim/internal/rotate"
)

// Sweep deletes rotated files which likes: "vim.1.exe", "vim.2.exe" or so.
func Sweep(targetDir, workDir string, srcPack SourcePack, arch Arch) error {
	ctx, err := newContext(targetDir, workDir, srcPack, arch)
	if err != nil {
		return err
	}
	logSetup(ctx.logDir, LogRotateCount)
	curr, err := loadFileInfo(ctx.recipePath())
	if err != nil {
		logLoadRecipeFailed(err)
		return err
	}

	for _, fi := range curr {
		err := sweepFile(fi.name)
		if err != nil {
			return err
		}
	}

	return nil
}

func sweepFile(name string) error {
	if !rotate.IsTarget(name) {
		return nil
	}
	err := rotate.Sweep(name, ExeRotateCount)
	if err != nil {
		return err
	}
	return nil
}
