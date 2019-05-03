package netup

import (
	"fmt"
	"os"
	"path/filepath"
)

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
func cleanFiles(dir string, prev, curr fileInfoTable) error {
	var lastError error
	for _, p := range prev {
		if _, ok := curr[p.name]; ok {
			continue
		}
		fpath := filepath.Join(dir, p.name)
		r, err := p.compareWithFile(fpath)
		if err != nil {
			logWarn("failed to compare %s: %s", fpath, err)
			lastError = err
			continue
		}
		if r != fileIsMatch {
			logInfo("keep %s because it is modified", fpath)
			continue
		}

		err = os.Remove(fpath)
		if err != nil {
			logWarn("failed to remove file %s: %s", fpath, err)
			lastError = err
		} else {
			logInfo("remove unused file %s", fpath)
		}
		err = sweepFile(fpath)
		if err != nil {
			logWarn("failed to sweep %s: %s", fpath, err)
			lastError = err
		}
	}
	return lastError
}

type dirInfo struct {
	path    string
	name    string
	keep    bool
	subdirs map[string]*dirInfo
}

func (di *dirInfo) add(fi fileInfo) {
	names := fi.dirList()
	if len(names) == 0 {
		return
	}
	curr := di
	for i, n := range names {
		sub, ok := curr.subdirs[n]
		if ok {
			curr = sub
			continue
		}
		sub = &dirInfo{
			path:    filepath.Join(names[:i+1]...),
			name:    n,
			subdirs: map[string]*dirInfo{},
		}
		curr.subdirs[n] = sub
		curr = sub
	}
}

// cleanDirs removes unused/untracked directory.
func cleanDirs(dir string, prev fileInfoTable) error {
	// build directory tree.
	root := &dirInfo{
		keep:    true,
		subdirs: map[string]*dirInfo{},
	}
	for _, fi := range prev {
		root.add(fi)
	}

	rmdir := func(di *dirInfo) error {
		if di.keep {
			return nil
		}
		for _, sub := range di.subdirs {
			if sub.keep {
				di.keep = true
				return nil
			}
		}
		path := filepath.Join(dir, di.path)
		err := os.Remove(path)
		if err != nil {
			di.keep = true
			logInfo("keep dir %s: %s", path, err)
			return err
		}
		logInfo("deleted empty dir %s", path)
		return nil
	}

	// remove dirs from edges.
	var lastError error
	var rmdirs func(di *dirInfo)
	rmdirs = func(di *dirInfo) {
		for _, sub := range di.subdirs {
			rmdirs(sub)
		}
		// ignore rmdir error
		err := rmdir(di)
		if err != nil {
			lastError = err
		}
	}
	rmdirs(root)
	return lastError
}

func extract(zipName, dir string, stripCount int, recipeName string) error {
	prev, err := loadFileInfo(recipeName)
	if err != nil {
		logLoadRecipeFailed(err)
		prev = make(fileInfoTable)
	}
	logInfo("extract archive: %s", zipName)
	msgPrintf("extract archive\n")
	last := -1
	curr, err := extractZip(zipName, dir, stripCount, prev, func(curr, max uint64) {
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
	cleanDirs(dir, prev)
	logInfo("extract completed successfully")
	return nil
}

func update(c *context) error {
	src := c.source
	logInfo("determined source: %s", src.String())
	at, err := c.anchor()
	if err != nil {
		return err
	}
	last := -1
	p, mt, err := src.download(c.tmpDir, at, func(curr, max int64) {
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
	if err := extract(p, c.targetDir, src.stripCount(), c.recipePath()); err != nil {
		return err
	}
	if err := c.updateAnchor(mt); err != nil {
		c.resetAnchor()
		return err
	}
	if err := os.Remove(p); err != nil {
		logCleanArchiveFailed(err)
	}
	return nil
}

func restore(c *context) error {
	if err := c.resetAnchor(); err != nil {
		return err
	}
	if err := c.resetRecipe(); err != nil {
		return err
	}
	logInfo("deleted anchor and recipe to restore")
	return update(c)
}
