package netup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
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
	logInfo("extract completed successfully")
	return nil
}

func update(c *context) error {
	src := c.source
	logInfo("determined source: %s", src.String())
	t, err := c.anchor()
	if err != nil {
		return err
	}
	last := -1
	p, err := src.download(c.tmpDir, t, func(curr, max int64) {
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
	if err := extract(p, c.targetDir, src.stripCount(), c.recipePath()); err != nil {
		return err
	}
	if err := c.updateAnchor(t); err != nil {
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
